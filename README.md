# cyqldog

A monitoring tool for periodically executing SQLs and sending metrics to Datadog.
# Features

* Monitor multiple SQLs at different execution intervals.
* Supported data sources are Postgres (including [Redshift](https://aws.amazon.com/redshift/)) and MySQL.
* Supported notifier to send metrics is [Datadog](https://www.datadoghq.com/) (using [DogStatsD](https://docs.datadoghq.com/guides/dogstatsd/)).

# Requirements
## DogStatsD

The cyqldog uses DogStatsD to send metrics to Datadog.

Install dd-agent referring to following documents.

https://docs.datadoghq.com/guides/basic_agent_usage/

or you can run dd-agent with Docker.

https://github.com/DataDog/docker-dd-agent

or you can run only DogStatsD with Docker.

https://github.com/DataDog/docker-dd-agent#dogstatsd

When using Docker, please make sure the UDP 8125 port is open by adding the option `-p 8125:8125/udp` to the `docker run` command.

# Install
## cyqldog

Download the latest compiled binaries and put it anywhere in your executable path.

https://github.com/crowdworks/cyqldog/releases

# Usage

```bash
$ cyqldog -C /path/to/cyqldog.yml
```

# Configuration

A sample configuration looks like the following:

```yaml
data_source:
  driver: postgres
  options:
    host: {{ .DB_HOST }}
    port: 5432
    user: cyqldog
    password: {{ .DB_PASSWORD }}
    dbname: cyqldogdb
    sslmode: disable

notifiers:
  dogstatsd:
    host: {{ .DD_HOST }}
    port: 8125
    namespace: playground.cyqldog.postgres
    tags:
      - "env:local"
      - "source:{{ .DB_HOST }}"

rules:
  - name: test1
    interval: 5s
    query: "SELECT COUNT(*) AS count FROM table1"
    notifier: dogstatsd
    value_cols:
      - count
  - name: test2
    interval: 10s
    query: "SELECT tag1, val1, tag2, val2 FROM table1"
    notifier: dogstatsd
    tag_cols:
      - tag1
      - tag2
    value_cols:
      - val1
      - val2
```

# License

All code is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
