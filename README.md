# Finance Tracker Gateway

The Finance Tracker Gateway acts as the central nexus for the Finance Tracker application, a system designed for meticulously managing expenses via a microservice architecture. Built with a modern stack including Go, gRPC, Protocol Buffers (protobuf), and SQLite3, this project embodies an efficient, scalable approach to personal finance management.

## Architecture Overview

Finance Tracker is composed of a series of microservices, with the Gateway serving as the entry point to the system. It communicates with various services, including the Expense service, through gRPC calls, ensuring high-performance and language-agnostic data exchange. Each component of the system, including the Gateway and all microservices, is containerized using Docker, facilitating seamless deployment and scalability.

### Components:

- **Finance Tracker Gateway**: The primary interface for client applications, routing requests to the appropriate services.
- **Expense Service**: A microservice responsible for managing expense records. [Repository](https://github.com/almaraz333/finance-tracker-expense-service)
- **Protocol Buffers**: The proto files defining the data structures and service interfaces used across the system are maintained in a separate repository for consistency and reusability. [Proto Files Repository](https://github.com/almaraz333/finance-tracker-proto-files)

## Getting Started

To get the entire system up and running, you'll need to clone the necessary repositories and use Docker Compose to orchestrate the containers. Here's how you can do it:

### Prerequisites

- Docker and Docker Compose installed on your machine.
- Git for cloning the repositories.

### Setup

1. **Clone the Repositories**:
    To ensure the Docker Compose setup functions correctly, clone this gateway repository and all related microservice repositories to the same directory.

    ```bash
    git clone https://github.com/almaraz333/finance-tracker-gateway.git
    git clone https://github.com/almaraz333/finance-tracker-expense-service.git
    ```

2. **Docker Compose**:
    Navigate to the root directory where you've cloned the repositories, rename the file `docker-compose.example.yaml` to `docker-compose.yaml`, and place the Docker Compose file provided in the gateway repository in the root of the directory where all the repos were cloned.

    Run the following command to start the services:

    ```bash
    docker-compose up --build
    ```

## Usage

The Finance Tracker Gateway exposes an API for interacting with the Expense service. Currently, it supports two methods on a single endpoint:

- `GET /api/expenses`: Retrieves a list of all expenses.
- `POST /api/expenses`: Creates a new expense. The request body must include a `category` (string) and an `amount` (float). The service will automatically assign an `id` and a `created_at` timestamp.
- `PUT /api/expenses/`: Updates an existing expense. The request body must include a `category` (string), an `amount` (float), and the `id` (int32).
- `DELETE /api/expenses/`: Deletes an existing expense. The request body must include the `id` (int32) of the expense to delete.

### Example: Creating an Expense

```bash
curl -X POST -d '{"category": "Utilities", "amount": 100.50}' http://127.0.0.1:8080/api/expenses
```

## Development

Should you wish to contribute or modify the Gateway or any services, remember that each component, including the Protocol Buffers, is version-controlled. Make sure to pull the latest changes and adhere to the predefined data structures and service interfaces when updating proto files.

## License

This project is licensed under the [MIT License](LICENSE.md). Feel free to fork and contribute to the development of the Finance Tracker.

## Acknowledgments

This project is built using several open-source technologies; we acknowledge and thank the creators and contributors of these projects.
