# trinity-proto

[![CI](https://github.com/InWamos/trinity-proto/actions/workflows/ci.yml/badge.svg)](https://github.com/InWamos/trinity-proto/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Code Quality](https://img.shields.io/badge/Code%20Quality-golangci--lint-blue)](https://golangci-lint.run/)
[![Security](https://img.shields.io/badge/Security-Gosec-orange)](https://github.com/securego/gosec)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg?logo=docker)](Dockerfile)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18.1-336791.svg?logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-8.4-DC382D.svg?logo=redis)](https://redis.io/)
[![Kafka](https://img.shields.io/badge/Apache%20Kafka-3.9.2-231F20.svg?logo=apachekafka)](https://kafka.apache.org/)

## Overview

Trinity is a production-ready backend service built with Go, following **modular monolith** principles to achieve the clarity of a monolith with the boundary discipline of microservices. Each module — `user`, `auth`, and `record` — owns its domain, application logic, infrastructure, and HTTP presentation layer in full isolation. Cross-module communication is strictly limited to well-defined client interfaces, making it trivial to extract any module into a standalone service if ever needed.

The system is **event-driven at its core**: when a user is deleted, a `UserRemovedEvent` is published to Apache Kafka and consumed asynchronously by the auth module, which immediately revokes all active sessions for that user — without either module coupled to the other.

Dependency injection is handled throughout via [Uber FX](https://github.com/uber-go/fx), enabling a clean and testable wiring of all components from infrastructure to presentation.

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                        HTTP API (chi)                        │
├───────────────┬──────────────────┬───────────────────────────┤
│  user module  │   auth module    │      record module        │
│               │                  │                           │
│  domain       │  domain          │  domain                   │
│  application  │  application     │  application              │
│  infra (pg)   │  infra (redis)   │  infra (pg)               │
└───────┬───────┴────────┬─────────┴───────────────────────────┘
        │   Kafka Events  │
        │ UserRemovedEvent│
        └────────────────►│ revoke all sessions
```

## Features

- 🔐 **Session-based authentication** — stateful sessions stored in Redis with full lifecycle management
- 👥 **Role-based authorization (RBAC)** — `User` and `Admin` roles with route-level enforcement
- 📨 **Event-driven session revocation** — user deletion triggers automatic session cleanup via Kafka
- 🏗️ **Modular monolith** — strict module boundaries with client interfaces as the only cross-module API
- 🐳 **Docker & Docker Compose** — fully containerized local development and production builds
- 📝 **End-to-end test coverage** — real containers (Postgres, Redis, Kafka) via Testcontainers
- 🔍 **Security scanning** — Gosec integrated in CI
- 📊 **Code quality** — golangci-lint enforced on every push

## HTTP API

### Users — `/v1/users`

| Method   | Path            | Description              | Auth required |
|----------|-----------------|--------------------------|---------------|
| `POST`   | `/`             | Create a new user        | Admin         |
| `GET`    | `/{id}`         | Get user by ID           | Admin         |
| `DELETE` | `/{id}`         | Delete a user            | Admin         |
| `PATCH`  | `/{id}/promote` | Promote user to Admin    | Admin         |
| `PATCH`  | `/{id}/demote`  | Demote Admin to User     | Admin         |

### Auth — `/v1/auth`

| Method | Path      | Description                        | Auth required |
|--------|-----------|------------------------------------|---------------|
| `POST` | `/login`  | Authenticate and create a session  | —             |
| `POST` | `/logout` | Terminate the current session      | User          |

### Sessions — `/v1/sessions`

| Method | Path | Description                          | Auth required |
|--------|------|--------------------------------------|---------------|
| `GET`  | `/`  | List all active sessions for a user  | User          |

### Records — `/v1/records`

| Method | Path                              | Description                          | Auth required |
|--------|-----------------------------------|--------------------------------------|---------------|
| `GET`  | `/telegram/{telegram_id}/records` | Get latest records by Telegram ID    | User          |
| `POST` | `/telegram/identity`              | Link a Telegram identity to a user   | User          |
| `POST` | `/telegram/user`                  | Register a Telegram user             | User          |
| `POST` | `/telegram/record`                | Submit a new Telegram record         | User          |