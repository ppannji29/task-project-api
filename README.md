## 🚀 Features Implemented

### 👤 User Dummy Creation

Create a dummy user for testing purposes (no frontend UI available).

**Endpoint:** `POST /api/user/mydummy/create`

```bash
curl --location 'https://task-project-api-production.up.railway.app/api/user/mydummy/create' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "Panji Pp",
  "email": "testingdev29@gmail.com", // ← change to your email to receive OTP
  "phone": "08123456789",
  "profile": {
    "address": "Jl. Dummy Raya No. 123",
    "photo": "",
    "bio": "Testing Dev Software Engineer",
    "birthdate": "1995-01-01T00:00:00Z",
    "age": 30,
    "gender": "male"
  }
}'
```

🔐 Mock Authentication
Simulated OTP-based login flow with token handling.

POST /api/auth/request-otp – Request OTP (mocked)

POST /api/auth/verify-otp – Verify OTP and return token

POST /api/auth/refresh-token – Refresh access token

DELETE /api/auth/logout – Clear session (mocked)

📋 Task CRUD Module

All task routes are protected and require a valid token.

GET /api/tasks – List tasks (supports filter, sort, pagination)

GET /api/task/{id} – View task detail, milestone

POST /api/task – Create new task

PATCH /api/task/{id} – Update existing task

DELETE /api/task/{id} – Delete task

🛡️ Protected Routes

All /api/task/* endpoints are secured token via middleware

## 📊 MongoDB Index Strategy

To ensure optimal query performance, indexes are automatically created through the `EnsureTaskIndexes()` function.

### 🧩 Index Overview

| **Field**         | **Purpose**                         |
|--------------------|-------------------------------------|
| `user_id`          | Filter tasks per user               |
| `status`           | Filter and sort by status           |
| `priority`         | Filter and sort by priority         |
| `due_date`         | Sort by deadline                    |
| `created_at`       | Sort by newest tasks                |
| `user_id + status` | Compound index for common queries   |

🧪 Tech Stack

Language: Go (Golang)

Database: MongoDB (mocked or real)

Deployment: Railway

🌐 Live Deployment

🔗 APP: https://task-project-fe-production.up.railway.app

