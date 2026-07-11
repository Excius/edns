
# Event Driven Notification System (EDNS)

A production-grade notification system built with Go that demonstrates modern backend architecture, event-driven processing, distributed workers, reliable message delivery, and real-time notifications.

The project is designed with scalability and fault tolerance in mind by separating responsibilities into independent services communicating through Redis Streams and Redis Pub/Sub.

---

# 📋 Table of Contents

* High-Level Architecture
* Key Features
* Development Phases
* Project Structure
* Tech Stack
* Getting Started
* System Design Concepts

---

# High-Level Architecture

```
                         Client
                            │
          ┌─────────────────┴─────────────────┐
          │                                   │
          ▼                                   ▼
      REST API                         WebSocket Client
          │                                   ▲
          ▼                                   │
     API Service                    WebSocket Service
          │                                   ▲
          │ Redis Streams                     │ Redis Pub/Sub
          ▼                                   │
      Worker Service ─────────────────────────┘
          │
          ▼
      PostgreSQL
```

---

# System Flow

1. Client sends a notification request through the REST API.
2. API validates the request and stores notification data in PostgreSQL.
3. API publishes the notification ID to Redis Streams.
4. Worker Service consumes the stream using Redis Consumer Groups.
5. Worker processes every delivery channel.
6. WebSocket deliveries are published to Redis Pub/Sub.
7. WebSocket Service delivers notifications to all active user connections.
8. Failed deliveries are retried automatically and eventually moved to the Dead Letter Queue.

---

# Key Features

* ✅ REST API
* ✅ PostgreSQL Persistence
* ✅ Event Driven Architecture
* ✅ Redis Streams
* ✅ Redis Consumer Groups
* ✅ Reliable Message ACK
* ✅ Retry Logic
* ✅ Dead Letter Queue (DLQ)
* ✅ Pending Message Recovery (XAUTOCLAIM)
* ✅ Multiple Worker Support
* ✅ Dedicated WebSocket Service
* ✅ Redis Pub/Sub
* ✅ Multiple Active Sessions Per User
* ✅ Graceful Shutdown
* ✅ Structured Logging

---

# Development Phases

## ✅ Phase 1 - Core Backend

Completed

Features:

* PostgreSQL integration
* Database migrations
* Repository layer
* Service layer
* REST APIs
* Configuration management
* Structured logging
* Docker environment

---

## ✅ Phase 2 - Event Driven Processing

Completed

Features:

* Redis Streams
* Queue Publisher
* Worker Service
* Notification Processor
* Delivery Processing
* Status Transitions
* Graceful Shutdown

Architecture:

```
API Service
      │
      ▼
Redis Streams
      │
      ▼
Worker Service
```

---

## ✅ Phase 3 - Reliable Queue Processing

Completed

Features:

* Redis Consumer Groups
* Message ACK
* Retry Count
* Dead Letter Queue
* Pending Message Recovery
* Worker Recovery
* XAUTOCLAIM Recovery

Goals Achieved:

* Prevent duplicate processing
* Prevent message loss
* Fault tolerant worker recovery

---

## ✅ Phase 4 - Real-Time Delivery

Completed

Features:

* Dedicated WebSocket Service
* Connection Hub
* User Session Tracking
* Multiple Active Connections
* Redis Pub/Sub
* Real-Time Notification Push

Architecture:

```
Worker Service
       │
Redis Pub/Sub
       │
       ▼
WebSocket Service
       │
       ▼
Connected Clients
```

---

## 🚧 Phase 5 - Email Delivery

Planned

Features:

* SMTP Integration
* Email Templates
* Delivery Tracking

---

## ⏳ Phase 6 - Observability

Planned

Features:

* Prometheus Metrics
* Grafana Dashboard
* Health Checks
* Distributed Logging

---

## ⏳ Phase 7 - Scaling

Planned

Features:

* Multiple Worker Instances
* Horizontal Scaling
* Distributed Processing
* High Availability

---

# Current Project Structure

```
api-service/
worker-service/
websocket-service/

internal/
├── config/
├── events/
├── logger/
├── models/
├── normalize/
├── queue/
└── repository/

configs/
docs/
migrations/
```

---

# Tech Stack

### Language

* Go

### Web Framework

* Gin

### Database

* PostgreSQL
* pgx

### Messaging

* Redis Streams
* Redis Pub/Sub

### Real-Time Communication

* Gorilla WebSocket

### Logging

* Uber Zap

### Containerization

* Docker
* Docker Compose

---

# Current System Architecture

```
REST API
    │
    ▼
PostgreSQL
    │
    ▼
Redis Streams
    │
    ▼
Worker Service
    │
    ├─────────────── Email (Phase 5)
    │
    └───────────────► Redis Pub/Sub
                            │
                            ▼
                  WebSocket Service
                            │
                            ▼
                  Connected Clients
```

---

# Reliability Features

* Redis Consumer Groups
* Message Acknowledgement
* Pending Message Recovery
* Retry Count Tracking
* Dead Letter Queue
* Graceful Shutdown
* Idempotent Processing

---

# Future Improvements

* SMTP Email Delivery
* Metrics & Monitoring
* Horizontal Scaling
* Authentication
* Rate Limiting
* Kubernetes Deployment

---

# System Design Concepts Demonstrated

* Event-Driven Architecture
* Distributed Workers
* Producer-Consumer Pattern
* Redis Streams
* Redis Consumer Groups
* Redis Pub/Sub
* Retry Mechanisms
* Dead Letter Queue
* Real-Time Notifications
* WebSocket Connection Management
* Graceful Shutdown
* Fault Tolerance
* Horizontal Scalability
