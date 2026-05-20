# Setup

## Local Development

Install Docker, Go 1.22+, and Node 20+.

```bash
cd infrasphere
docker compose up --build
```

Frontend: `http://localhost:5173`  
Backend: `http://localhost:8080/healthz`

## PostgreSQL

Docker Compose starts PostgreSQL 16 and loads `backend/migrations/001_init.sql`. The MVP API currently uses an in-memory seeded store while the schema establishes the production persistence model.

## Redis

Redis is available at `redis://localhost:6379` for caching, rate limiting, queue metadata, and session-related extensions.

## Environment Variables

Use:

```bash
APP_ENV=development
APP_PORT=8080
DATABASE_URL=postgres://infrasphere:infrasphere@localhost:5432/infrasphere?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=change-me
ENCRYPTION_KEY=change-me-32-byte-key
CORS_ALLOWED_ORIGINS=http://localhost:5173
```

AI, OIDC, and provider credentials are configured with the variables listed in `docs/HOW_TO_USE_AND_SETUP.md`.

## Running Backend

```bash
cd backend
go run ./cmd/api
```

## Running Frontend

```bash
cd frontend
npm install
npm run dev
```

## Running Workers

```bash
cd backend
go run ./cmd/worker
```

## First Admin User

The MVP seeds a default in-memory admin:

```text
admin@infrasphere.local / ChangeMe123!
```

In production, create the first admin through a bootstrap CLI or one-time migration that stores a bcrypt/Argon2 password hash.

