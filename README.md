# Animal API

REST API eläinsuojan eläinten hallintaan. Projekti on toteutettu **Go-kielellä (Gin)** ja tietokantana toimii **MongoDB**.

**Live demo:** https://animal-api-fyfq.onrender.com/docs

Palvelu on julkaistu Renderin ilmaisella palvelutasolla. Ensimmäinen pyyntö voi kestää noin minuutin, koska palvelu menee käyttämättömänä lepotilaan.

## Palautus ja projektin linkit

| Palautettava osa | Linkki |
|---|---|
| Lähdekoodi | https://github.com/Abdualj/animal-api |
| Julkaistu API | https://animal-api-fyfq.onrender.com |
| Swagger-dokumentaatio | https://animal-api-fyfq.onrender.com/docs |

**Yhteenveto:** Kaikki viisi vaihetta (Step 1–5) on suoritettu Go-toteutuksena.

---

## Step 1 – REST API -tutoriaali 

Ensimmäisessä vaiheessa suoritin TutorialEdgen Go REST API -tutoriaalin:

[Creating a RESTful API With Go](https://tutorialedge.net/golang/creating-restful-api-with-golang/)

Tutoriaalissa toteutin:
- HTTP-palvelimen Go-kielellä.
- Article-rakenteen ja JSON-datan käsittelyn.
- Reitityksen gorilla/mux-kirjastolla.
- GET-, POST- ja DELETE-pyynnöt.
- PUT-pyynnön, jolla olemassa olevaa artikkelia voidaan päivittää.

Lisäksi toteutin tutoriaalin lopussa olevan Challenge-tehtävän. Jos päivitettävää artikkelia ei löydy, API palauttaa 404-virheen.

**Koodi löytyy:** `tutorial/main.go`

Tutoriaalin jälkeen siirryin toteuttamaan varsinaista Animal API -projektia Gin-kehyksellä ja MongoDB-tietokannalla.

## Step 2 – API ja CRUD-toiminnot 

Toteutin eläinsuojan eläinten hallintaan tarkoitetun REST API:n.

Teknologiat:
- Go 1.23
- Gin
- MongoDB
- MongoDB Go Driver

Jokaisella eläimellä on seuraavat tiedot:
- Nimi
- Laji
- Rotu
- Ikä
- Paino
- Adoptiotilanne
- Luontiaika ja päivitysaika
- Eläimen luoneen käyttäjän tiedot

API sisältää kaikki CRUD-toiminnot:

| HTTP-metodi | Osoite | Toiminto |
|---|---|---|
| GET | `/api/animals` | Hakee kaikki eläimet |
| GET | `/api/animals/{id}` | Hakee yhden eläimen |
| POST | `/api/animals` | Lisää uuden eläimen |
| PUT | `/api/animals/{id}` | Päivittää eläimen tiedot |
| DELETE | `/api/animals/{id}` | Poistaa eläimen |

**Koodi:** `internal/handlers/animals.go`

## Step 3 – Edistyneet ominaisuudet 

Toteutin kaikki kuusi edistynyttä ominaisuutta, vaikka tehtävän vähimmäisvaatimus oli kolme.

**1. Filtering – Suodatus**

Eläimiä voi suodattaa esimerkiksi lajin, adoptiotilanteen ja iän perusteella.

Myös eläimen nimeä voi hakea. Haku ei erottele isoja ja pieniä kirjaimia.

Esimerkki:
`/api/animals?species=dog&adopted=false&minAge=1&maxAge=10&search=rex`

**2. Sorting – Lajittelu**

Eläimiä voi järjestää esimerkiksi nimen, iän, painon tai lajin mukaan.

Järjestys voi olla nouseva (asc) tai laskeva (desc).

Esimerkki:
`?sort=age&order=desc`

**3. Pagination – Sivutus**

Eläinlistaus on jaettu sivuihin, jotta suuria tietomääriä on helpompi käsitellä.

Käyttäjä voi valita sivun ja sivulla näytettävien eläinten määrän. Enimmäismäärä on 100 eläintä sivua kohden.

Esimerkki:
`?page=2&limit=10`

Vastaus sisältää myös kokonaismäärän ja sivujen lukumäärän.

**4. Authentication – Käyttäjän tunnistautuminen**

API:ssa on käyttäjärekisteröinti ja kirjautuminen.

- Käyttäjä voi luoda tunnuksen.
- Kirjautumisen jälkeen API palauttaa JWT-tokenin.
- Token on voimassa 24 tuntia.
- Salasanat tallennetaan turvallisesti bcrypt-tiivisteinä.

Token lähetetään API-pyyntöjen mukana Authorization-headerissa.

**5. Authorization – Käyttöoikeudet**

API:ssa on kaksi käyttäjäroolia: `user` ja `admin`.

Käyttöoikeudet toimivat seuraavasti:

- Kaikki voivat tarkastella eläimiä.
- Kirjautunut käyttäjä voi lisätä eläimiä.
- Käyttäjä voi muokata vain itse lisäämiään eläimiä.
- Admin voi muokata ja poistaa kaikkia eläimiä.

Jos käyttäjällä ei ole tarvittavia oikeuksia, API palauttaa 403-virheen.

Admin-tunnus luodaan automaattisesti palvelimen käynnistyessä ympäristömuuttujien avulla.

**6. Validation – Datan validointi**

API tarkistaa käyttäjän lähettämän datan ennen sen tallentamista.

Esimerkkejä tarkistuksista:
- Nimen pituus on 2–50 merkkiä.
- Ikä on välillä 0–100.
- Painon täytyy olla suurempi kuin nolla.
- Käyttäjänimen tulee sisältää vain kirjaimia ja numeroita.
- Salasanan tulee olla vähintään kuusi merkkiä pitkä.
- Eläimen lajin tulee kuulua sallittuihin arvoihin.

Virheellisestä datasta palautetaan 400-virhe sekä tarkemmat tiedot siitä, mikä kenttä ei läpäissyt validointia.

**Keskeiset kooditiedostot:**
- `internal/middleware/auth.go`
- `internal/handlers/auth.go`
- `internal/models/models.go`

## Step 4 – Dokumentaatio 

Dokumentoin API:n OpenAPI 3.0 -standardin mukaisesti.

Dokumentaatio sisältää:
- Kaikki API-endpointit.
- Käytettävissä olevat parametrit.
- Pyyntöjen ja vastausten rakenteet.
- Mahdolliset virheilmoitukset.
- JWT-tunnistautumisen ohjeet.

API:n dokumentaatiota voi tarkastella ja testata suoraan selaimessa Swagger UI:n avulla.

**Swagger UI:** https://animal-api-fyfq.onrender.com/docs

Swaggerissa voi esimerkiksi kirjautua sisään, kopioida JWT-tokenin ja käyttää Authorize-toimintoa suojattujen endpointien testaamiseen.

**Dokumentaation tiedosto:** `docs/openapi.yaml`

## Step 5 – Docker ja julkaisu pilveen 

Projektin voi käynnistää Dockerin avulla yhdellä komennolla.

Toteutin:
- Dockerfilen, jossa käytetään kaksivaiheista build-prosessia.
- Pienemmän Alpine-pohjaisen tuotantoimagen.
- Kontin, joka toimii ilman root-oikeuksia.
- docker-compose.yml-tiedoston, joka käynnistää sekä API:n että MongoDB:n.

**Pilvijulkaisu**

API on julkaistu Render-palvelussa Docker-konttina. Tietokantana toimii MongoDB Atlas.

Julkaisussa käytetyt palvelut:
- Render: API:n hosting.
- MongoDB Atlas: Pilvitietokanta.
- GitHub: Lähdekoodin hallinta ja julkaisun lähde.

Salaisuudet, kuten tietokantayhteys ja JWT-avain, määritellään ympäristömuuttujina. Niitä ei tallenneta GitHub-repositorioon.

Valitsin Renderin Azuren sijaan, koska Render tarjoaa ilmaisen palvelutason ilman luottokorttivaatimusta.

**Huomio:** Renderin ilmainen palvelu menee lepotilaan 15 minuutin käyttämättömyyden jälkeen. Tämän vuoksi ensimmäinen pyyntö voi kestää noin 30–60 sekuntia.

---

## Testaus

Testasin API:n endpointit sekä paikallisesti että pilvessä curl-komentojen avulla.

Testauksessa käytiin läpi:

- Eläinten luominen, hakeminen, päivittäminen ja poistaminen.
- Suodatuksen, lajittelun ja sivutuksen toiminta.
- Kirjautuminen ja JWT-tokenin käyttö.
- 401-virhe, kun tunnistautuminen puuttuu.
- 403-virhe, kun käyttäjällä ei ole tarvittavia oikeuksia.
- 400-virhe, kun lähetetty data on virheellistä.
- 404-virhe, kun haettua eläintä ei löydy.

---

## Projektin käynnistäminen paikallisesti

Projektin voi käynnistää Docker Composella:

```bash
docker compose up -d --build
```

Käynnistyksen jälkeen palvelut löytyvät seuraavista osoitteista:

| Palvelu | Osoite |
|---|---|
| API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/docs |

Oletusarvoinen admin-tunnus paikallisessa ympäristössä:

- Käyttäjätunnus: `admin`
- Salasana: `admin123`

Tunnukset voi vaihtaa ympäristömuuttujilla.

## API-endpointit

| Metodi | Osoite | Käyttöoikeus | Tarkoitus |
|---|---|---|---|
| GET | `/health` | Ei vaadita | Palvelun tilan tarkistus |
| POST | `/api/auth/register` | Ei vaadita | Rekisteröityminen |
| POST | `/api/auth/login` | Ei vaadita | Kirjautuminen |
| GET | `/api/animals` | Ei vaadita | Eläinten listaus |
| GET | `/api/animals/{id}` | Ei vaadita | Yksittäisen eläimen haku |
| POST | `/api/animals` | User / Admin | Eläimen lisääminen |
| PUT | `/api/animals/{id}` | Omistaja / Admin | Eläimen päivittäminen |
| DELETE | `/api/animals/{id}` | Admin | Eläimen poistaminen |

## Projektin rakenne

```text
main.go                      Palvelimen käynnistys ja reitit
tutorial/                    Tutoriaalin toteutus
internal/config              Ympäristömuuttujien asetukset
internal/db                  MongoDB-yhteys
internal/models              Tietomallit ja validointi
internal/middleware/auth.go  Tunnistautuminen ja käyttöoikeudet
internal/handlers            API:n käsittelylogiikka
docs/openapi.yaml            API-dokumentaatio
```

---

**Lopputulos:** Kaikki viisi projektin vaihetta on toteutettu. API sisältää CRUD-toiminnot, kuusi edistynyttä ominaisuutta, dokumentaation, testauksen sekä toimivan pilvijulkaisun.

Käytin apuna tekoälyä readme.md tekstissä sekä lähdekoodissa. Tarkistin ja omaan kuitenkin koodin ymmärtämistä