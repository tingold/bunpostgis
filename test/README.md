## Testing
Tests have been refactored into their own package to avoid unessesary imports when using the library.

This test expects (Postgres environment variables)[https://www.postgresql.org/docs/current/libpq-envars.html] to be present to configure the database connection -- intended mostly for CI. 
