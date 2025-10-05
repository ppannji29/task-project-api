🚀 Features Implemented
User Dummy Creation
POST /api/user/mydummy/create – Create dummy user for testing (no UI)

curl --location 'https://task-project-api-production.up.railway.app/api/user/mydummy/create' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "Panji Pp",
  "email": "testingdev29@gmail.com", // -> change to your email to get request otp from email
  "phone": "08123456789",
  "profile": {
    "address": "Jl. Dummy Raya No. 123",
    "photo": "",
    "bio": "Testing Dev Software Engineer",
    "birthdate": "1995-01-01T00:00:00Z",
    "age": 30,
    "gender": "male"
  }
}
'

🔐 Mock Authentication
Simulated OTP-based login flow with token handling.
POST /api/auth/request-otp – Request OTP (mocked)
POST /api/auth/verify-otp – Verify OTP and return token
POST /api/auth/refresh-token – Refresh access token
DELETE /api/auth/logout – Clear session (mocked)

📋 Task CRUD Module
All task routes are protected and require a valid token.
GET /api/tasks – List tasks (supports filter, sort, pagination)
GET /api/task/{id} – View task detail
POST /api/task – Create new task
PATCH /api/task/{id} – Update existing task
DELETE /api/task/{id} – Delete task

🛡️ Protected Routes
All /api/task/* endpoints are secured via middleware. Token must be present in cookies (access_token, refresh_token, or token) to access them.

🧪 Tech Stack
Language: Go (Golang)
Database: MongoDB (mocked or real)
Deployment: Railway

🌐 Live Deployment
🔗 APP: https://task-project-fe-production.up.railway.app