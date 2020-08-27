#!/bin/bash

export DOCKER_MYSQL_QUERYLOG=$(\
    for i in $(ls -1 /var/lib/docker/containers); \
    do \
      if [[ \
          "$(jq .Config.Image /var/lib/docker/containers/$i/config.v2.json)" \
          == \
          "\"mrecco/mysqllog-collector:v1.0.0\"" \
      ]]; \
      then \
        jq .LogPath /var/lib/docker/containers/$i/config.v2.json | tr -d '"'; \
      fi; \
    done; \
)
echo "Using ${DOCKER_MYSQL_QUERYLOG} for scrape nginx logs"
sed -ie "s#@@@DOCKER_MYSQL_QUERYLOG@@@#"$DOCKER_MYSQL_QUERYLOG"#g" /fluentd/etc/*.conf
