# Event Driven Notification System (EDNS)

## Overview

EDNS is a distributed notification system built using Go.

The system is designed around an event-driven architecture where API requests are decoupled from notification delivery using Redis Streams and background workers.

Current architecture:

```
Client
↓
API Service
↓
Redis Stream
↓
Worker Service
↓
PostgreSQL
```
---

## Services

### API Service

Responsibilities:

* Handle HTTP requests
* Validate incoming data
* Persist notifications
* Publish notification events

The API service does not send notifications directly.

---

### Worker Service

Responsibilities:

* Consume notification events
* Process delivery records
* Update delivery status
* Update notification status

The worker operates asynchronously and independently of the API service.

---

### WebSocket Service (Planned)

Responsibilities:

* Maintain client connections
* Push real-time notifications
* Track connected users

This service will be introduced in a later phase.

---

## Database Schema

### users

Stores registered users.

Fields:

* id
* email
* created_at

---

### notifications

Stores notification metadata.

Fields:

* id
* user_id
* message
* status
* created_at

---

### notification_deliveries

Stores per-channel delivery state.

Fields:

* id
* notification_id
* channel
* status
* retry_count
* created_at

---

## Notification Flow

1. User creates notification
2. API stores notification
3. API stores delivery records
4. API publishes event to Redis Stream
5. Worker consumes event
6. Worker processes deliveries
7. Worker updates delivery status
8. Worker updates notification status

---

## Status Lifecycle

Notification:

pending
→ completed
→ failed

Delivery:

pending
→ processing
→ completed
→ failed

---

## Technology Stack

Language:

* Go

Database:

* PostgreSQL

Queue:

* Redis Streams

Logging:

* Zap

Configuration:

* Viper

Containerization:

* Docker
* Docker Compose
