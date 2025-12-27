# trinity-proto

Proto for the Trinity project

[![CI](https://github.com/InWamos/trinity-proto/actions/workflows/ci.yml/badge.svg)](https://github.com/InWamos/trinity-proto/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Code Quality](https://img.shields.io/badge/Code%20Quality-golangci--lint-blue)](https://golangci-lint.run/)
[![Security](https://img.shields.io/badge/Security-Gosec-orange)](https://github.com/securego/gosec)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg?logo=docker)](Dockerfile)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18.1-336791.svg?logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-8.4-DC382D.svg?logo=redis)](https://redis.io/)

## Features

- 🔐 Session-based authentication with Redis
- 👥 Role-based authorization (User, Admin)
- 🏗️ Modular monolith architecture
- 📝 Comprehensive test coverage
- 🐳 Docker & Docker Compose support
- 🔍 Multiple security scanning tools
- 📊 Code quality and linting checks

# TODO
- [x] Session-based auth
- [x] RBAC authorization
- [x] golint and build + tests CI

# Use Cases
- User
    - [x] Get User by ID
    - [x] Create User
    - [x] Promote User
    - [x] Demote User
    - [x] Delete User

- Auth
    - [x] Login
    - [ ] Logout
    - [x] Verify
    - [ ] Logout specific session

# REFACTORING:
- [ ] Fix interactors (remove transaction logic from query interactors)
- [ ] Fix linter Errors
- [x] Rely on chi router

# Talking with the outside 
In terms of visibility, a module is allowed to import and use other modules' clients. And that's the only single piece of code they can import from the other modules. Ideally a client is defined as an interface, allowing to go with a direct code call implementation or an over-the-network implementation, in case it's needed (for instance, by an actual external application). ([Source](https://dev.to/xoubaman/modular-monolith-3fg1))