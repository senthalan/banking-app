# 🏦 Banking Application

A modern, full-stack banking application demonstrating microservices architecture with secure authentication, transaction management, and automated reporting.

## 🚀 Features

- User management and authentication
- Bank account operations
- Transaction processing
- Daily email reports
- RESTful API with OpenAPI docs
- Responsive React UI

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   React SPA     │───▶│   Go Lang        │───▶│   MySQL DB      │
│  (Frontend)     │    │  (Backend)       │    │   (Database)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                ▲
                                │
                                │
                       ┌──────────────────┐
                       │  Email Service   │
                       │ (Daily Reports)  │
                       └──────────────────┘
```

## 📁 Project Structure

```
banking-app/
├── backend/          # Go backend service (Gin framework)
├── frontend/         # React frontend (Vite + TypeScript)
├── task/            # Email notification service
├── tester/          # Postman collection for API testing
└── README.md        # This file
```

## 🛠️ Technology Stack

### Frontend
- React 19, TypeScript, Material-UI, Vite

### Backend
- Go 1.24, Gin, GORM, MySQL

### Deployment
- WSO2 Choreo, Docker, OpenAPI 3.0

## 🚀 Quick Start

### Prerequisites

- Node.js >= 21.0.0, Go >= 1.21, MySQL >= 8.0
- npm or yarn
- WSO2 Choreo account (for deployment)

### Setup

1. Clone the repository
2. Backend: `cd backend && go mod tidy && go run main.go`
3. Frontend: `cd frontend && npm install && npm run dev`
4. Email service: `cd task && go mod tidy && go run main.go` (optional)

Configure environment variables for database and SMTP settings as needed.
