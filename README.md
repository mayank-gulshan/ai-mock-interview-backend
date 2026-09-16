# AI Mock Interview Backend — Go port

A 1:1 behavioral port of the original Kotlin/Spring Boot backend
(`Ai_Mock_Interview_Backend`) to Go, using Gin + GORM.

## Original architecture (summary)

- **Framework:** Spring Boot 4.1.0, Kotlin 2.3.21, Java 21 toolchain.
- **Web layer:** `spring-boot-starter-webmvc` (Servlet MVC), two
  `@RestController`s: `AuthController` (`/api/auth/**`) and
  `InterviewController` (`/api/interview/**`).
- **Persistence:** Spring Data JPA / Hibernate over PostgreSQL (hosted on
  Supabase's connection pooler), `ddl-auto=update` (auto schema migration).
  Three `@Entity` classes: `User`, `Session`, `Answer`. Only `User` is
  actually read/written by any service today — `Session`/`Answer` and their
  repositories, plus the `InterviewDto`/`HistoryDto` request/response types,
  exist in the repo but aren't wired to any controller yet (a half-built
  "save session history" feature). This port preserves that exact state.
- **Auth:** Custom JWT (jjwt 0.13.0), **not** Spring Security's own
  UserDetailsService flow beyond satisfying the `UserDetails` interface. A
  `OncePerRequestFilter` (`JwtAuthFilter`) parses a Bearer access token,
  loads the user, and puts it in the `SecurityContext`; `SecurityConfig`
  permits `/api/auth/**` and requires authentication on everything else. A
  `CustomAuthEntryPoint` returns a JSON 401 for unauthenticated access.
  Passwords are hashed with `BCryptPasswordEncoder`.
- **Third-party integration:** Google Gemini (`gemini-2.5-flash`)
  `generateContent` REST API, called via `RestTemplate`, for two things:
  generating 10 interview Q&As with a strict JSON response schema, and
  evaluating a candidate's answer against a hand-written prompt (Gemini's
  raw text response is then parsed as loose JSON).
- **Errors:** A `@RestControllerAdvice` (`GlobalExceptionHandler`) maps three
  custom exceptions to specific HTTP statuses (409 for "user exists", 401 for
  bad credentials/bad tokens) and Spring Security's own
  `AuthenticationException` to 401.

## Go stack

| Concern              | Kotlin/Spring                     | Go                                      |
|-----------------------|-----------------------------------|------------------------------------------|
| Web framework         | Spring MVC                        | **Gin** — minimal, fast, huge ecosystem, closest match to Spring MVC's controller/middleware model |
| ORM / DB              | Spring Data JPA + Hibernate       | **GORM** + `gorm.io/driver/postgres`, `AutoMigrate` in place of `ddl-auto=update` |
| Auth / JWT             | jjwt + Spring Security filter     | **golang-jwt/jwt/v5** + a Gin middleware (`internal/security`) |
| Password hashing       | `BCryptPasswordEncoder`           | `golang.org/x/crypto/bcrypt` |
| Config                 | `application.properties` + `${ENV_VAR}` | plain environment variables via `internal/config` (same variable names: `DB_PASS`, `JWT_SECRET_64`, `GEMINI_API_KEY`, `PORT`) |
| HTTP client (Gemini)   | `RestTemplate`                    | `net/http.Client` |

Gin was chosen over Echo mainly for its wider adoption/ecosystem and
middleware model that maps cleanly onto Spring's filter chain concept used
here (`JwtAuthFilter` → `security.RequireAuth` middleware).

## Project layout

```
ai-mock-interview-backend/
├── go.mod
├── .env.example
├── README.md
├── cmd/
│   └── server/
│       └── main.go              # wiring + startup (AiMockInterviewApplication.kt)
└── internal/
    ├── config/config.go         # env var loading (application.properties)
    ├── db/db.go                 # GORM connection + logger (ddl-auto=update)
    ├── models/                  # Entity/*.kt
    │   ├── user.go
    │   ├── session.go
    │   └── answer.go
    ├── dto/                     # Dto/*.kt
    │   ├── auth_dto.go
    │   ├── interview_dto.go
    │   ├── history_dto.go
    │   ├── answer_dto.go
    │   └── gemini_dto.go
    ├── apperrors/errors.go      # Security/Exception.kt
    ├── repository/              # Repository/*.kt
    │   ├── user_repository.go
    │   ├── session_repository.go
    │   └── answer_repository.go
    ├── security/                # Security/*.kt (minus Auth*, Exception.kt)
    │   ├── hash.go              # HashEncoder.kt
    │   ├── jwt.go                # JWTService.kt
    │   └── middleware.go        # JWTAuthFilter.kt + SecurityConfug.kt + CustonAuthEntryPoint.kt
    ├── auth/                    # Security/AuthController.kt + AuthService.kt
    │   ├── service.go
    │   └── handler.go
    ├── gemini/service.go        # Rest/GeminiService.kt + AppConfig.kt
    ├── interview/               # Rest/GeminiController.kt + Security/EvalutionService.kt
    │   ├── evaluation_service.go
    │   └── handler.go
    └── server/router.go         # route table + CORS/logging/recovery middleware
```

## API (unchanged from the original)

| Method | Path                     | Auth required | Notes |
|--------|--------------------------|----------------|-------|
| POST   | `/api/auth/register`     | no             | 200 on success, 409 if email taken |
| POST   | `/api/auth/login`        | no             | 200 on success, 401 on bad credentials |
| POST   | `/api/auth/refresh`      | no             | 200 on success, 401 on invalid/expired refresh token |
| GET    | `/api/interview/questions?role=&difficulty=medium` | yes | `role` required, `difficulty` defaults to `medium` |
| POST   | `/api/interview/evaluate`| yes            | Body: `EvaluationRequest` |

Auth uses `Authorization: Bearer <accessToken>`. Access tokens last 15
minutes, refresh tokens 30 days (hardcoded, same as the original — the
`app.jwt.expiration` property in the source repo is defined but never
actually read).

## Notable behavioral / porting notes

- **CORS**: the original `SecurityConfig` never enables Spring Security's
  CORS support, so as written it has no CORS headers at all. This port adds
  a permissive `Access-Control-Allow-Origin: *` middleware
  (`internal/server/router.go`) so a browser frontend can actually call it —
  tighten this to your real origin(s) before deploying.
- **`Session`/`Answer` entities and repositories, plus `InterviewDto`/
  `HistoryDto`**, are translated for parity but unused, exactly matching the
  original (see NOTE comments in `internal/models/session.go`).
- **Validation**: `RegisterRequest` imports `jakarta.validation` annotations
  in the original but never applies them (no `@field:` targets, no `@Valid`
  on the controller param) — no validation is enforced, and this port
  matches that (no server-side validation on register/login fields beyond
  JSON well-formedness).
- **`GET /api/interview/questions`** requires `role` as a plain
  (non-nullable, no-default) query param in Spring, which returns 400 if
  missing; this port returns the same 400 for a missing `role`.

## Running locally

1. Install Go 1.22+.
2. `cp .env.example .env` and fill in `DB_PASS`, `JWT_SECRET_64`,
   `GEMINI_API_KEY` (and optionally override the DB host/user or set
   `DATABASE_URL` directly).
3. Fetch dependencies (requires network access to the Go module proxy,
   which wasn't reachable from the sandbox this was written in — run this
   yourself once):
   ```bash
   go mod tidy
   ```
4. Run it:
   ```bash
   go run ./cmd/server
   ```
   The server listens on `0.0.0.0:$PORT` (default `8080`), matching the
   original `server.address=0.0.0.0` / `server.port=${PORT:8080}`.

## Building a binary

```bash
go build -o bin/server ./cmd/server
./bin/server
```
