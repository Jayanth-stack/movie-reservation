# Movie Reservation API

A backend service for a movie reservation system: users browse movies and showtimes, reserve seats, and admins manage the catalog.

> Status: **early scaffolding**. DB connection + `User` model are in place; auth, movies, showtimes, and reservations are stubs.

## Stack

- **Language:** Go 1.25
- **HTTP:** [Gin](https://github.com/gin-gonic/gin)
- **ORM / DB:** [GORM](https://gorm.io) + PostgreSQL 16 (pgx driver)
- **Auth:** JWT ([`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt))
- **Config:** [`godotenv`](https://github.com/joho/godotenv)
- **Testing:** [`go.uber.org/mock`](https://github.com/uber-go/mock)
- **Infra (dev):** Docker Compose (Postgres only)

## Project layout

```
.
├── cmd/
│   └── api/                # entrypoint (main.go)
├── internal/
│   ├── auth/               # registration, login, JWT issuing  (TODO)
│   ├── config/             # typed config loader               (TODO)
│   ├── db/                 # gorm.DB bootstrap
│   ├── middleware/         # auth / role guards                (TODO)
│   ├── movie/              # movie domain                      (TODO)
│   ├── reservation/        # reservation / seats               (TODO)
│   ├── shared/             # cross-cutting helpers             (TODO)
│   ├── showtime/           # showtime domain                   (TODO)
│   └── user/               # User model (+ role)
├── docker-compose.yml      # postgres:16
├── go.mod / go.sum
└── .env                    # local secrets (gitignored)
```

Conventions:

- Each domain package owns its `model.go`, `repository.go`, `service.go`, `handler.go`, and `routes.go`.
- HTTP layer lives only in `handler.go` / `routes.go`. Business logic in `service.go`. SQL only in `repository.go`.
- `internal/shared` is for things genuinely shared (errors, pagination, response helpers) — not a dumping ground.

## Getting started

### 1. Prerequisites

- Go **1.25+**
- Docker + Docker Compose
- (optional) `make`

### 2. Configure environment

Create a `.env` in the repo root:

```env
DB_HOST=localhost
DB_USER=movie_user
DB_PASSWORD=movie_pass
DB_NAME=movie_db
DB_PORT=5432
JWT_SECRET=change-me
```

> The Postgres container reads `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` from the same file — keep the keys in `.env` aligned with `docker-compose.yml` (see TODO below).

### 3. Start the database

```bash
docker compose up -d
```

### 4. Run the API

```bash
go run ./cmd/api
```

Server listens on **`:8080`**.

Sanity check:

```bash
curl http://localhost:8080/
# {"message":"Welcome to the Movie Reservation API"}
```

## Domain model (planned)

| Entity        | Purpose                                                 |
| ------------- | ------------------------------------------------------- |
| `User`        | Authenticated principal. Roles: `admin`, `user`.        |
| `Movie`       | Title, description, poster, genre, duration.            |
| `Showtime`    | Movie × theater × start time. Owns a seat map.          |
| `Seat`        | Per-showtime seat with status (`free`, `held`, `sold`). |
| `Reservation` | A user's hold/purchase of one or more seats.            |

## API surface (planned)

```
POST   /auth/register
POST   /auth/login

GET    /movies
GET    /movies/:id
POST   /movies                 (admin)
PATCH  /movies/:id             (admin)
DELETE /movies/:id             (admin)

GET    /movies/:id/showtimes
POST   /showtimes              (admin)

GET    /showtimes/:id/seats
POST   /reservations           (auth)
GET    /reservations/me        (auth)
DELETE /reservations/:id       (auth)
```

## Development

### Useful commands

```bash
go run ./cmd/api          # run
go build ./...            # compile everything
go test ./...             # run tests
go vet ./...              # static checks
gofmt -w .                # format
```

### Migrations

Currently none. The plan is to call `db.AutoMigrate(&user.User{}, &movie.Movie{}, ...)` from `cmd/api/main.go` during bootstrap until/unless we adopt a real migration tool ([`golang-migrate`](https://github.com/golang-migrate/migrate) or [`goose`](https://github.com/pressly/goose)).

## Known issues / TODO

- [ ] `internal/user/model.go`: `Role` default tag is `'USER'` (uppercase) but constants are `"admin"`/`"user"` — won't match.
- [ ] `internal/db/db.go`: `TimeZone=Asia/Shanghai` is hardcoded — move to env or use UTC.
- [ ] `internal/db/db.go`: replace the global `var DB` with a returned `*gorm.DB` injected into services (better for tests).
- [ ] No `AutoMigrate` call yet — `users` table doesn't exist on a fresh DB.
- [ ] `docker-compose.yml` doesn't set `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` — the container will fall back to defaults and ignore `.env`.
- [ ] Add healthcheck + an `app` service to compose.
- [ ] Add `Makefile` and CI.
- [ ] Add tests.

## License

TBD.
