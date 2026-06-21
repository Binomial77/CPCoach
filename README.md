# cp://coach

An AI competitive programming coach. Log the problems you solve, track your rating, and let Gemini turn your submission history into a concrete training plan instead of generic advice.

![Go](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go\&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Web%20Framework-00ADD8)
![SQLite](https://img.shields.io/badge/SQLite-GORM-003B57?logo=sqlite\&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-lightgrey)

## What it does

Most "track your CP progress" tools stop at a spreadsheet. This one watches your rating curve and per-topic accuracy (graphs, trees, greedy, DP, binary search, number theory), then asks Gemini to read that pattern and tell you what to actually grind next: strengths, weaknesses, a recommended difficulty band, and a short training plan. The guidance is regenerated every time you log a solve.

* Email/password authentication with bcrypt-hashed passwords and JWT stored in an httpOnly cookie
* Codeforces-style rating ladder (Newbie → Legendary Grandmaster) computed client-side from your current rating
* One-tap problem logging by topic tag and difficulty band
* AI-generated guidance via the Gemini API, based on your actual stats, not a static template
* Server-rendered dashboard (Gin's `html/template`) with a black, white, and grey "judge console" interface
* No frontend framework and no build step

## Tech Stack

| Layer          | Choice                                    |
| -------------- | ----------------------------------------- |
| Language       | Go                                        |
| Web Framework  | Gin                                       |
| ORM / Database | GORM + SQLite                             |
| Authentication | JWT (httpOnly cookie) + bcrypt            |
| AI             | Google Gemini (`google.golang.org/genai`) |
| Frontend       | Vanilla HTML, CSS, and JavaScript         |

## Project Structure

```text
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
├── utils/
│   └── jwt.go              # JWT generation and validation
├── templates/
│   ├── home.html
│   ├── login.html
│   ├── signup.html
│   └── dashboard.html
├── static/
│   ├── css/style.css
│   └── js/app.js
├── main.go
├── go.mod
├── go.sum
├── .env.example            # template, copy to .env and fill in real values
└── .env                    # not committed
```

## Getting Started

### Prerequisites

* Go 1.26 or newer
* A Gemini API key from Google AI Studio

### 1. Clone the Repository

```bash
git clone https://github.com/<your-username>/<your-repo>.git
cd <your-repo>
```

### 2. Configure Environment Variables

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

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Application

```bash
go run main.go
```

GORM's `AutoMigrate()` creates the SQLite tables on first run.

Visit:

```text
http://localhost:8080
```

## API Reference

| Method | Path            | Auth | Description                                                   |
| ------ | --------------- | ---- | ------------------------------------------------------------- |
| GET    | `/`             | No   | Landing page                                                  |
| GET    | `/signup`       | No   | Signup page                                                   |
| POST   | `/signup`       | No   | Create account, sets JWT cookie, redirects to `/dashboard`    |
| GET    | `/login`        | No   | Login page                                                    |
| POST   | `/login`        | No   | Authenticate user, sets JWT cookie, redirects to `/dashboard` |
| GET    | `/dashboard`    | Yes  | Renders rating, total solved, and topic breakdown             |
| POST   | `/postproblem`  | Yes  | Logs a solved problem (rating + topic tags)                   |
| POST   | `/updaterating` | Yes  | Updates current rating                                        |
| GET    | `/getguidance`  | Yes  | Returns Gemini-generated coaching guidance                    |

Routes marked "Yes" require a valid `jwt` cookie.

## Rating Tiers

The dashboard colors the rating pill using a Codeforces-style ladder, computed in `static/js/app.js`.

| Tier                      | Range     |
| ------------------------- | --------- |
| Newbie                    | < 1200    |
| Pupil                     | 1200–1399 |
| Specialist                | 1400–1599 |
| Expert                    | 1600–1899 |
| Candidate Master          | 1900–2099 |
| Master                    | 2100–2299 |
| International Master      | 2300–2399 |
| Grandmaster               | 2400–2599 |
| International Grandmaster | 2600–3499 |
| Legendary Grandmaster     | ≥ 3500    |

## Security Notes

* Passwords are hashed with bcrypt before storage and are never stored in plaintext.
* Authentication uses JWT tokens stored in an httpOnly cookie.
* Protected routes are validated through `middleware.Authentication()`.
* `.env` and SQLite database files are excluded from version control through `.gitignore`.

