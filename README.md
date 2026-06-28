# Notification Processor Stateless Service

A production-ready **Go microservice** that consumes user action events from Kafka, builds personalized notifications, and guarantees reliable delivery through an external service.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Notification Routing Rules](#notification-routing-rules)
- [Kafka Message Format](#kafka-message-format)
- [Guarantees](#guarantees)
- [Requirements](#requirements)
- [Failure Simulation](#failure-simulation)
- [Getting Started](#getting-started)

---

## Overview

The **Notification Processor Service** is designed to:

- Consume user action events from Kafka
- Build personalized notifications based on event type
- Guarantee delivery through an external notification service
- Handle failures gracefully with retry mechanisms
- Ensure exactly-once processing semantics

---

## Features

| Feature | Description |
| :--- | :--- |
| 📨 **Kafka Integration** | Consumes events with manual offset management (autocommit disabled) |
| 📧 **Smart Routing** | Routes notifications to appropriate channels (Email, Push, SMS) |
| 💾 **PostgreSQL Storage** | Stores notification state and enables deduplication |
| 🔄 **Retry Logic** | Exponential backoff with max 3 retries for temporary failures |
| 🛡️ **Idempotency** | Duplicate messages don't create duplicate notifications |
| ⚡ **Concurrent Processing** | Handles messages in parallel without deadlocks |
| 🧹 **Graceful Shutdown** | Cleanly completes processing on SIGTERM |
| 🚫 **Poison Pill Resilience** | Malformed JSON doesn't block other messages |
| 📊 **No Memory Leaks** | Stable memory usage even under continuous failures |

---

## Notification Routing Rules

The service routes notifications based on event type:

| Event Type | Delivery Channels |
| :--- | :--- |
| `order_created` | 📧 **Email** + 📱 **Push** |
| `payment_received` | 📧 **Email** |
| `order_shipped` | 📱 **Push** + 💬 **SMS** |

> **Note:** Actual sending is not required — simply log who the notification was sent to and via which channel.

---

## Kafka Message Format

### Sample Message

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": 12345,
  "event_type": "order_created",
  "timestamp": "2024-01-15T14:30:00Z",
  "payload": {
    "order_id": "550e8400-e29b-41d4-a716-446655440000",
  }
}
```

### Field Descriptions

| Field | Type | Description |
| :--- | :--- | :--- |
| `event_id` | `string` | Unique event identifier (UUID) |
| `user_id` | `int64` | ID of the user receiving the notification |
| `event_type` | `string` | Type of event (one of 3 types) |
| `timestamp` | `string` | Event timestamp in RFC3339 format |
| `payload` | `object` | Arbitrary JSON payload with event details |

### Supported Event Types

- `order_created`
- `payment_received`
- `order_shipped`

---

## Guarantees

### 1. Idempotency

> Duplicate messages with the same `event_id` do not create duplicate notifications.

**Implementation:** PostgreSQL `UNIQUE` constraint on `event_id`.

### 2. Exactly-Once Semantics

> Messages are not lost if the service crashes between DB write and offset commit.

**Implementation:** Transactional write + manual offset commit after successful persistence.

### 3. Fault Tolerance

> Poison pill (malformed JSON) does not block processing of other messages.

**Implementation:** Per-message error handling with continue on failure.

### 4. No Memory Leaks

> Service runs for hours without memory growth, even under constant temporary delivery failures.

**Implementation:** Proper resource cleanup, bounded worker pools, and goroutine management.

### 5. Graceful Shutdown

> Processing of current messages completes cleanly upon SIGTERM.

**Implementation:** Signal handling with context cancellation and sync.WaitGroup.

---

## Requirements

| Requirement | Details |
| :--- | :--- |
| **Database** | PostgreSQL for state storage and deduplication |
| **Message Broker** | Kafka with **manual offset management** (autocommit prohibited) |
| **Concurrency** | Concurrent message processing without deadlocks |
| **Retry Logic** | Max 3 retries with exponential backoff |
| **Language** | Go (must compile and run with `go run`) |

---

## Failure Simulation

The delivery simulation function must emulate:

- **10%** of calls → **Temporary errors** (retry with exponential backoff)
- **1%** of calls → **Permanent errors** (stop after 3 attempts)
- **89%** of calls → **Success** (log delivery)

---

## Getting Started

### Prerequisites

- Go 1.26.4
- Docker & Docker Compose
- Kafka (or use provided docker-compose)
- PostgreSQL (or use provided docker-compose)

### Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/notification-processor.git
cd notification-processor

# 2. Start dependencies (Kafka + PostgreSQL)
docker-compose up -d

# 3. Run the service
go run cmd/processor/main.go

# 4. (Optional) Send test events
go run cmd/generator/main.go
```