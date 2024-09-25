# 🎟️ Tickets API

A simple, efficient API for managing ticket operations, including creation, retrieval, and purchases.

## ✨ Features

- **Create Tickets**: Easily create a new ticket.
- **Retrieve Tickets**: Fetch the details of an existing ticket.
- **Purchase Tickets**: Facilitate the purchase of tickets.
- **Swagger Documentation**: Fully documented API with Swagger for easier integration.

## 🛠️ Technologies Used

- **Web Framework**: [Echo v4](https://echo.labstack.com/) - Fast, minimalist Go web framework.
- **Logging**: [ZapLogger](https://github.com/uber-go/zap) - High-performance logging.
- **Database**: [PostgreSQL](https://www.postgresql.org/) - Reliable relational database system.
- **Configuration**: `yaml` & [viper](https://github.com/spf13/viper) for environment mapping and configuration management.

## 🔗 API Endpoints

### 🎫 Tickets

- `POST /tickets` - **Create a new ticket**  
- `GET /tickets/:id` - **Retrieve ticket details** by ticket ID  
- `POST /tickets/:id/purchases` - **Purchase a ticket** by ticket ID

## 📜 Swagger Documentation

Access the interactive Swagger documentation for a full overview of the API at:

```sh
/swagger/index.html
```

## 🚀 Installation

### Prerequisites

Ensure you have the following tools installed on your system:

- Docker: For containerization.
- Go (version 1.22.3 or later): For building and running the application.

### 📦 Building and Running the API

####  Using Docker Compose

1. Clone the repository:

   ```sh
   git clone https://github.com/fleimkeipa/tickets-api.git
   cd tickets-api
   ```

2. Copy the example configuration file and modify it if necessary:

   ```sh
   cp config_example.yaml config.yaml
   ```

3. Build and run the application using Docker Compose:

   ```sh
   docker-compose up --build -d
   ```

The application will now run in a Docker container and be accessible at the configured port!

#### Running Locally

If you prefer to run the application locally without Docker, follow these steps:

1. Install dependencies:

   ```sh
   go mod download
   ```

2. Build the application:

   ```sh
   go build
   ```

3. Run the application:

   ```sh
   ./tickets-api
   ```

Now the API should be running and accessible at your configured port!
