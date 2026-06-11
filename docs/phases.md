# EDNS Development Roadmap

## Phase 1 - Core Backend

Completed

Features:

* PostgreSQL integration
* Migrations
* User APIs
* Notification APIs
* Repository layer
* Service layer
* Config management
* Structured logging
* Docker environment

---

## Phase 2 - Event Driven Processing

Completed

Features:

* Redis Streams
* Queue publisher
* Worker service
* Notification processor
* Delivery processing
* Status transitions
* Graceful shutdown

Architecture:

API
↓
Redis Stream
↓
Worker

---

## Phase 3 - Reliable Queue Processing

Planned

Features:

* Consumer Groups
* Message ACK
* Retry Logic
* Dead Letter Queue
* Worker Recovery

Goals:

* Prevent duplicate processing
* Prevent message loss
* Improve fault tolerance

---

## Phase 4 - Real Time Delivery

Planned

Features:

* WebSocket Service
* Connection Management
* User Session Tracking
* Real-Time Notification Push

---

## Phase 5 - Email Delivery

Planned

Features:

* SMTP Integration
* Email Templates
* Delivery Tracking

---

## Phase 6 - Observability

Planned

Features:

* Metrics
* Health Checks
* Prometheus
* Grafana

---

## Phase 7 - Scaling

Planned

Features:

* Multiple Workers
* Horizontal Scaling
* Distributed Processing

---

## Long Term Goals

* Production-grade notification system
* Event-driven architecture
* Distributed worker model
* Real-time delivery support
