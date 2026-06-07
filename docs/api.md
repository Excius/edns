# API Documentation

## Health Check

### GET /health

Response:

{
"status": "ok"
}

---

## Users

### Create User

POST /users

Request:

{
"email": "[user@example.com](mailto:user@example.com)"
}

Response:

{
"id": "uuid",
"email": "[user@example.com](mailto:user@example.com)"
}

---

### Get User

GET /users/:id

Response:

{
"id": "uuid",
"email": "[user@example.com](mailto:user@example.com)"
}

---

## Notifications

### Create Notification

POST /notifications

Request:

{
"user_id": "uuid",
"message": "Hello",
"channels": [
"email",
"websocket"
]
}

Response:

{
"id": "uuid",
"user_id": "uuid",
"message": "Hello",
"status": "pending"
}

---

### Get Notification

GET /notifications/:id

Response:

{
"id": "uuid",
"user_id": "uuid",
"message": "Hello",
"status": "completed"
}

---

## Status Values

Notification Status:

* pending
* completed
* failed

Delivery Status:

* pending
* processing
* completed
* failed

---

## Supported Channels

* email
* websocket
