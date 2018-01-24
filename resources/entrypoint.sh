#!/bin/sh
if [ ! -e /var/db/sqlite.db ]; then
    sqlite3 /var/db/sqlite.db < /resources/sql/schema.sql
fi

exec /app "$@"