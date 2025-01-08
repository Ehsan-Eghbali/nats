# NATS JetStream Example

This repository showcases how to use **NATS** and its **JetStream** capabilities in a Hexagonal Architecture (a.k.a. Ports & Adapters). The demo includes:

- An **inbound adapter** to listen for messages published to a specific subject
- An **outbound adapter** to publish events back to NATS
- A **domain layer** demonstrating a simple Order entity and related business logic

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [How to Run](#how-to-run)
- [Testing Messages](#testing-messages)
- [Configuration](#configuration)
- [License](#license)

---

## Features

- **NATS JetStream**: Demonstrates message persistence (storage) and guaranteed delivery (retries).
- **Hexagonal Architecture**: Separates business logic (domain) from infrastructure concerns (NATS, storage, etc.).
- **Go-based**: Written in Go for high performance and simplicity.

---

## Project Structure

Below is an example layout:

```
nats/
├── cmd
│   └── main.go              # Entry point for the application
├── config
│   └── config.go            # Loads configuration (NATS URL, Stream name, etc.)
├── internal
│   ├── infra
│   │   ├── nats_listener.go # Inbound adapter that subscribes to NATS
│   │   └── nats_publisher.go# Outbound adapter to publish events
│   ├── repository
│   │   └── memory_order_repository.go
│   └── services
│       └── order_service.go
├── go.mod
├── go.sum
└── README.md                # You're reading it!
```

- **cmd/**: Contains the `main.go` file which composes all parts and starts the application.
- **config/**: Handles environment variables or default settings.
- **internal/infra/**: Infrastructure adapters for NATS (listener and publisher).
- **internal/repository/**: In-memory (or other) repositories.
- **internal/services/**: Domain services and use cases.

---

## Prerequisites

1. **Go 1.18+** (or a recent Go version).
2. **NATS Server** with JetStream enabled. You can run a local NATS instance using Docker:
   ```bash
   docker run --name nats-server -p 4222:4222 -p 8222:8222 nats:latest -js
   ```
3. **Docker** (optional, but helpful if you want to run via Docker Compose).

---

## How to Run

1. **Clone** this repository:
   ```bash
   git clone https://github.com/Ehsan-Eghbali/nats.git
   cd nats
   ```

2. **Install dependencies** (if needed) and **build**:
   ```bash
   go mod tidy
   go build -o nats-app ./cmd/main.go
   ```

3. **Run the application**:
   ```bash
   ./nats-app
   ```
   By default, it will try to connect to `nats://localhost:4222` (or whatever is set in your config).

4. **(Optional) Docker Compose**:
    - If you have a `docker-compose.yml` in the project, you can spin up NATS + the Go service together:
      ```bash
      docker-compose up --build
      ```

---

## Testing Messages

- **Publish a test message** (using the `nats` CLI or any client library):
  ```bash
  nats pub "orders.created" '{"id":"order-101","items":["itemA","itemB"]}'
  ```
    - The inbound adapter (listener) should receive this message and create an in-memory order.

- **Check logs**: You should see something like `"New order created: {order-101 ...}"`.

- **Simulated Processing**: By default, after a few seconds, the `main.go` code processes the order and publishes an `orders.processed` event. You can subscribe to it:
  ```bash
  nats sub "orders.processed"
  ```
    - The console will show the processed order data.

---

## Configuration

- **NATS_URL**: The NATS connection URL (e.g., `nats://localhost:4222`).
- **NATS_STREAM**: Name of the JetStream (e.g., `MY_STREAM`).

In the `config/config.go`, you can see how these are loaded from environment variables or default values.

---

## License

This project is open source. Feel free to use or modify it for your own needs. See the [LICENSE](LICENSE) file for details (if present).
