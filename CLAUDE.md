# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**Manutenção deste arquivo:** sempre que um trabalho nesta sessão alterar algo que uma seção abaixo descreve (arquitetura, comandos, comportamento das edições 4e/5e, camadas, deploy, etc.), atualize a seção correspondente antes de considerar a tarefa concluída — não crie um changelog separado; o histórico de mudanças já vive no `git log`.

## Project overview

RPG Manager is a character sheet manager for tabletop D&D that supports **two editions side by side: 4e and 3.5/5e**. A character's `Edition` field ("4e" or "5e") drives which rules, fields, and calculations apply throughout the backend — this is the central architectural fact of the codebase. There are three independent apps in this repo:

- `backend/` — Go + Gin + GORM REST API (the core app, port 8080)
- `frontend/` — React 19 + TypeScript + Vite SPA (port 5173)
- `ai-service/` — Python FastAPI microservice (port 8000) that RAG-queries indexed D&D rulebook PDFs via ChromaDB + Ollama to auto-generate class powers/skills

Comments, error messages, and commit messages in this codebase are in Portuguese; keep new ones consistent with that.

## Commands

### Backend (Go)
Run from `backend/`:
```
go run cmd/api/main.go     # start API server on :8080
go build ./...             # compile everything
go vet ./...                # lint
go test ./...                # run tests (no test files currently exist)
```
Requires a Postgres database and a `.env` file (see `backend/.env`, gitignored) with `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `SERVER_PORT`, `JWT_SECRET`. On every boot, `main.go` calls `config.ConnectDatabase()` (GORM `AutoMigrate`) followed by `seed.Run(db)`, which is idempotent (upserts) and repopulates classes/races/skills/pericias/talentos — expect this to run on every local start.

### Frontend (React/Vite)
Run from `frontend/`:
```
npm run dev       # start dev server on :5173
npm run build     # tsc -b && vite build
npm run lint      # eslint .
npm run preview   # preview production build
```
`VITE_API_URL` env var overrides the API base URL (defaults to `http://localhost:8080/api`).

### AI service (Python)
Run from `ai-service/`:
```
python main.py                 # or: uvicorn main:app --reload  (port 8000)
python ingest.py               # OCRs PDFs from ai-service/books/ and indexes them into ./chroma_db
```
Requires a local Ollama instance running the `nomic-embed-text` (embeddings) and `llama3.2` (generation) models, and Tesseract OCR installed (path is hardcoded in `ingest.py` for Windows). `ingest.py` must be run once to populate `chroma_db/` before `/skills` queries return results. There is no `requirements.txt`; dependencies are `fastapi`, `chromadb`, `langchain-ollama`, `ollama`, `pytesseract`, `pillow`, `pymupdf`, `langchain-text-splitters`.

## Architecture

### Backend layering
Strict `handler → service → repository → domain` layering, all wired manually (no DI framework) in `cmd/api/main.go`:
- `internal/domain/` — GORM models (`Character`, `Class`, `Race`, `Skill`, `Pericia`, `Talento`, `Antecedent`, `Background`, `Armor`, `User`)
- `internal/repository/` — GORM queries, one file per entity
- `internal/service/` — business logic (edition-specific math lives here, see below)
- `internal/handler/` — Gin handlers, one file per entity, only parse/validate HTTP and delegate to services
- `internal/seed/` + `pkg/config/seed*.go` — large static data files (classes, races, skills, pericias, talentos per edition); `seed.go` alone is 300+KB, read it with `offset`/`limit` rather than in full
- `pkg/middleware/auth.go` — JWT bearer middleware, sets `userID`/`userRole` in Gin context
- `pkg/config/database.go` — Postgres connection + `AutoMigrate`

Routes are registered in `main.go`: most reads are public, writes to `characters`/`classes`/`races`/`skills` and pericias/talentos-on-character are behind `AuthMiddleware`. Auth uses `Authorization: Bearer <jwt>`.

### Dual-edition logic (4e vs 5e)
`internal/service/character_service.go` is the heart of the ruleset differences — read it before touching any character math:
- Separate XP tables (`xpTable4e`, `xpTable5e`) and max level (30 vs 20)
- HP: 4e uses `BaseHP + CON` then flat `HPPerLevel` per level; 5e uses hit die + CON mod, half-die+1+CON mod per level
- 4e has Surge Value/Surges-per-day and four defenses (`Defense_AC/Fort/Refl/Will`, all `10 + level/2 + ability mod + class bonus`); 5e has `ProficiencyBonus` (scales by level breakpoints) and Death Saving Throws (`DeathSaveSuccesses/Failures`)
- ASI (Ability Score Improvement) levels differ: 5e is fixed levels {4,8,12,16,19}, 4e is every even level; `checkAndApplyLevelUps` auto-advances level-by-level on XP gain and pauses at ASI levels waiting for a player choice via `ApplyASI`
- `domain.Character` carries both rule sets' fields simultaneously (e.g. `AntecedentID`/`Alignment`/`PersonalityTraits` are 5e-only, `Defense_*`/`SurgeValue` are 4e-only) rather than being split into subtypes — when adding character fields, follow this pattern and gate the logic by `character.Edition`, not the schema
- "Pericias" (5e skill proficiencies, e.g. Stealth/Perception) and "Talentos" (5e feats) are distinct from the older `Skill` model (4e-style class powers with `power_type`: at-will/encounter/daily/utility) — don't conflate the two systems

### Frontend structure
Flat structure, no feature folders: `pages/` (routed screens), `components/` (shared UI), `services/` (one Axios-based module per backend resource, all importing the shared `services/api.ts` client), `store/authStore.ts` (Zustand, persists `token`/`user` to `localStorage`), `types/index.ts` (shared TS types). Routing is in `App.tsx` via `react-router-dom`; all routes except `/login` and `/register` are wrapped in `PrivateRoute`, which checks `useAuthStore`. `services/api.ts` attaches the JWT to every request, force-logs-out on 401, and shows a "waking up server" toast if a request takes >4s (the backend is hosted on a free tier with cold starts).

### Deployment
Frontend deploys to Vercel (SPA rewrite in `frontend/vercel.json`); backend is expected to be a separately hosted service (CORS in `main.go` explicitly allowlists `http://localhost:5173` and `https://rpg-manager-smoky.vercel.app` — update this list, not a wildcard, when adding origins).
