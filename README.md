# Chirpy — Twitter-like backend (server)

Chirpy is a small Go backend that provides a Twitter-like API service: short messages are called "chirps" (equivalent to tweets). This repository contains the HTTP handlers, simple authentication (JWT + refresh tokens), a Polka webhook to upgrade users to "Chirpy Red", and Swagger annotations to generate interactive API docs.

**Requirements**
- Go 1.20+ (or compatible)
- PostgreSQL for the application database
- `swag` CLI (optional, for generating Swagger docs) — `github.com/swaggo/swag/cmd/swag`

**Environment variables**
- `DB_URL` — database URL (the code uses a local Postgres connection by default)
- `PLATFORM` — `dev` or `prod` (some admin endpoints are restricted to `dev`)
- `SECRET` — JWT secret used to sign tokens
- `POLKA_KEY` — API key expected by the Polka webhook

**Generate Swagger docs (optional)**
1. Install swag: `go install github.com/swaggo/swag/cmd/swag@latest`
2. From the repo root run: `swag init -g main.go`

This creates the `docs` package referenced by the server. If you don't generate the `docs` package, remove the `_ "github.com/tsironi93/WebServer/docs"` import in `main.go` or generate docs as shown.

**Run the server**
Set required environment variables and start the server:

```bash
export DB_URL="postgres://postgres:@localhost:5432/chirpy?sslmode=disable"
export PLATFORM=dev
export SECRET="your_jwt_secret"
export POLKA_KEY="your_polka_key"
go run .
```

The server listens on port `8080` by default. Swagger UI (when docs are generated) is available at:

`http://localhost:8080/swagger/index.html`

**Primary API resources**
- Chirps (short messages):
  - `GET /api/chirps` — list chirps (optional query params: `author_id`, `sort`)
  - `POST /api/chirps` — create a chirp (requires `Authorization: Bearer <jwt>`)
  - `GET /api/chirps/{chirpID}` — retrieve a single chirp
  - `DELETE /api/chirps/{chirpID}` — delete a chirp (owner only)

- Users & auth:
  - `POST /api/users` — create a user (`email`, `password`)
  - `PUT /api/users` — update the authenticated user's info (requires `Authorization: Bearer <jwt>`)
  - `POST /api/login` — authenticate and receive `token` (JWT) and `refresh_token`
  - `POST /api/refresh` — exchange a refresh token for a new JWT (send `Authorization: Bearer <refresh_token>`)
  - `POST /api/revoke` — revoke a refresh token (send `Authorization: Bearer <refresh_token>`)

- Webhooks and admin:
  - `POST /api/polka/webhooks` — Polka webhook to upgrade users (expects `Authorization: ApiKey <key>`)
  
    Note: Polka in this project is a fictional payment platform used to illustrate webhook-driven upgrades. When Polka sends a `user.upgraded` event the webhook handler verifies the `ApiKey` and upgrades the user to the "Chirpy Red" tier.
  - `GET /api/healthz` — readiness check
  - `GET /admin/metrics` — simple HTML admin metrics
  - `POST /admin/reset` — dev-only reset (clears hits counter and deletes all users)

**Quick examples**
Create a user:

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"me@example.com","password":"secret"}'
```

Login and get tokens:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"me@example.com","password":"secret"}'
```

Create a chirp (use returned JWT):

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT>" \
  -d '{"body":"Hello from Chirpy!"}'
```
