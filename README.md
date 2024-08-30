# habitmap

Heatmap esque webapp for habit tracking

## Local Development

Some helpful commands

```zsh
# Start database
docker-compose up -d

# Initial database setup
docker-compose exec -T postgres psql -U myuser -d habitmap -f scripts/db_init.sql
docker-compose exec -T postgres psql -U myuser -d habitmap -f scripts/populate_tables.sql

# Remove database
docker-compose down -v

# Access database from the terminal
docker-compose exec postgres psql -U myuser -d habitmap
```
