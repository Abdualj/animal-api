package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"animal-api/internal/middleware"
	"animal-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Users     *mongo.Collection
	JWTSecret string
}

// Register creates a new account with the "user" role.
func (h *AuthHandler) Register(c *gin.Context) {
	var in models.Credentials
	if err := c.ShouldBindJSON(&in); err != nil {
		bindError(c, err)
		return
	}

	user, err := CreateUser(c.Request.Context(), h.Users, in.Username, in.Password, models.RoleUser)
	if mongo.IsDuplicateKeyError(err) {
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// Login checks credentials and returns a signed JWT.
func (h *AuthHandler) Login(c *gin.Context) {
	var in models.Credentials
	if err := c.ShouldBindJSON(&in); err != nil {
		bindError(c, err)
		return
	}

	var user models.User
	err := h.Users.FindOne(c.Request.Context(), bson.M{"username": in.Username}).Decode(&user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	expires := time.Now().Add(24 * time.Hour)
	claims := middleware.Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.Hex(),
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "expiresAt": expires, "user": user})
}

// CreateUser hashes the password and inserts a user document.
func CreateUser(ctx context.Context, users *mongo.Collection, username, password, role string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{Username: username, PasswordHash: string(hash), Role: role, CreatedAt: time.Now().UTC()}
	res, err := users.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

// EnsureAdmin seeds the admin account if it does not exist yet.
func EnsureAdmin(ctx context.Context, users *mongo.Collection, username, password string) error {
	err := users.FindOne(ctx, bson.M{"username": username}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		_, err = CreateUser(ctx, users, username, password, models.RoleAdmin)
	}
	return err
}
