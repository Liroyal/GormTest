# PostgreSQL Docker Setup

This project provides a simple PostgreSQL Docker setup for local development.

## Configuration Constants

The following constants are defined in the `Makefile` and used throughout the project:

| Parameter | Value | Description |
|-----------|-------|-------------|
| **Docker Image** | `postgres:latest` | Official PostgreSQL Docker image |
| **Container Name** | `pg-local` | Name of the Docker container |
| **Port Mapping** | `5439` → `5432` | Host port 5439 maps to container port 5432 |
| **Username** | `admin` | Database admin username |
| **Password** | `1234` | Database admin password |
| **Database Name** | `appdb` | Default database name (optional) |

> **Note**: To modify these values, update the constants at the top of the `Makefile`. This ensures consistency across all operations.

## Quick Start

1. **Start PostgreSQL container:**
   ```bash
   make start
   ```

2. **Connect to the database:**
   ```bash
   make connect
   ```

3. **View container logs:**
   ```bash
   make logs
   ```

4. **Stop and cleanup:**
   ```bash
   make stop
   ```

## Connection Details

When the container is running, you can connect to PostgreSQL using:

- **Host**: `localhost`
- **Port**: `5439`
- **Username**: `admin`
- **Password**: `1234`
- **Database**: `appdb`

### Connection String Example
```
postgresql://admin:1234@localhost:5439/appdb
```

## Available Make Targets

Run `make help` to see all available targets with descriptions.

## Customization

To change any configuration values:

1. Edit the constants section at the top of the `Makefile`
2. Update this README if needed
3. Run `make stop` and `make start` to apply changes

This approach ensures all configuration is centralized and easy to maintain.
