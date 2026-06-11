# Event Notification System

A production-grade, scalable notification system built with Go that demonstrates modern system design principles including event-driven architecture, message queues, distributed workers, and retry mechanisms.

## 📋 Table of Contents

- [High-Level Architecture](#high-level-architecture)
- [Key Features](#key-features)
- [Development Phases](#development-phases)
- [Project Structure](#project-structure)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [System Design Interview Topics](#system-design-interview-topics)

---

## High-Level Architecture

```
              ┌──────────────┐
              │   Client     │
              │  (REST API)  │
              └──────┬───────┘
                     │
                     ▼
             ┌─────────────┐
             │ API Service │
             └──────┬──────┘
                    │ publish event
                    ▼
               ┌─────────┐
               │  Queue  │
               │ (Redis) │
               └────┬────┘
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
  Email Worker   Websocket    Retry Worker
                 Worker
```

### Flow

1. **Client Request**: Client sends notification request via REST API
2. **API Validation**: API validates and stores notification data in database
3. **Event Publishing**: API publishes event to Redis queue
4. **Worker Consumption**: Workers consume events from queue
5. **Delivery**: Workers send email/websocket notifications
6. **Failure Handling**: Failed notifications are moved to retry queue with exponential backoff

---

## Key Features

- ✅ **REST API** - Create and retrieve notifications
- ✅ **Event-Driven Architecture** - Decoupled services via message queue
- ✅ **Message Queue** - Redis Streams for reliable event distribution
- ✅ **Worker Services** - Multiple specialized worker types (Email, WebSocket, Retry)
- ✅ **Retry Logic** - Exponential backoff strategy for failed deliveries
- ✅ **Rate Limiting** - Token bucket algorithm to prevent abuse
- ✅ **WebSocket Notifications** - Real-time push notifications to connected clients
- ✅ **Email Notifications** - SendGrid/SMTP integration
- ✅ **Horizontal Scaling** - Worker pools automatically distribute load via consumer groups
- ✅ **Persistence** - PostgreSQL for notification storage and audit trails

---

## Development Phases

Follow this **exact order**. Complete each phase before moving to the next.

### Phase 1: Core API (Foundation)

**Goal**: Build the Notification API service and establish data persistence

**Features**:

- `POST /notifications` - Create notification
- `GET /notifications/:id` - Retrieve notification details
- `GET /users/:id/notifications` - List user notifications

**Data Model**:

```go
// User
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex"`
    CreatedAt time.Time
}

// Notification
type Notification struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index"`
    Type      string    // "email" or "websocket"
    Message   string
    Status    string    // "pending", "sent", "failed"
    CreatedAt time.Time
}
```

**Tech Stack**: Go, Gin, PostgreSQL, GORM

**Acceptance Criteria**:

- ✓ Store notifications in PostgreSQL
- ✓ CRUD operations functional
- ✓ Error handling implemented

---

### Phase 2: Introduce Event Queue

**Goal**: Convert system to event-driven architecture

**Changes**:

- Instead of sending notifications immediately, publish events to queue
- API becomes a producer
- Database stores event metadata

**Example Event**:

```json
{
  "notification_id": 123,
  "user_id": 88,
  "type": "email",
  "message": "Payment successful"
}
```

**Queue Technology**: Redis Streams

**Acceptance Criteria**:

- ✓ API publishes events to Redis
- ✓ Queue stores events reliably
- ✓ Event schema defined

---

### Phase 3: Worker Service

**Goal**: Create worker service to consume queue events

**Implementation**:

```go
// Pseudocode
while true {
    message := redisStream.ReadGroup("notification-group")
    if message != nil {
        processNotification(message)
        redisStream.Acknowledge(message)
    }
}
```

**Consumer Pattern**: XREADGROUP for distributed consumption

**Acceptance Criteria**:

- ✓ Worker consumes Redis queue messages
- ✓ Consumer groups distribute load
- ✓ Message acknowledgment implemented

---

### Phase 4: Email Service

**Goal**: Implement email delivery via worker

**Features**:

- Email worker consumes queue events
- Uses SMTP or SendGrid provider
- Updates notification status in database

**Flow**:

```
Queue Event
    ↓
Email Worker
    ↓
Send Email
    ↓
Update Notification Status
```

**Implementation File**: `worker/email_worker.go`

**Acceptance Criteria**:

- ✓ Email worker consumes events
- ✓ Email delivery functional
- ✓ Database status updates

---

### Phase 5: WebSocket Notifications

**Goal**: Add real-time notification delivery

**Architecture**:

```
Client connects: ws://server/ws?user_id=123
         ↓
Server stores: user_id → connection
         ↓
Worker sends message to WebSocket connection
         ↓
Client receives in real-time
```

**Features**:

- Real-time push notifications
- Per-user subscriptions
- Connection pooling

**Acceptance Criteria**:

- ✓ WebSocket connections established
- ✓ Real-time message delivery
- ✓ Connection cleanup on disconnect

---

### Phase 6: Retry System

**Goal**: Implement robust failure recovery

**Strategy**: Exponential backoff with retry limits

**Retry Schedule**:

- 1st retry: 2 seconds
- 2nd retry: 4 seconds
- 3rd retry: 8 seconds
- 4th retry: 16 seconds
- 5th retry: 32 seconds

**Storage**: Database tracks retry attempts per notification

**Flow**:

```
Email Failed
    ↓
Move to Retry Queue
    ↓
Retry Worker Picks Up
    ↓
Re-attempt Delivery
```

**Acceptance Criteria**:

- ✓ Retry queue functional
- ✓ Exponential backoff implemented
- ✓ Attempt tracking in database

---

### Phase 7: Rate Limiting

**Goal**: Prevent abuse and system overload

**Strategy**: Token bucket algorithm

**Example Configuration**:

- 10 notifications per minute per user
- Sliding window tracking

**Implementation**:

```go
// Using golang.org/x/time/rate
limiter := rate.NewLimiter(rate.Limit(10)/rate.Minute, 10)
if !limiter.Allow() {
    return errors.New("rate limit exceeded")
}
```

**Storage**: Redis for distributed token management

**Acceptance Criteria**:

- ✓ Rate limiting enforced
- ✓ Distributed across instances
- ✓ User-friendly error responses

---

### Phase 8: Worker Scaling

**Goal**: Scale workers horizontally

**Architecture**:

```
Worker 1  ─┐
Worker 2  ─┤
Worker 3  ├─→ Redis Consumer Group
Worker 4  ─┤
Worker N  ─┘
```

**Key Concept**: Redis consumer groups automatically distribute messages among connected workers. Each message is processed by exactly one worker.

**Benefits**:

- Automatic load distribution
- Fault tolerance (worker crash doesn't lose messages)
- Easy horizontal scaling

**Scaling Steps**:

1. Ensure consumer group configured
2. Start new worker instances
3. Workers automatically join consumer group
4. Redis distributes load evenly

**Interview Point**: "Workers scale horizontally with consumer groups. When a new worker joins, Redis automatically rebalances message distribution."

**Acceptance Criteria**:

- ✓ Multiple workers running
- ✓ Load distributed evenly
- ✓ No message loss on worker failure

---

## Project Structure

```
edns/
│
├── README.md
├── LICENSE
├── Makefile                    # Build targets and commands
├── docker-compose.yml          # Production Docker setup
├── docker-compose.dev.yml      # Development Docker setup
├── .gitignore
│
├── api-service/
│   ├── main.go                 # Entry point
│   ├── config.go               # Configuration loading
│   ├── handlers/
│   │   ├── notification.go     # REST endpoint handlers
│   │   └── health.go           # Health check endpoint
│   ├── services/
│   │   ├── notification.go     # Business logic
│   │   └── queue.go            # Queue publishing
│   ├── models/
│   │   ├── user.go             # User model
│   │   └── notification.go     # Notification model
│   └── repository/
│       └── notification.go     # Database queries
│
├── worker/
│   ├── main.go                 # Worker entry point
│   ├── config.go               # Configuration
│   ├── email_worker.go         # Email delivery logic
│   ├── websocket_worker.go     # WebSocket delivery
│   └── retry_worker.go         # Retry logic
│
├── websocket-server/
│   ├── main.go                 # WebSocket entry point
│   └── server.go               # Connection management
│
├── pkg/
│   ├── queue/
│   │   └── redis.go            # Queue operations
│   ├── logger/
│   │   └── zap.go              # Structured logging
│   ├── config/
│   │   └── config.go           # Config management
│   └── models/
│       └── event.go            # Event schema
│
└── migrations/
    └── 001_initial_schema.sql  # Database schema
```

---

## Tech Stack

### Framework & HTTP

- **Gin** - Lightweight HTTP framework
  - Fast routing
  - Middleware support
  - Built-in error handling

### Database

- **PostgreSQL** - Primary data store
- **GORM** - ORM for database operations
- **pgx** - PostgreSQL driver

### Queue & Caching

- **Redis** - Message queue and caching
- **go-redis** - Redis client library
  - Used for Streams
  - Consumer group operations

### Real-Time Communication

- **Gorilla WebSocket** - WebSocket library
  - Full-duplex communication
  - Connection pooling

### Email

- **Gomail** or **SendGrid** - Email delivery
  - Template support
  - Batch sending

### Logging

- **Uber Zap** - Structured logging
  - High performance
  - Contextual logging

### Configuration

- **Viper** - Configuration management
  - Environment variables
  - Config file support
  - Type casting

### Rate Limiting

- **golang.org/x/time/rate** - Token bucket implementation

### Testing

- **Go testing package** - Built-in
- **Testify** - Assertions and mocking

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Redis 6.0+
- Docker & Docker Compose (optional)

### Installation

1. **Clone repository**

   ```bash
   git clone https://github.com/yourusername/edns.git
   cd edns
   ```

2. **Install dependencies**

   ```bash
   cd api-service && go mod download
   cd ../worker-service && go mod download
   cd ../websocket-service && go mod download
   ```

3. **Set up environment**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start services using Docker Compose**

   ```bash
   make docker-up-dev
   ```

   Or manually with docker-compose:

   ```bash
   docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
   ```

5. **Run database migrations**
   ```bash
   cd api-service
   go run cmd/server/main.go migrate
   ```

### Running Services

**Using Makefile**:

```bash
# Run API service
make run-api
# Runs on :8080

# Build API service binary
make build-api
# Output: api-service/bin/notification-api

# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Vet code
make vet

# Clean build artifacts
make clean

# Docker - Start services
make docker-up-dev

# Docker - Stop services
make docker-down-dev
```

**Manual Execution**:

```bash
# API Service
cd api-service
go run cmd/server/main.go
# Runs on :8080

# Worker Service
cd worker-service
go run cmd/worker/main.go
# Processes queue events

# WebSocket Service
cd websocket-service
go run cmd/websocket/main.go
# Runs on :8081
```

### API Endpoints

**Create Notification**:

```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "type": "email",
    "message": "Payment successful"
  }'
```

**Get Notification**:

```bash
curl http://localhost:8080/notifications/1
```

**Get User Notifications**:

```bash
curl http://localhost:8080/users/1/notifications
```

**WebSocket Connection**:

```javascript
const ws = new WebSocket("ws://localhost:8081/ws?user_id=123");
ws.onmessage = (event) => {
  console.log("Notification:", event.data);
};
```

---

## Example System Walkthrough

### Create and Deliver Email Notification

1. **Client**: POST `/notifications`

   ```json
   {
     "user_id": 42,
     "type": "email",
     "message": "Your order #123 has shipped"
   }
   ```

2. **API Service**:
   - Validates request
   - Creates notification record (status: "pending")
   - Publishes event to Redis: `notification:create`

3. **Queue** (Redis Stream):

   ```
   {
     "notification_id": 789,
     "user_id": 42,
     "type": "email",
     "message": "Your order #123 has shipped"
   }
   ```

4. **Email Worker**:
   - Consumes event from consumer group
   - Looks up user email from database
   - Sends email via SMTP/SendGrid
   - Updates notification status → "sent"
   - Acknowledges message in Redis

5. **Success**: Client can poll `/notifications/789` and see status "sent"

### Failure & Retry Flow

1. **Email Worker**: Email send fails (network error)
   - Updates notification status → "failed"
   - Publishes retry event with `attempt: 1`

2. **Retry Worker**:
   - Consumes retry event
   - Calculates backoff: 2^1 = 2 seconds
   - Schedules re-delivery in 2 seconds

3. **Retry Worker** (after 2 seconds):
   - Re-attempts email send
   - Success: status → "sent"
   - Failure: status → "failed", attempts → 2

4. **Max Attempts**: After 5 failures, moved to dead-letter queue for manual review

---

## Configuration

Environment variables (`.env`):

```env
# API Service
API_PORT=8080
API_ENV=development

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=notifications

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=app-password

# Rate Limiting
RATE_LIMIT_PER_MINUTE=10

# Logging
LOG_LEVEL=info
```

---

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./api-service/...

# Run with verbose output
go test -v ./...
```

---

## Deployment

### Docker Compose (Local Development)

```bash
docker-compose up -d
```

### Production Deployment

1. **Build images**:

   ```bash
   docker build -t notification-api ./api-service
   docker build -t notification-worker ./worker
   docker build -t notification-ws ./websocket-server
   ```

2. **Use Kubernetes** (optional):
   - Deploy API service with multiple replicas
   - Deploy workers with horizontal pod autoscaling
   - Persistent volumes for PostgreSQL

3. **Environment-specific configs**:
   - Separate `.env.prod`, `.env.staging`
   - Use secrets management (HashiCorp Vault, AWS Secrets)

---

## Monitoring & Observability

### Metrics to Track

- API response times
- Queue depth / consumer lag
- Worker processing rate
- Failure rates (by type)
- Retry attempts distribution
- WebSocket connection count
- Rate limit violations

### Logging

Structured logs help debugging:

```json
{
  "timestamp": "2026-03-09T10:30:45Z",
  "level": "error",
  "notification_id": 789,
  "user_id": 42,
  "error": "email delivery failed",
  "reason": "SMTP connection timeout"
}
```

### Health Checks

All services expose `/health` endpoint:

- API: Checks DB connectivity
- Worker: Checks Redis connectivity
- WebSocket: Checks Redis pub/sub

---

## Contributing

1. Create feature branch: `git checkout -b feature/new-feature`
2. Commit changes: `git commit -am 'Add feature'`
3. Push branch: `git push origin feature/new-feature`
4. Create Pull Request

---

## License

Licensed under MIT License. See [LICENSE](LICENSE) file for details.

---
