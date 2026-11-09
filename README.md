# Animal API

REST API for managing shelter animals, written in **Go (Gin)** with **MongoDB**.

**Live demo:** https://animal-api-fyfq.onrender.com/docs (hosted on Render free tier + MongoDB Atlas; the first request may take ~1 min to wake up)

## Palautus – suoritetut stepit

| | Linkki |
|---|---|
| Lähdekoodi | https://github.com/Abdualj/animal-api |
| API pilvessä | https://animal-api-fyfq.onrender.com |
| Swagger-dokumentaatio | https://animal-api-fyfq.onrender.com/docs |

> Render free -palvelu nukahtaa 15 min käyttämättömyyden jälkeen, joten ensimmäinen pyyntö voi kestää noin minuutin.

**Yhteenveto:** Stepit 1–5 suoritettu (Go-toteutus).

### Step 1 – Tutoriaali ✅
Tehty **Go REST API Tutorial 1**: TutorialEdge, *Creating a RESTful API With Go* (https://tutorialedge.net/golang/creating-restful-api-with-golang/).

- Koodi on kansiossa [`tutorial/`](tutorial/main.go).
- Toteutettu kaikki tutoriaalin vaiheet: `net/http`-palvelin, `Article`-struct ja JSON-marshalling, `gorilla/mux`-reititin ja path-muuttujat sekä GET-, POST- ja DELETE-endpointit.
- Lisäksi tutoriaalin lopun **Challenge** eli `PUT /article/{id}` (päivitys, palauttaa 404 jos artikkelia ei löydy).
- Ajo: `cd tutorial && go run main.go`, jonka jälkeen API vastaa osoitteessa http://localhost:10000/articles.

Tutoriaalin pohjalta varsinainen API (stepit 2–5) tehtiin Gin-kehyksellä ja MongoDB:llä.

### Step 2 – API ja CRUD ✅
- Toteutettu **Go 1.23 + Gin** -kehyksellä, ja tiedot tallennetaan **MongoDB**:hen (virallinen `mongo-driver`).
- Aiheena on eläinsuojan eläimet (Animal API). Eläimellä on kentät `name`, `species`, `breed`, `age`, `weight`, `adopted` sekä automaattiset `createdBy`, `createdAt` ja `updatedAt`.
- CRUD-operaatiot:

| Metodi | Polku | Toiminto |
|---|---|---|
| GET | `/api/animals` | Listaa eläimet |
| GET | `/api/animals/{id}` | Hakee yhden eläimen |
| POST | `/api/animals` | Luo eläimen |
| PUT | `/api/animals/{id}` | Päivittää eläimen |
| DELETE | `/api/animals/{id}` | Poistaa eläimen |

Koodi: [`internal/handlers/animals.go`](internal/handlers/animals.go)

### Step 3 – Edistyneet ominaisuudet ✅ (kaikki 6, vaatimus oli 3)
1. **Filtering:** `GET /api/animals?species=dog&adopted=false&minAge=1&maxAge=10&search=rex`. Haku nimellä (`search`) ei erottele isoja ja pieniä kirjaimia.
2. **Sorting:** `?sort=name|age|weight|species|createdAt&order=asc|desc`.
3. **Pagination:** `?page=2&limit=10` (enintään 100 per sivu). Vastaus sisältää `data`, `page`, `limit`, `total` ja `totalPages`.
4. **Authentication:** käyttäjä rekisteröityy (`POST /api/auth/register`) ja kirjautuu (`POST /api/auth/login`), jolloin hän saa **JWT-tokenin** (voimassa 24 h). Salasanat tallennetaan **bcrypt**-tiivisteinä. Token annetaan headerissa `Authorization: Bearer <token>`.
5. **Authorization (roolit `user` ja `admin`):**
   - Kuka tahansa voi lukea eläimiä.
   - Kirjautunut käyttäjä voi luoda eläimiä ja muokata **vain itse luomiaan** eläimiä (muuten vastaus on `403`).
   - Vain **admin** voi poistaa eläimiä ja muokata kaikkia eläimiä.
   - Admin-tunnus luodaan automaattisesti käynnistyksessä ympäristömuuttujista `ADMIN_USERNAME` ja `ADMIN_PASSWORD`.
6. **Validation:** pyyntöjen rungot tarkistetaan (esim. nimi 2–50 merkkiä, laji sallittujen arvojen listasta, ikä 0–100, paino > 0, käyttäjänimi vain kirjaimia ja numeroita, salasana vähintään 6 merkkiä). Virheestä palautuu `400` ja kenttäkohtaiset virheilmoitukset. Myös query-parametrit ja id:n muoto tarkistetaan.

Koodi: [`internal/middleware/auth.go`](internal/middleware/auth.go), [`internal/handlers/auth.go`](internal/handlers/auth.go), [`internal/models/models.go`](internal/models/models.go)

### Step 4 – Dokumentaatio ✅
- API on dokumentoitu **OpenAPI 3.0** -muodossa: [`docs/openapi.yaml`](docs/openapi.yaml). Mukana ovat kaikki endpointit, parametrit, skeemat, virhevastaukset ja JWT-autentikointi.
- API tarjoilee **Swagger UI**:n osoitteessa `/docs`. Siellä endpointeja voi testata suoraan selaimessa: kirjaudu sisään, kopioi token ja paina **Authorize**.

### Step 5 – Docker ja pilvideployment ✅
- **Dockerfile** on kaksivaiheinen (multi-stage): Go-binääri käännetään `golang`-imagessa ja ajetaan pienessä `alpine`-imagessa ilman root-oikeuksia.
- **docker-compose.yml** käynnistää paikallisesti API:n ja MongoDB:n yhdellä komennolla.
- **Pilvi:** API pyörii Docker-konttina **Renderissä** (ilmainen taso). Asetukset ovat tiedostossa [`render.yaml`](render.yaml), ja Render rakentaa imagen GitHub-repon Dockerfilestä. Tietokantana on **MongoDB Atlas** (ilmainen klusteri). Salaisuudet (`MONGO_URI`, `JWT_SECRET`, `ADMIN_PASSWORD`) annetaan ympäristömuuttujina, eivätkä ne ole repossa.
- Azuren sijaan käytin Renderiä, koska se on ilmainen eikä vaadi luottokorttia.

### Testaus
Kaikki endpointit on testattu sekä paikallisesti että Renderissä curlilla. Testatut tapaukset: CRUD, suodatus, lajittelu, sivutus, `401` ilman tokenia, `403` väärällä roolilla tai toisen käyttäjän eläimellä, `400` virheellisellä datalla ja `404` olemattomalla id:llä.

---

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
tutorial/                    Step 1: TutorialEdge REST API tutorial (gorilla/mux)
internal/config              environment config
internal/db                  Mongo connection + unique index
internal/models              Animal / User structs + validation rules
internal/middleware/auth.go  JWT authentication + role authorization
internal/handlers            animal CRUD/listing, auth, error formatting
docs/openapi.yaml            OpenAPI 3 specification
```
