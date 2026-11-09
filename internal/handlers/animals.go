package handlers

import (
	"errors"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"animal-api/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnimalHandler struct {
	Animals *mongo.Collection
}

var sortableFields = map[string]string{
	"name":      "name",
	"age":       "age",
	"weight":    "weight",
	"species":   "species",
	"createdAt": "createdAt",
}

// List returns animals with filtering, sorting and pagination.
//
//	GET /api/animals?species=dog&adopted=false&minAge=1&maxAge=10&search=rex
//	                &sort=age&order=desc&page=1&limit=10
func (h *AnimalHandler) List(c *gin.Context) {
	// ---- Filtering ----
	filter := bson.M{}
	if s := c.Query("species"); s != "" {
		filter["species"] = s
	}
	if a := c.Query("adopted"); a != "" {
		b, err := strconv.ParseBool(a)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "adopted must be true or false"})
			return
		}
		filter["adopted"] = b
	}
	if q := c.Query("search"); q != "" {
		filter["name"] = bson.M{"$regex": regexp.QuoteMeta(q), "$options": "i"}
	}
	ageRange := bson.M{}
	for param, op := range map[string]string{"minAge": "$gte", "maxAge": "$lte"} {
		if v := c.Query(param); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil || n < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": param + " must be a non-negative integer"})
				return
			}
			ageRange[op] = n
		}
	}
	if len(ageRange) > 0 {
		filter["age"] = ageRange
	}

	// ---- Sorting ----
	sortField := "createdAt"
	if s := c.Query("sort"); s != "" {
		f, ok := sortableFields[s]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "sort must be one of: name, age, weight, species, createdAt"})
			return
		}
		sortField = f
	}
	sortDir := 1
	switch c.DefaultQuery("order", "asc") {
	case "asc":
	case "desc":
		sortDir = -1
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "order must be asc or desc"})
		return
	}

	// ---- Pagination ----
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be a positive integer"})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 100"})
		return
	}

	ctx := c.Request.Context()
	total, err := h.Animals.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	opts := options.Find().
		SetSort(bson.D{{Key: sortField, Value: sortDir}, {Key: "_id", Value: 1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cur, err := h.Animals.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	animals := []models.Animal{}
	if err := cur.All(ctx, &animals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       animals,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": int(math.Ceil(float64(total) / float64(limit))),
	})
}

func (h *AnimalHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var animal models.Animal
	err := h.Animals.FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&animal)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, gin.H{"error": "animal not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(http.StatusOK, animal)
}

func (h *AnimalHandler) Create(c *gin.Context) {
	var in models.AnimalInput
	if err := c.ShouldBindJSON(&in); err != nil {
		bindError(c, err)
		return
	}
	now := time.Now().UTC()
	animal := models.Animal{
		Name: in.Name, Species: in.Species, Breed: in.Breed, Age: *in.Age,
		Weight: in.Weight, Adopted: in.Adopted,
		CreatedBy: c.GetString("username"), CreatedAt: now, UpdatedAt: now,
	}
	res, err := h.Animals.InsertOne(c.Request.Context(), animal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	animal.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, animal)
}

// Update replaces an animal. Users can update only animals they created; admins can update any.
func (h *AnimalHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in models.AnimalInput
	if err := c.ShouldBindJSON(&in); err != nil {
		bindError(c, err)
		return
	}

	filter := bson.M{"_id": id}
	if c.GetString("role") != models.RoleAdmin {
		filter["createdBy"] = c.GetString("username")
	}
	update := bson.M{"$set": bson.M{
		"name": in.Name, "species": in.Species, "breed": in.Breed, "age": *in.Age,
		"weight": in.Weight, "adopted": in.Adopted, "updatedAt": time.Now().UTC(),
	}}

	var animal models.Animal
	err := h.Animals.FindOneAndUpdate(c.Request.Context(), filter, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&animal)
	if errors.Is(err, mongo.ErrNoDocuments) {
		h.notFoundOrForbidden(c, id)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(http.StatusOK, animal)
}

// Delete removes an animal (admin only, enforced by middleware).
func (h *AnimalHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	res, err := h.Animals.DeleteOne(c.Request.Context(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "animal not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AnimalHandler) notFoundOrForbidden(c *gin.Context, id primitive.ObjectID) {
	n, _ := h.Animals.CountDocuments(c.Request.Context(), bson.M{"_id": id})
	if n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "animal not found"})
		return
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "you can only modify animals you created"})
}

func parseID(c *gin.Context) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return id, false
	}
	return id, true
}
