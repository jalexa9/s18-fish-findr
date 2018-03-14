FROM postgres:9.6

RUN apt-get clean && apt-get update && apt-get install -y --no-install-recommends \
    postgresql-9.6-postgis-2.3 \
    postgis-2.3 \
&& rm -rf /var/lib/apt/lists/*

COPY init-user-db.sh /docker-entrypoint-initdb.d/init-user-db.sh
COPY . ./