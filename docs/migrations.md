## Database Migrations with golang-migrate

This project uses the [`migrate/migrate`](https://github.com/golang-migrate/migrate) CLI inside Docker to automatically apply SQL migrations located in the `migrations/` directory.

### Running inside Docker

```bash
docker compose up migrate
```

- Waits for the `postgres` service by probing the TCP port with `nc -z postgres 5432`.
- Runs `migrate -path=/migrations -database=$DATABASE_URL up`.
- Exits after finishing (status 0 when successful).

The `app` service depends on the `migrate` job, so running `docker compose up` will automatically run migrations before starting the API.

### Re-running migrations

If you add new migration files:

```bash
docker compose run --rm migrate migrate -path=/migrations -database=$DATABASE_URL up
```

To roll back the last migration:

```bash
docker compose run --rm migrate migrate -path=/migrations -database=$DATABASE_URL down 1
```

> Ensure the `DATABASE_URL` environment variable in `docker-compose.yml` matches your database credentials.
