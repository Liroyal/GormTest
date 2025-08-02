# Employee API with PostgreSQL

This project provides a complete employee management API with PostgreSQL database backend, including Docker setup for local development.

## Command Reference

| Target | Description | Example |
|--------|-------------|----------|
| `make db-up` | Start PostgreSQL container | `make db-up` |
| `make db-down` | Stop and remove PostgreSQL container | `make db-down` |
| `make db-logs` | Follow PostgreSQL container logs | `make db-logs` |
| `make db-rebuild` | Stop container, remove it, start fresh, run migrations | `make db-rebuild` |
| `make migrate-up` | Apply all pending database migrations | `make migrate-up` |
| `make migrate-down` | Roll back the last database migration | `make migrate-down` |
| `make migrate-create NAME=<name>` | Create a new migration file | `make migrate-create NAME=add_employee_table` |
| `make api-build` | Build the API binary | `make api-build` |
| `make api-up` | Start the API server in background | `make api-up` |
| `make api-down` | Stop the API server | `make api-down` |
| `make help` | Show all available targets with descriptions | `make help` |

## Quick Start

Get the entire application up and running in three simple steps:

```bash
make db-up
make migrate-up
make api-up
```

That's it! The API will be available at `http://localhost:8080`.

## Configuration

### Environment Variables

The application uses the following environment variables with their defaults (defined in [`config/config.go`](config/config.go)):

**Database Configuration:**
- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `5432`)
- `DB_USER` (default: `postgres`)
- `DB_PASSWORD` (default: empty)
- `DB_NAME` (default: `gormtest`)
- `DB_SSLMODE` (default: `disable`)

**Server Configuration:**
- `SERVER_HOST` (default: `localhost`)
- `SERVER_PORT` (default: `8080`)

### Makefile Constants

The `Makefile` overrides these defaults for local development:

| Parameter | Value | Description |
|-----------|-------|-------------|
| **Docker Image** | `postgres:latest` | Official PostgreSQL Docker image |
| **Container Name** | `pg-local` | Name of the Docker container |
| **Port Mapping** | `5439` → `5432` | Host port 5439 maps to container port 5432 |
| **Username** | `admin` | Database admin username |
| **Password** | `1234` | Database admin password |
| **Database Name** | `appdb` | Default database name |

> **Note**: To modify these values, update the constants at the top of the `Makefile`. This ensures consistency across all operations.

## API Endpoints

### Health Check

Check if the API and database are running properly:

```bash
curl http://localhost:8080/health
```

**Expected Response (200 OK):**
```json
{
  "status": "ok",
  "database": "healthy",
  "message": "GormTest application is running"
}
```

**Degraded Response (503 Service Unavailable):**
```json
{
  "status": "degraded",
  "database": "unhealthy",
  "message": "Application is running but database is unavailable"
}
```

### Employee Management

#### Create Employee

Create a new employee record:

```bash
curl -X POST http://localhost:8080/employees \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Ada",
    "last_name": "Lovelace"
  }'
```

**Expected Response (201 Created):**
```json
{
  "id": 1,
  "first_name": "Ada",
  "last_name": "Lovelace",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "first_name is required"
}
```

#### Get Employee

Retrieve an employee by ID:

```bash
curl http://localhost:8080/employees/1
```

**Expected Response (200 OK):**
```json
{
  "id": 1,
  "first_name": "Ada",
  "last_name": "Lovelace",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Employee not found"
}
```

#### Update Employee

Update an existing employee:

```bash
curl -X PUT http://localhost:8080/employees/1 \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Ada",
    "last_name": "Lovelace-King"
  }'
```

**Expected Response (200 OK):**
```json
{
  "id": 1,
  "first_name": "Ada",
  "last_name": "Lovelace-King",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:35:00Z"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Employee not found"
}
```

## Development Workflow

### Starting Everything

```bash
# Start database
make db-up

# Apply migrations
make migrate-up

# Start API server
make api-up
```

### Making Changes

```bash
# Create a new migration
make migrate-create NAME=add_employee_email

# Apply new migrations
make migrate-up

# Restart API to pick up code changes
make api-down
make api-up
```

### Teardown

When you're done developing, clean up all resources:

```bash
# Stop API server
make api-down

# Stop and remove database container
make db-down
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

## Troubleshooting

### View API Logs
```bash
tail -f logs/api.log
```

### View Database Logs
```bash
make db-logs
```

### Reset Everything
```bash
make api-down
make db-rebuild
make api-up
```

## Customization

To change any configuration values:

1. Edit the constants section at the top of the `Makefile`
2. Update this README if needed
3. Run `make api-down && make db-down` to stop everything
4. Run `make db-up && make migrate-up && make api-up` to restart with new config

This approach ensures all configuration is centralized and easy to maintain.
