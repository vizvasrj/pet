#!/bin/bash

# Wait for the database to be ready 
# (Important if your database starts in another container)
echo "$POSTGRES_HOST $POSTGRES_PORT $POSTGRES_USER $POSTGRES_DB"
export PGPASSWORD="$POSTGRES_PASSWORD"
until psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT 1" > /dev/null 2>&1; do
  echo "Waiting for database to be ready..."
  sleep 2
done

# Run sqlx migrations
sqlx migrate run --database-url="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB"