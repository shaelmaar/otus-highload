FROM postgres:17

RUN apt-get update && apt-get install -y \
        python3 python3-venv python3-pip curl jq gcc libpq-dev \
    && python3 -m venv /opt/venv \
    && /opt/venv/bin/pip install --no-cache-dir patroni[etcd] psycopg2-binary \
    && rm -rf /var/lib/apt/lists/*

ENV PATH="/opt/venv/bin:$PATH"
ENV PATRONI_POSTGRESQL_DATA_DIR=/var/lib/postgresql/data

RUN mkdir -p /var/lib/postgresql/data \
    && chown -R postgres:postgres /var/lib/postgresql/data \
    && chmod 700 /var/lib/postgresql/data

USER postgres

CMD ["patroni", "/etc/patroni.yaml"]