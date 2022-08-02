#!bin/sh

# Exit immediately if any command returns a non-zero exit code
set -e

echo "Running database migration..."
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"