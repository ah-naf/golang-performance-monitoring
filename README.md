# Go Metrics Monitoring System

This project is a Go-based application designed for task management, instrumented with Prometheus for metrics collection and Grafana for visualization. The entire system is containerized using Docker and Docker Compose for easy setup and deployment. It also includes a k6 script for load testing the API.

## Technologies Used

- **Go:** Backend API development using the Gin framework.
- **PostgreSQL:** Database for storing task information.
- **Docker & Docker Compose:** For containerizing and orchestrating the application services.
- **Prometheus:** For collecting application and system metrics.
- **Grafana:** For visualizing metrics collected by Prometheus.
- **k6:** For load testing the API.

## Folder Structure

```
.
├── docker-compose.yml      # Docker Compose file to orchestrate all services.
├── go-app/                 # Contains the Go application source code.
│   ├── Dockerfile          # Dockerfile for building the Go application image.
│   ├── go.mod              # Go module definition.
│   ├── go.sum              # Go module checksums.
│   ├── init.sql            # SQL script to initialize the PostgreSQL database.
│   ├── internals/          # Internal Go packages.
│   │   ├── database/       # Database connection and query logic.
│   │   ├── metrics/        # Prometheus metrics collection logic.
│   │   ├── models/         # Go struct definitions for data models (e.g., Task).
│   │   └── routers/        # API route handlers.
│   └── main.go             # Main entry point for the Go application.
├── grafana/                # Grafana provisioning files.
│   └── provisioning/
│       ├── dashboard_providers/ # Defines where to find dashboard definitions.
│       │   └── providers.yml
│       ├── dashboards/         # Grafana dashboard JSON definitions.
│       │   └── grafana_dashboard.json
│       └── datasources/        # Grafana datasource definitions (e.g., Prometheus).
│           └── prometheus.yml
├── load-test.js            # k6 script for load testing the API.
├── prometheus.yml          # Prometheus configuration file.
└── README.md               # This file.
```

## Running the Project

**Prerequisites:**

- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Docker Compose: Usually comes with Docker Desktop. If not, [Install Docker Compose](https://docs.docker.com/compose/install/)

**Steps:**

1.  **Clone the repository (if you haven't already):**

    ```bash
    git clone https://github.com/ah-naf/golang-performance-monitoring
    cd golang-performance-monitoring
    ```

2.  **Start all services using Docker Compose:**
    Open a terminal in the project's root directory (where `docker-compose.yml` is located) and run:

    ```bash
    docker-compose up -d
    ```

    This command will build the Go application image (if not already built) and start all the services (PostgreSQL, Go App, Prometheus, Grafana) in detached mode (`-d`).

3.  **To stop the services:**
    ```bash
    docker-compose down
    ```

## Accessing Services

Once the services are running, you can access them at the following URLs:

- **Go API Server:** `http://localhost:8080`
  - This is the main application providing task management functionalities.
- **Prometheus:** `http://localhost:9090`
  - Use this to view collected metrics, explore metric data, and check target status.
- **Grafana:** `http://localhost:3000`
  - Login with default credentials:
    - Username: `admin`
    - Password: `admin`
  - A pre-configured Prometheus datasource and a sample dashboard should be available.

## Load Testing

The project includes a k6 script (`load-test.js`) to simulate user traffic and test the performance of the API.

**Prerequisites:**

- **k6:** [Install k6](https://k6.io/docs/getting-started/installation/)

**Running the Load Test:**

1.  Make sure the application stack is running (see "Running the Project").
2.  Open a terminal in the project's root directory.
3.  Run the following command:

    ```bash
    k6 run load-test.js
    ```

    This will execute the script, which simulates creating, reading, updating, and deleting tasks. The script is configured with stages to gradually ramp up virtual users.

    You can modify `load-test.js` to change the test parameters (e.g., duration, number of virtual users, API endpoints).
    The `BASE_URL` for the test is set to `http://localhost:8080` by default but can be overridden using an environment variable:

    ```bash
    k6 run -e BASE_URL=http://your-custom-host:8080 load-test.js
    ```

## API Endpoints (Go Server)

The Go API server runs on `http://localhost:8080`.

### General

- **`GET /`**

  - **Description:** Returns a welcome message for the Go Metric Monitoring System.
  - **Response (200 OK):**
    ```
    Golang Metric Monitoring System
    ```

- **`GET /health`**

  - **Description:** Performs a health check of the application. Used by Docker to verify service health.
  - **Response (200 OK):**
    ```json
    {
      "status": "healthy",
      "message": "Service is running smoothly"
    }
    ```
  - **Response (503 Service Unavailable if DB is down):**
    ```json
    {
      "status": "unhealthy",
      "message": "Database connection failed: <error details>"
    }
    ```

- **`GET /metrics`**
  - **Description:** Exposes application metrics in Prometheus format.
  - **Response (200 OK):** Text-based Prometheus metrics.

### Task Management (`/task`)

- **`POST /task`**

  - **Description:** Adds a new task.
  - **Request Body (application/json):**
    ```json
    {
      "id": 123,
      "title": "My New Task",
      "description": "Details about the task.",
      "status": "pending"
    }
    ```
    - `id` (integer, required): Unique identifier for the task.
    - `title` (string, required): Title of the task.
    - `description` (string, optional): Detailed description of the task.
    - `status` (string, required): Current status of the task (e.g., "pending", "in-progress", "completed").
  - **Response (201 Created):**
    ```json
    {
      "success": "Task added successfully",
      "task": {
        "id": 123,
        "title": "My New Task",
        "description": "Details about the task.",
        "status": "pending"
      }
    }
    ```
  - **Response (400 Bad Request):** If the request body is invalid or a database error occurs (e.g., duplicate ID).
    ```json
    {
      "error": "<error details>"
    }
    ```

- **`GET /task`**

  - **Description:** Retrieves all tasks.
  - **Response (200 OK):**
    ```json
    {
      "tasks": [
        {
          "id": 123,
          "title": "My New Task",
          "description": "Details about the task.",
          "status": "pending"
        },
        {
          "id": 124,
          "title": "Another Task",
          "description": "More details.",
          "status": "completed"
        }
      ]
    }
    ```
    If no tasks exist, `tasks` will be an empty array.
  - **Response (400 Bad Request):** If a database error occurs.

- **`GET /task/:id`**

  - **Description:** Retrieves a specific task by its ID.
  - **Path Parameter:**
    - `id` (integer): The ID of the task to retrieve.
  - **Response (200 OK):**
    ```json
    {
      "task": {
        "id": 123,
        "title": "My New Task",
        "description": "Details about the task.",
        "status": "pending"
      }
    }
    ```
  - **Response (400 Bad Request):** If the `id` is not a valid integer.
    ```json
    {
      "error": "invalid task ID"
    }
    ```
  - **Response (404 Not Found):** If no task with the given ID exists.
    ```json
    {
      "error": "task not found"
    }
    ```
  - **Response (500 Internal Server Error):** For other database errors.

- **`PUT /task/:id`**

  - **Description:** Updates an existing task by its ID.
  - **Path Parameter:**
    - `id` (integer): The ID of the task to update.
  - **Request Body (application/json):**
    ```json
    {
      "title": "Updated Task Title",
      "description": "Updated description.",
      "status": "in-progress"
    }
    ```
    - `title` (string, required)
    - `description` (string, optional)
    - `status` (string, required)
      _(Note: The `id` field in the request body is ignored; the `id` from the URL path is used.)_
  - **Response (200 OK):**
    ```json
    {
      "success": "task updated successfully",
      "task": {
        "id": 123, // The ID from the path parameter
        "title": "Updated Task Title",
        "description": "Updated description.",
        "status": "in-progress"
      }
    }
    ```
  - **Response (400 Bad Request):** If the `id` is invalid or the request body is malformed.
  - **Response (404 Not Found):** If no task with the given ID exists.
  - **Response (500 Internal Server Error):** For other database errors.

- **`DELETE /task/:id`**
  - **Description:** Deletes a task by its ID.
  - **Path Parameter:**
    - `id` (integer): The ID of the task to delete.
  - **Response (200 OK):**
    ```json
    {
      "message": "task deleted successfully"
    }
    ```
  - **Response (400 Bad Request):** If the `id` is not a valid integer.
  - **Response (404 Not Found):** If no task with the given ID exists.
  - **Response (500 Internal Server Error):** For other database errors.
