# LootForge

Fair loot distribution system for gaming groups. Supports multiple distribution strategies with group management, roll tracking, and fairness metrics.

## Architecture

```
core/     Go REST API (net/http, Firestore)
ui/       Next.js 14 dashboard (React 18, MUI v9)
```

## Features

- **Google Sign-In** via Firebase Authentication
- **Group management** — create, update, delete groups with configurable loot strategies
- **Invite system** — generate time-limited invite links, accept/revoke invites
- **4 loot strategies**:
  - **Pure Random** — equal chance for all participants
  - **Weighted Luck** — luck stat increases for losers, decreases for winners
  - **Round Robin** — priority-based rotation
  - **DKP** — Dragon Kill Points; highest bidder wins, points deducted
- **Roll history** with pagination and filtering
- **Statistics dashboard** — winner distribution charts, fairness metrics (std dev)
- **Role-based access** — Owner, Admin, Member roles per group
- **Rate limiting** — per-IP request throttling

## Quick Start

### Prerequisites

- Docker & Docker Compose
- (Optional) Go 1.24+, Node.js 20+

### Run with Docker Compose

```bash
docker compose up --build
```

This starts:
- Firestore emulator on `:8081`
- Go API on `:8080`
- Next.js UI on `:3000`

Open http://localhost:3000 and sign in with Google.

### Local Development

**Backend:**

```bash
cd core
export FIRESTORE_EMULATOR_HOST=localhost:8081
export GOOGLE_CLOUD_PROJECT=lootforge-dev
go run ./cmd/server
```

**Frontend:**

```bash
cd ui
cp .env.example .env.local   # configure Firebase keys
npm install
npm run dev
```

### Environment Variables

| Variable | Service | Description |
|---|---|---|
| `GOOGLE_CLOUD_PROJECT` | core | GCP/Firebase project ID |
| `FIRESTORE_EMULATOR_HOST` | core | Firestore emulator address (dev only) |
| `CORS_ORIGINS` | core | Comma-separated allowed origins |
| `PORT` | core | API listen port (default: 8080) |
| `NEXT_PUBLIC_API_URL` | ui | Backend API URL |
| `NEXT_PUBLIC_FIREBASE_API_KEY` | ui | Firebase client API key |
| `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN` | ui | Firebase auth domain |
| `NEXT_PUBLIC_FIREBASE_PROJECT_ID` | ui | Firebase project ID |

## API Endpoints

### Auth
| Method | Path | Description |
|---|---|---|
| `GET` | `/api/health` | Health check (public) |
| `GET` | `/api/me` | Current user profile |

### Groups
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/groups` | Create group |
| `GET` | `/api/groups` | List user's groups |
| `GET` | `/api/groups/{id}` | Get group |
| `PUT` | `/api/groups/{id}` | Update group |
| `DELETE` | `/api/groups/{id}` | Delete group (owner only) |

### Members
| Method | Path | Description |
|---|---|---|
| `GET` | `/api/groups/{gid}/members` | List members |
| `POST` | `/api/groups/{gid}/members` | Add member |
| `PUT` | `/api/groups/{gid}/members/{id}` | Update member |
| `DELETE` | `/api/groups/{gid}/members/{id}` | Delete member |
| `PATCH` | `/api/groups/{gid}/members/{id}/points` | Adjust DKP points |

### Rolls
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/groups/{gid}/rolls` | Execute a roll |
| `GET` | `/api/groups/{gid}/rolls` | List roll history (paginated) |
| `GET` | `/api/groups/{gid}/rolls/stats` | Get win statistics |

### Invites
| Method | Path | Description |
|---|---|---|
| `POST` | `/api/groups/{gid}/invites` | Create invite (admin+) |
| `POST` | `/api/invites/{token}/accept` | Accept invite |
| `DELETE` | `/api/invites/{token}` | Revoke invite (admin+) |

## Testing

```bash
# Backend
cd core && go test ./...

# Frontend
cd ui && npm test
```

## Project Structure

```
core/
  cmd/server/         Entry point
  internal/
    auth/             Firebase token verification, authorization helpers
    handlers/         HTTP request handlers
    httputil/         JSON response helpers
    middleware/       Auth, CORS, logging, rate limiting
    models/           Domain types
    router/           Route definitions
    store/            Firestore data access
    strategies/       Loot distribution algorithms
ui/
  src/
    app/              Next.js pages
    components/       React components
    hooks/            Custom React hooks
    lib/              API client, Firebase config, theme
    types/            TypeScript type definitions
```

## License

Private — all rights reserved.
