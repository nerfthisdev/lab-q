# LabQ

Telegram bot for managing subject queues in a laboratory. The bot stores data in
PostgreSQL and automatically applies database migrations at startup.

## Usage

1. **Prepare configuration.** Create a YAML file based on the structure below:

   ```yaml
   server:
     host: "0.0.0.0"
     port: ":8080"

   bot:
     token: "<telegram-token>"
     webhookURL: "https://example.com/bot"

   database:
     connectionURI: "postgres://user:pass@localhost:5432/labq?sslmode=disable"
     maxConn: 5
     maxConnLifetime: 1h
   ```

2. **Run the bot.**

   ```bash
   go run ./cmd/main.go -c path/to/config.yaml
   ```

   When starting, the application connects to the database and runs all SQL
   migrations found in the `migrations` directory.

## Development

Migrations are managed using [golang-migrate](https://github.com/golang-migrate/migrate).
New migrations can be added as numbered `*.up.sql` and `*.down.sql` files in the
`migrations` folder. They will be executed in order on startup.

## Requirements

### For admin

- Add classes
- Add dates (with weekly recurrence)

Optional:

- Edit queue
- Update the queue according to time

### For regular student

- Add yourself to the queue
- Check the queue

Optional:

- Swap places
- Check in yourself to update the queue
