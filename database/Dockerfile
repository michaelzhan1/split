# Start from the official Postgres 17 image
FROM postgres:17

# Copy initialization scripts
COPY ./scripts /docker-entrypoint-initdb.d

# run
CMD ["postgres", "-p", "5432", "-c", "listen_addresses=*"]