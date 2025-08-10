# RARA Membership Management System

A Go CLI application for managing RARA (Radio Amateur club) membership database with point-based volunteer tracking system.

## Features

- **Member Management**: Create, update, and list amateur radio club members
- **Event Management**: Create and manage club events with point allocations
- **Volunteer Point System**: Track points earned by members for volunteering at events
- **Point Redemption**: Allow members to redeem earned points for rewards
- **Complete History**: View comprehensive transaction history for each member
- **HTTP API**: REST API server for web-based interactions
- **SQLite Database**: Lightweight, file-based database storage

## Installation

### Prerequisites

- Go 1.21 or later
- Git

### Build from Source

```bash
git clone <repository-url>
cd rara-membership
go build -o rara-membership .
```

## Quick Start

### 1. Create a Member

```bash
./rara-membership member create
# Follow the interactive prompts to add member details
```

### 2. Create an Event

```bash
./rara-membership event create
# Enter event details including default points for volunteers
```

### 3. Allocate Points for Volunteering

```bash
./rara-membership volunteer --event 1 --member 1 --points 25 --notes "Setup crew leader"
```

### 4. Redeem Points

```bash
./rara-membership redeem --member 1 --points 50 --notes "Club t-shirt"
```

### 5. View Member Point History

```bash
./rara-membership member points W1ABC
```

## Commands

### Member Management

```bash
# List all active members
./rara-membership member list

# Create a new member (interactive)
./rara-membership member create

# Update member information
./rara-membership member update W1ABC

# View complete point history for a member
./rara-membership member points W1ABC
```

### Event Management

```bash
# List all events
./rara-membership event list

# Create a new event (interactive)
./rara-membership event create

# Update an existing event
./rara-membership event update 1
```

### Point System

```bash
# Allocate points to a member for volunteering
./rara-membership volunteer --event 1 --member 1 --points 25 --notes "Setup crew"

# Redeem points (with balance validation)
./rara-membership redeem --member 1 --points 50 --notes "Club merchandise"
```

### HTTP Server

```bash
# Start the REST API server
./rara-membership serve
# Server starts on http://localhost:8080
```

## Database Schema

### Members
- **Basic Info**: name, email, phone, mailing_address
- **Radio Info**: call_sign (unique), frn, license_class
- **Status**: is_active, is_silent_key
- **License Classes**: Technician, General, Extra, Novice, Advanced

### Events
- **Details**: name, description, date, location
- **Points**: default_points_allocated, max_points_per_member

### Point System
- **PointAllocations**: Points earned by members for volunteering
- **PointDeductions**: Points redeemed by members for rewards
- **Balance Tracking**: Automatic calculation of current point balance

## Point System Rules

1. **Earning Points**: Members earn points by volunteering at events
2. **Default Points**: Each event has default points, but custom amounts can be assigned
3. **Point Limits**: Events can set maximum points per member
4. **Balance Validation**: Members cannot redeem more points than they have earned
5. **Transaction History**: Complete chronological history of all point transactions

## API Endpoints

When running the HTTP server (`./rara-membership serve`):

- `GET /members` - List all members
- `POST /members` - Create a new member
- `GET /members/{id}` - Get member details
- `PUT /members/{id}` - Update member
- `GET /events` - List all events
- `POST /events` - Create a new event
- `GET /events/{id}` - Get event details

## Database

The application uses SQLite with the database file `membership.db` created automatically in the current directory (or path specified by `RARA_DB_PATH` environment variable). The database includes:

- Automatic migrations on startup
- Indexes on frequently queried fields (call_sign, is_active)
- Foreign key constraints for data integrity
- Unique constraints to prevent duplicate allocations

## Development

### Code Structure

```
├── cmd/                 # CLI commands (Cobra framework)
├── ent/schema/         # Database schema definitions
├── services/           # Business logic layer
├── handlers/           # HTTP request handlers
├── database/           # Database connection management
├── validation/         # Input validation utilities
└── main.go            # Application entry point
```

### Key Design Patterns

1. **Schema-First**: Database schema defined in `ent/schema/`, code generated automatically
2. **Service Layer**: Business logic separated in `services/` package
3. **Thin Controllers**: CLI commands and HTTP handlers are lightweight
4. **Input Validation**: Centralized validation for emails, call signs, etc.

### Development Commands

```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run tests
go test ./...

# Generate ent code after schema changes
go generate ./ent

# Build application
go build -o rara-membership .
```

### Adding New Models

1. Create schema in `ent/schema/`
2. Run `go generate ./ent` to generate code
3. Add business logic in `services/`
4. Create CLI commands in `cmd/`
5. Add HTTP handlers if needed

## Configuration

The application uses sensible defaults and requires no configuration files. Configuration is handled through environment variables:

### Environment Variables

- **`RARA_DB_PATH`** (optional): Path to the SQLite database file
  - Default: `membership.db` (created in current working directory)
  - Example: `export RARA_DB_PATH="/path/to/custom/database.db"`
  - Used by tests to isolate database files in temporary directories

## License Classes

The system recognizes these amateur radio license classes:
- **Technician**: Entry-level license
- **General**: Most common license class
- **Extra**: Highest license class
- **Novice**: Legacy license class
- **Advanced**: Legacy license class

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For issues or questions:
- Check existing issues in the repository
- Create a new issue with detailed description
- Include steps to reproduce any bugs

---

**RARA Membership Management System** - Streamlining amateur radio club operations with modern tooling.