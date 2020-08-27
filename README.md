# MySQL query log collector

This is little app for collect logs from mysql log files to json. This is implement like pipeline-to-pipeline formatter. Logs ready for inject to graylog or elasticsearch (via fluent, by example).

Repository staff can be not ready to go, but code is complete.

## Docker

See images here:

```
mrecco/logrotate:v1.0.0
mrecco/mysqllog-collector:v1.0.0
mrecco/fluentd:v1.3.3
```

## Configuration

### MySQL

Just redefine few variables in `my.cnf`:

```conf
general_log=ON
general-log-file=/opt/mysqllog/general.log
slow_query_log=ON
slow_query_log_file=/opt/mysqllog/slowquery.log
long_query_time=0.02
```

### Fluentd

All configured and ready to send to graylog. See **volumes/fluentd** directory.

## Solution (why this way and dont otherwise)

My code just convert ugly mysql logging formats to JSON. JSON can be pushed to any persistence provider by other applications and I satisfied how fluentd work for this.

## Reference

### Dependencies

https://github.com/hpcloud/tail
