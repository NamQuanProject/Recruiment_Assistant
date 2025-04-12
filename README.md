# AI-Powered Recruitment System

## 🧠 Introduction

The **AI-Powered Recruitment System** is an intelligent recruitment assistant designed to automate the analysis and evaluation of resumes (CVs) based on job descriptions (JDs) provided by employers. By leveraging cutting-edge AI models, the system can:

- Match candidates to job descriptions with high accuracy.
- Assess the validity and relevance of information in resumes.
- Minimize bias in the hiring process.
- Offer transparent scoring and explainable decision-making to improve fairness and clarity for both recruiters and applicants.

This platform not only streamlines the recruitment pipeline but also enables companies to focus on finding the most suitable candidate — not just the most accomplished one. Additionally, the system is extensible and can be adapted to various use cases, such as automatically generating promotional videos from job descriptions or supporting a network of companies with unique recruitment features.

---

## 🚀 How to Run

### 🧰 Prerequisites

Ensure you have the following installed:

- **Python** (≥ 3.12)
- **Go** (≥ 1.19)
- **Node.js** and **npm**
- **Required Libraries**:
  - `PyMuPDF` (`fitz`) for PDF parsing
  - `react` for frontend
  - Other Go/Python dependencies as specified in respective modules

Install Python packages:

```bash
pip install PyMuPDF
```

Install Node.js dependencies (in webs directory):

```bash
cd webs
npm install
npm install react react-dom react-router-dom react-select
npm install framer-motion
```

### 🖥️ Running the Project
Start each service in a separate terminal:
1. Start the API Service (Go):
```bash
go run cmd/api/main.go
```
2. Start the Backend Services (Go):
```bash
go run cmd/backend/main.go
```
3. Start the Frontend (React):
```bash
cd webs
npm run dev
```
Then, open your browser and go to: http://localhost:5173/

---

## 🏗️ Architecture

```
recruitment-system/
├── cmd/                  # Application entry points
├── internal/             # Private application code
├── pkg/                  # Public libraries that can be used by external applications
├── webs/                 # Web frontend assets and components
├── test/                 # Test data and test utilities
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── README.md             # Project readme
└── Makefile              # Build automation
```

### Detailed Component Structure

#### cmd/ - Application Entry Points

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

#### internal/ - Private Application Code

```
internal/
├── apigateway/              # B: API Gateway
│   ├── middleware/
│   │   ├── auth.go          # JWT authentication
│   │   ├── ratelimit.go     # Rate limiting
│   │   └── routing.go       # Request routing
│   ├── handlers/
│   │   ├── get_Hl_CV.go        # Handle highlighting
│   │   ├── submit_CVS.go       # Handle submit CVs
│   │   └── submit_jd.go        # Handle submit JD
│   └── server.go            # API Gateway server setup
│
├── backend/                     # C: Backend Services
│   ├── evaluation/              # Evaluation service
│   │   ├── evaluator.go         # (scoring, bias, explanation, authentication, final_scores)
│   │   └── server.go            # gin server for evaluation
│   ├── parsing/                 # Input processing service
│   │   ├── extract_pdf.py       # Extract from pdf
│   │   ├── helper.py            # Some supporting function
│   │   ├── parse.go             # JD and CVs parser implementation
│   │   └── server.go            # gin server for parsing
│   ├── highlight/               # Highlighting pdf services
│   │   ├── calibrate.go         # Handles calibration of y-offsets for highlighting in PDFs
│   │   ├── calibrate_offset.py  # A Python script for calibrating offsets in PDF rendering
│   │   ├── extract_pdf_text.py  # A Python script for extracting text and positions from a PDF
│   │   ├── find_areas.go        # Contains logic to identify areas in a CV highlighting based on job details and evaluation
│   │   ├── highlight_pdf.py     # A Python script for adding highlights and annotations to a PDF
│   │   ├── pdf_extractor.go     # Extracts text blocks from a PDF file
│   │   └── server.go            # gin server for highlight services
│   └── output/              # Final output services
│       ├── process.go       # Implements logic for processing data and generating output
│       └── server.go        # Implements a server to handle requests related to output processing
│
├── aiservices/              # D: AI Services
│   ├── chatbot_singleton.go     # Implements a singleton pattern for managing chatbot instances
│   ├── gemini_areas.go          # Handles area-related logic for the Gemini AI service
│   ├── gemini_call.go           # Manages API calls to the Gemini service
│   ├── gemini_category.go       # Implements category-related logic for Gemini
│   ├── gemini_chatbot.go        # Contains chatbot functionalities specific to Gemini
│   ├── gemini_evaluate.go       # Handles evaluation logic for Gemini
│   ├── gemini_initialize.go     # Manages initialization processes for Gemini
│   ├── gemini_output.go         # Handles output generation for Gemini
│   ├── gemini_parsing.go        # Implements parsing logic for Gemini
│   ├── gemini_reading_links.go  # Processes and reads links for Gemini
│   ├── model.go                 # Defines data models for the AI services
│   ├── output.json              # Likely contains output data generated by the AI services
│   ├── prompts.go               # Manages prompts for AI interactions
│   ├── server.go                # Implements the server for AI services
│   └── utils.go                 # Contains utility functions for the AI services
```

#### webs/ - Web Frontend

```
webs/
├── src/                      # D: Source Code for the Web Application
│   ├── App.tsx               # The main application component for the frontend
│   ├── assets/               # Directory for static assets such as images, fonts, or icons
│   ├── components/           # D: Reusable UI Components
│   │   ├── chatbox.tsx       # Component for displaying and managing chat interactions
│   │   ├── criteria.tsx      # Component for displaying evaluation criteria
│   │   ├── datacontext.tsx   # Context provider for managing shared data across components
│   │   ├── footer.tsx        # Footer component for the application
│   │   ├── inputbox.tsx      # Input box component for user inputs
│   │   ├── navbar.tsx        # Navigation bar component
│   │   └── personDisplay.tsx # Component for displaying person-related information
│   ├── index.css             # The main CSS file for styling the application
│   ├── main.tsx              # The entry point for the React application
│   ├── pages/                # D: Page-Level Components
│   │   ├── candidateDetailPage.tsx  # Page for displaying detailed information about a candidate
│   │   ├── dashboard.tsx            # Dashboard page for an overview of application data
│   │   └── inputPage.tsx            # Page for handling user inputs
```
