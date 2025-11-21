FROM mongo:8-noble

COPY --chown=mongodb:mongodb init.sh /usr/local/bin/
RUN chmod u+x /usr/local/bin/init.sh

ENTRYPOINT ["/usr/local/bin/init.sh"]