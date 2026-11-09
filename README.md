# Animal API

REST API for managing shelter animals, written in **Go (Gin)** with **MongoDB**.

**Live demo:** https://animal-api-fyfq.onrender.com/docs (hosted on Render free tier + MongoDB Atlas; the first request may take ~1 min to wake up)

| Assignment step | What was done |
|---|---|
| 2. CRUD | `GET` / `POST` / `PUT` / `DELETE` on `/api/animals` |
| 3. Advanced features (6/6) | Filtering, Sorting, Pagination, JWT Authentication, Role-based Authorization, Validation |
| 4. Documentation | OpenAPI 3 spec in [`docs/openapi.yaml`](docs/openapi.yaml), Swagger UI served at `/docs` |
| 5. Deployment | Multi-stage `Dockerfile` + `docker-compose.yml`; free deployment on Render (`render.yaml`) |

## Run locally

```bash
docker compose up -d --build
```

- API: http://localhost:8080
- Swagger UI: http://localhost:8080/docs
- Default admin: `admin` / `admin123` (change it with `ADMIN_USERNAME` / `ADMIN_PASSWORD`)

## Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/health` | – | Health check |
| POST | `/api/auth/register` | – | Create a user account (role `user`) |
| POST | `/api/auth/login` | – | Returns a JWT (valid for 24h) |
| GET | `/api/animals` | – | List with filters / sort / pagination |
| GET | `/api/animals/{id}` | – | Get one animal |
| POST | `/api/animals` | user, admin | Create an animal |
| PUT | `/api/animals/{id}` | owner or admin | Update an animal |
| DELETE | `/api/animals/{id}` | admin | Delete an animal |

### Query parameters for `GET /api/animals`

| Param | Example | Feature |
|---|---|---|
| `species` | `dog` | filtering |
| `adopted` | `false` | filtering |
| `minAge`, `maxAge` | `1`, `10` | filtering (range) |
| `search` | `rex` | filtering (case-insensitive name search) |
| `sort` | `name`, `age`, `weight`, `species`, `createdAt` | sorting |
| `order` | `asc` / `desc` | sorting |
| `page`, `limit` | `2`, `10` (max 100) | pagination |

Response:
```json
{ "data": [ ... ], "page": 1, "limit": 10, "total": 42, "totalPages": 5 }
```

### Authorization rules
- Anyone can read animals.
- A logged-in **user** can create animals and update only the animals they created (other animals return `403`).
- An **admin** can update and delete any animal.

### Validation
Request bodies are validated, and errors come back as `400` responses with details for each field:
```json
{ "error": "validation failed",
  "details": { "name": "must be at least 2 characters", "weight": "must be > 0" } }
```

## Example requests

```bash
# login
TOKEN=$(curl -s -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}' \
  localhost:8080/api/auth/login | sed -E 's/.*"token":"([^"]+)".*/\1/')

# create
curl -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"Rex","species":"dog","breed":"Husky","age":3,"weight":25}' \
  localhost:8080/api/animals

# list dogs, oldest first, page 1
curl "localhost:8080/api/animals?species=dog&sort=age&order=desc&page=1&limit=5"
```

## Configuration (environment variables)

| Variable | Default |
|---|---|
| `PORT` | `8080` |
| `MONGO_URI` | `mongodb://localhost:27017` |
| `MONGO_DB` | `animaldb` |
| `JWT_SECRET` | `change-me-in-production` |
| `ADMIN_USERNAME` / `ADMIN_PASSWORD` | `admin` / `admin123` |

## Deploy for free (Render + MongoDB Atlas)

Both services have free tiers and neither needs a credit card.

**1. Database: MongoDB Atlas (free M0 cluster)**
1. Sign up at https://www.mongodb.com/cloud/atlas/register and create a **free M0** cluster.
2. *Database Access*: add a database user and password.
3. *Network Access*: add `0.0.0.0/0`. Render's free tier has no fixed IP, so this is required.
4. *Connect → Drivers*: copy the connection string, e.g.
   `mongodb+srv://<user>:<password>@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority`

**2. API: Render (free web service, built from the Dockerfile)**
1. Push this project to a GitHub repository.
2. On https://render.com, sign in with GitHub and choose **New → Blueprint**, then select the repo. Render reads [`render.yaml`](render.yaml).
3. When asked, enter `MONGO_URI` (the Atlas string) and `ADMIN_PASSWORD`. `JWT_SECRET` is generated automatically.
4. Wait for the build to finish, then open `https://animal-api-xxxx.onrender.com/docs`.

> Free-tier note: the service sleeps after 15 minutes without traffic, so the first request after that takes about 30–60 s to wake it up.

## Project structure

```
main.go                      routes, Swagger UI, startup
internal/config              environment config
internal/db                  Mongo connection + unique index
internal/models              Animal / User structs + validation rules
internal/middleware/auth.go  JWT authentication + role authorization
internal/handlers            animal CRUD/listing, auth, error formatting
docs/openapi.yaml            OpenAPI 3 specification
```
