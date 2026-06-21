# cp://coach

An AI competitive programming coach. Log the problems you solve, track your rating, and let Gemini turn your submission history into a concrete training plan instead of generic advice.

![Go](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Web%20Framework-00ADD8)
![SQLite](https://img.shields.io/badge/SQLite-GORM-003B57?logo=sqlite&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-lightgrey)

## What it does

Most "track your CP progress" tools stop at a spreadsheet. This one watches your rating curve and per-topic accuracy (graphs, trees, greedy, DP, binary search, number theory), then asks Gemini to read that pattern and tell you what to actually grind next — strengths, weaknesses, a recommended difficulty band, and a short training plan, regenerated every time you log a solve.

- Email/password auth with bcrypt-hashed passwords and JWT stored in an httpOnly cookie
- Codeforces-style rating ladder (Newbie → Legendary Grandmaster) computed client-side from your current rating
- One-tap problem logging by topic tag and difficulty band
- AI-generated guidance via the Gemini API, based on your actual stats — not a static template
- Server-rendered dashboard (Gin's `html/template`) with a black/white/grey "judge console" interface — no frontend framework, no build step

## Tech stack

| Layer | Choice |
|---|---|
| Language | Go |
| Web framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM / DB | [GORM](https://gorm.io/) over SQLite |
| Auth | JWT (httpOnly cookie) + bcrypt |
| AI | Google Gemini (`google.golang.org/genai`) |
| Frontend | Vanilla HTML/CSS/JS, served as Gin templates + static assets |

## Project structure

```
.
├── controllers/
│   ├── controllers.go      # signup, login, dashboard, log problem, update rating
│   └── guidance.go         # builds the Gemini prompt and returns AI guidance
├── database/
│   └── config.go           # SQLite connection
├── initializers/
│   └── initializers.go     # loads .env
├── middleware/
│   └── auth.go             # JWT cookie verification, sets userID in context
├── models/
│   └── models.go           # User, UserRating, ProblemStat, request DTOs
├── routes/
│   └── routes.go           # route table
├── utils/                  # JWT generation/validation (referenced, not shown above)
├── templates/
│   ├── home.html
│   ├── login.html
│   ├── signup.html
│   └── dashboard.html
├── static/
│   ├── css/style.css
│   └── js/app.js
├── main.go
├── go.mod / go.sum
├── .env.example              # template — copy to .env and fill in real values
└── .env                       # not committed — see below
```

## Getting started

### Prerequisites

- Go 1.26 or newer
- A [Gemini API key](https://aistudio.google.com/apikey)

### 1. Clone the repo

```bash
git clone https://github.com/<your-username>/<your-repo>.git
cd <your-repo>
```

### 2. Configure environment variables

Copy the example file and fill in your own values:

```bash
cp .env.example .env
```

```env
PORT=8080
DB_NAME=your_db_name
JWT_SECRET=your_jwt_secret
API_KEY=gemini_api_key
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run it

```bash
go run main.go
```

GORM's `AutoMigrate` creates the SQLite tables on first run. Visit `http://localhost:8080`.

## API reference

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/` | — | Landing page |
| GET | `/signup` | — | Signup page |
| POST | `/signup` | — | Create account → sets JWT cookie → redirects to `/dashboard` |
| GET | `/login` | — | Login page |
| POST | `/login` | — | Authenticate → sets JWT cookie → redirects to `/dashboard` |
| GET | `/dashboard` | ✅ | Renders rating, total solved, and topic breakdown |
| POST | `/postproblem` | ✅ | Logs a solved problem (rating + topic tags) |
| POST | `/updaterating` | ✅ | Updates current rating |
| GET | `/getguidance` | ✅ | Returns Gemini-generated coaching guidance |

Routes marked ✅ require the `jwt` cookie set at login/signup.

## Rating tiers

The dashboard colors your rating pill using a Codeforces-style ladder, computed in `static/js/app.js`:

| Tier | Range |
|---|---|
| Newbie | < 1200 |
| Pupil | 1200–1399 |
| Specialist | 1400–1599 |
| Expert | 1600–1899 |
| Candidate Master | 1900–2099 |
| Master | 2100–2299 |
| International Master | 2300–2399 |
| Grandmaster | 2400–2599 |
| International Grandmaster | 2600–3499 |
| Legendary Grandmaster | ≥ 3500 |

## Security notes

- Passwords are hashed with bcrypt before storage — never stored or logged in plaintext.
- Auth uses a JWT in an httpOnly cookie, checked by `middleware.Authentication()` on every protected route.
- `.env` and the SQLite database file are excluded from version control (see `.gitignore` below) — both can contain real secrets or user data.

## License

[MIT](LICENSE) — or pick whatever you prefer; update this section and add a `LICENSE` file to match.
