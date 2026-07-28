# Event Driven Notification System (EDNS v1.0.0)

A production-grade event-driven notification system built with Go that demonstrates asynchronous processing, reliable message delivery, real-time notifications, and a scalable microservice architecture.

The system is built around independent services communicating through Redis Streams and Redis Pub/Sub, making it easy to extend with additional notification channels.

---

# Table of Contents

- Overview
- Architecture
- Notification Flow
- Features
- Project Structure
- Tech Stack
- Running the Project
- Reliability & Observability
- Future Improvements

---

# Overview

EDNS is composed of three independent services:

- API Service
- Worker Service
- WebSocket Service

Notifications are created through the REST API, processed asynchronously by workers, and delivered through multiple notification channels.

---

# Architecture

```text
                Client
                  │
                  ▼
            API Service
                  │
        PostgreSQL + Redis Streams
                  │
                  ▼
           Worker Service
           │             │
           ▼             ▼
      SMTP Email   Redis Pub/Sub
                         │
                         ▼
                 WebSocket Service
                         │
                         ▼
                 Connected Clients
```

---

# Notification Flow

1. Client creates a notification through the REST API.
2. API stores notification data in PostgreSQL.
3. Notification ID is published to Redis Streams.
4. Worker consumes the message using Redis Consumer Groups.
5. Worker processes each delivery channel.
6. Email notifications are sent through SMTP.
7. WebSocket notifications are published through Redis Pub/Sub.
8. WebSocket Service pushes notifications to all connected user sessions.
9. Failed deliveries are retried and eventually moved to the Dead Letter Queue (DLQ).

---

# Features

### Core

- REST API
- PostgreSQL Persistence
- Repository-Service Architecture
- Configuration Management
- Docker & Docker Compose

### Event Processing

- Redis Streams
- Redis Consumer Groups
- Reliable Message Acknowledgement
- Automatic Retry Logic
- Dead Letter Queue (DLQ)
- Pending Message Recovery (`XAUTOCLAIM`)
- Multiple Worker Support

### Notification Channels

- Email Notifications (SMTP)
- WebSocket Notifications
- Redis Pub/Sub
- Multiple Active Connections Per User

### Observability

- Structured Logging (Zap)
- Health Checks
- Readiness Checks
- Prometheus Metrics

### Reliability

- Graceful Shutdown
- Fault Tolerant Workers
- Delivery Status Tracking
- Notification Status Synchronization

---

# Project Structure

```text
api-service/
worker-service/
websocket-service/

internal/
├── config/
├── events/
├── logger/
├── models/
├── observability/
├── repository/
├── stream/
└── validation/

configs/
deploy/
migrations/
```

---

# Tech Stack

**Language**

- Go

**Framework**

- Gin

**Database**

- PostgreSQL
- pgx

**Messaging**

- Redis Streams
- Redis Pub/Sub

**Real-Time Communication**

- Gorilla WebSocket

**Email**

- SMTP
- Mailpit

**Observability**

- Prometheus
- Grafana
- Zap Logger

**Containerization**

- Docker
- Docker Compose

---

# Running the Project

Production

```bash
make docker-up-prod
```

Development

```bash
make docker-up-dev
```

Available services:

| Service    | Address                    |
| ---------- | -------------------------- |
| API        | http://localhost:8080      |
| Worker     | http://localhost:8081      |
| WebSocket  | ws://localhost:8082/api/ws |
| Mailpit    | http://localhost:8025      |
| Prometheus | http://localhost:9090      |
| Grafana    | http://localhost:3001      |

---

# Reliability & Observability

### Reliability

- Redis Consumer Groups
- Automatic Retries
- Dead Letter Queue
- Pending Message Recovery
- Graceful Shutdown
- Multiple Worker Support

### Observability

- Structured Logging
- Health Endpoint
- Readiness Endpoint
- Prometheus Metrics

---

# System Design Concepts

- Event-Driven Architecture
- Producer-Consumer Pattern
- Repository Pattern
- Service Layer
- Distributed Workers
- Redis Streams
- Redis Consumer Groups
- Redis Pub/Sub
- Dead Letter Queue
- Retry Mechanism
- WebSocket Connection Management
- Fault Tolerance
- Horizontal Scalability

---

# Future Improvements

- Unit & Integration Tests
- Grafana Dashboards
- OpenTelemetry Tracing
- CI/CD Pipeline
- Kubernetes Deployment
- Authentication & Authorization
- Rate Limiting
- Additional Notification Channels
