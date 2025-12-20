FROM valkey/valkey:7.2-alpine

COPY base.conf /config/sentinel-base.conf
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]