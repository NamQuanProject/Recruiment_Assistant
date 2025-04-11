# API Documentation

## Overview

This document provides detailed information about all API endpoints in the Recruitment Assistant project. The system consists of multiple microservices running on different ports:

- API Gateway: `:8080`
- AI Services: `:8081`
- Evaluation Service: `:8082`
- Highlight Service: `:8083`
- Output Service: `:8084`
- Parsing Service: `:8085`

## Authentication Endpoints

### Signup

- **Endpoint**: `POST /signup`
- **Description**: Creates a new user account
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response**:
  ```json
  {
    "message": "User created successfully"
  }
  ```

### Login

- **Endpoint**: `POST /login`
- **Description**: Authenticates a user and returns a JWT token
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response**:
  ```json
  {
    "message": "Login successful",
    "token": "string"
  }
  ```

### Logout

- **Endpoint**: `POST /logout`
- **Description**: Logs out the current user
- **Response**:
  ```json
  {
    "message": "Logout successful"
  }
  ```

## Job Description Endpoints

### Submit Job Description

- **Endpoint**: `POST /submitJD`
- **Description**: Uploads a job description PDF and processes it
- **Request**: Form-data
  - `job_name`: string (required)
  - `pdf_file`: file (required, PDF format)
- **Response**:
  ```json
  {
    "message": "Job description processed successfully"
  }
  ```

## CV Processing Endpoints

### Submit CVs

- **Endpoint**: `POST /submitCVs`
- **Description**: Uploads and processes CVs (single PDF or ZIP containing multiple PDFs)
- **Request**: Form-data
  - `file`: file (required, PDF or ZIP format)
- **Response**:
  ```json
  {
    "message": "PDF uploaded and parsed successfully",
    "path": "string"
  }
  ```
  or for ZIP files:
  ```json
  {
    "message": "ZIP processed",
    "pdf_count": number,
    "saved_pdfs": ["string"]
  }
  ```

### Get Highlighted CV

- **Endpoint**: `POST /getHlCV`
- **Description**: Retrieves a highlighted version of a CV
- **Request Body**:
  ```json
  {
    "index": number
  }
  ```
- **Response**:
  ```json
  {
    "highlighted_pdf_path": "string"
  }
  ```

## AI Services Endpoints

### AI Chat

- **Endpoint**: `GET /ai`
- **Description**: General AI chat endpoint
- **Response**:
  ```json
  {
    "Question": "string",
    "Response": "string"
  }
  ```

### Job Category List

- **Endpoint**: `GET /ai/jd_category/`
- **Description**: Retrieves list of job categories
- **Response**: JSON object containing job categories and their data

### Job Category Details

- **Endpoint**: `GET /ai/jd_category/:job_name`
- **Description**: Retrieves details for a specific job category
- **Response**:
  ```json
  {
    "Response": object
  }
  ```

### CV Parsing

- **Endpoint**: `POST /ai/parsing`
- **Description**: Parses raw CV text
- **Request Body**:
  ```json
  {
    "job_raw_text": "string"
  }
  ```
- **Response**:
  ```json
  {
    "Question": "string",
    "Response": object
  }
  ```

### CV Evaluation

- **Endpoint**: `POST /ai/evaluate`
- **Description**: Evaluates a CV against job requirements
- **Request Body**:
  ```json
  {
    "job_name": "string",
    "jd_main_quiteria": "string",
    "cv_raw_text": "string",
    "evaluation_id": "string",
    "cv_id": "string"
  }
  ```
- **Response**:
  ```json
  {
    "evaluation": object
  }
  ```

### Chatbot Endpoints

#### Initialize Chatbot

- **Endpoint**: `POST /ai/chatbot/init`
- **Description**: Initializes a new chatbot session
- **Request Body**:
  ```json
  {
    "eval_id": "string"
  }
  ```
- **Response**:
  ```json
  {
    "message": "Chatbot initialized successfully",
    "evaluation_id": "string"
  }
  ```

#### Ask Chatbot

- **Endpoint**: `POST /ai/chatbot/ask`
- **Description**: Sends a question to the chatbot
- **Request Body**:
  ```json
  {
    "cv_id": "string",
    "question": "string"
  }
  ```
- **Response**:
  ```json
  {
    "answer": "string"
  }
  ```

## Output Service Endpoints

### Get Output

- **Endpoint**: `POST /output`
- **Description**: Retrieves evaluation results
- **Request Body**:
  ```json
  {
    "evaluation_folder": "string"
  }
  ```
- **Response**:
  ```json
  {
    "list": [
      {
        "full_name": "string",
        "worked_for": "string",
        "experience_level": "string",
        "authenticity": number,
        "final_score": number,
        "path_to_cv": "string",
        "path_to_evaluation": "string"
      }
    ]
  }
  ```

## Highlight Service Endpoints

### Highlight CV

- **Endpoint**: `POST /highlight`
- **Description**: Highlights areas in a CV
- **Request Body**:
  ```json
  {
    "pdf_path": "string",
    "areas": [
      {
        "text": "string",
        "page": number,
        "x": number,
        "y": number,
        "width": number,
        "height": number,
        "description": "string",
        "type": "string"
      }
    ]
  }
  ```
- **Response**:
  ```json
  {
    "highlighted_pdf_path": "string",
    "message": "string"
  }
  ```

## Error Responses

All endpoints may return the following error responses:

- **400 Bad Request**:

  ```json
  {
    "error": "string"
  }
  ```

- **401 Unauthorized**:

  ```json
  {
    "error": "string"
  }
  ```

- **404 Not Found**:

  ```json
  {
    "error": "string"
  }
  ```

- **500 Internal Server Error**:
  ```json
  {
    "error": "string"
  }
  ```
