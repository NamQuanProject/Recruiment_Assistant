# AI-Powered Recruitment System - Golang Project Structure

## Project Root Structure

```
recruitment-system/
├── cmd/                  # Application entry points
├── internal/             # Private application code
├── pkg/                  # Public libraries that can be used by external applications
├── api/                  # API definitions (OpenAPI/Swagger specs, protocol definitions)
├── web/                  # Web frontend assets and components
├── configs/              # Configuration files
├── deployments/          # Deployment configurations and templates
├── scripts/              # Scripts for development, CI/CD, etc.
├── test/                 # Test data and test utilities
├── docs/                 # Documentation
├── vendor/               # Dependencies (if using vendoring)
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── README.md             # Project readme
└── Makefile              # Build automation
```

## Detailed Component Structure

### cmd/ - Application Entry Points

```
cmd/
├── api/
│   └── main.go           # API Gateway service entry point
├── frontend/
│   └── main.go           # Frontend service entry point
├── backend/
│   └── main.go           # Backend service entry point
└── worker/
    └── main.go           # Background worker entry point
```

### internal/ - Private Application Code

```
internal/
├── apigateway/              # B: API Gateway
│   ├── middleware/
│   │   ├── auth.go          # JWT authentication
│   │   ├── ratelimit.go     # Rate limiting
│   │   └── routing.go       # Request routing
│   ├── handlers/
│   │   ├── proxy.go         # Service proxying
│   │   └── router.go        # Route definitions
│   └── server.go            # API Gateway server setup
│
├── backend/                 # C: Backend Services
│   ├── evaluation/          # Evaluation service
│   │   ├── criteria.go      # Criteria service implementation
│   │   ├── evaluator.go     # (scoring, bias, explanation, authentication, final_scores)
│   │   └── server.go        # gin server for evaluation
│   ├── parsing/             # Input processing service
│   │   ├── cvparser.go      # CV parser implementation
│   │   ├── jdparser.go      # JD parser implementation
│   │   └── server.go        # gin server for parsing
│   ├── user/                # User services
│   │   ├── account.go       # User account management
│   │   ├── preferences.go   # User preferences
│   │   └── server.go        # gin server for user services
│   └── output/              # Final output services
│       ├── formatter.go     # Output formatting
│       ├── generator.go     # Report generation
│       └── server.go        # gin server for output services
│
├── aiservices/              # D: AI Services
│   ├── parsing/
│   │   ├── cvparser.go      # CV parsing with AI
│   │   └── jdparser.go      # JD parsing with AI
│   ├── bias/
│   │   ├── detector.go      # Bias detection algorithms
│   │   └── mitigation.go    # Bias mitigation strategies
│   ├── scoring/
│   │   ├── matcher.go       # CV-JD matching algorithm
│   │   └── ranker.go        # Candidate ranking
│   ├── explanation/
│   │   └── generator.go     # Explanation generation
│   └── client.go            # AI service client
│
├── database/                # E: Database
│   ├── models/
│   │   ├── job.go           # Job posting model
│   │   ├── candidate.go     # Candidate model
│   │   ├── resume.go        # Resume storage model
│   │   └── company.go       # Company model
│   ├── repositories/
│   │   ├── job_repo.go      # Job repository
│   │   ├── candidate_repo.go # Candidate repository
│   │   ├── resume_repo.go   # Resume repository
│   │   └── company_repo.go  # Company repository
│   └── db.go                # Database connection and setup
│
├── external/                # F: External APIs
│   ├── gemini/
│   │   └── client.go        # Gemini AI client
│   ├── github/
│   │   └── client.go        # GitHub API client
│   └── extapi.go            # External API client interfaces
│
└── common/                  # Shared code
    ├── config/
    │   └── config.go        # Configuration loading
    ├── logger/
    │   └── logger.go        # Logging utilities
    ├── auth/
    │   └── jwt.go           # JWT handling
    ├── errors/
    │   └── errors.go        # Error definitions
    └── types/
        └── types.go         # Common type definitions
```

### pkg/ - Public Libraries

```
pkg/
├── cvanalysis/             # CV analysis utilities
│   ├── parser.go           # CV parsing utilities
│   ├── extractor.go        # Information extraction
│   └── validator.go        # CV validation
├── jdanalysis/             # Job description analysis
│   ├── parser.go           # JD parsing utilities
│   └── extractor.go        # Requirements extraction
├── bias/                   # Bias detection utilities
│   ├── detector.go         # Bias detection algorithms
│   └── metrics.go          # Bias measurement metrics
└── matching/               # CV-JD matching utilities
    ├── algorithm.go        # Matching algorithm
    └── scoring.go          # Candidate scoring
```

### api/ - API Definitions

```
api/
├── swagger/                # OpenAPI/Swagger specs
│   ├── api.yaml            # Main API specification (combines all services)
│   ├── evaluation.yaml     # Evaluation service endpoints (optional split)
│   ├── user.yaml           # User service endpoints (optional split)
│   └── ...                 # Other modular specs
└── schemas/                # Shared request/response schemas (optional)
    ├── User.json
    ├── Evaluation.json
    └── ...
```

### web/ - Web Frontend

```
web/
├── static/                 # Static assets
│   ├── css/                # CSS files
│   ├── js/                 # JavaScript files
│   └── images/             # Image files
├── templates/              # HTML templates
│   ├── admin/              # Admin dashboard templates
│   ├── upload/             # Upload form templates
│   └── results/            # Results display templates
└── components/             # Web components
    ├── form/               # Form components
    ├── dashboard/          # Dashboard components
    └── evaluation/         # Evaluation result components
```
