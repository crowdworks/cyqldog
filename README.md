# cyqldog
[![CircleCI](https://circleci.com/gh/crowdworks/cyqldog.svg?style=svg)](https://circleci.com/gh/crowdworks/cyqldog)
[![GitHub release](http://img.shields.io/github/release/crowdworks/cyqldog.svg?style=flat-square)](https://github.com/crowdworks/cyqldog/releases)
[![GoDoc](https://godoc.org/github.com/crowdworks/cyqldog/cyqldog?status.svg)](https://godoc.org/github.com/crowdworks/cyqldog/cyqldog)

A monitoring tool for periodically executing SQLs and sending metrics to Datadog.

# Features

* Execute multiple SQLs at different intervals and send metrics.
* Supported data sources are PostgreSQL (including Redshift) and MySQL.
* Supported notifier to send metrics is Datadog (using DogStatsD).

# Requirements
## DogStatsD

The cyqldog uses [DogStatsD](https://docs.datadoghq.com/guides/dogstatsd/) to send metrics to Datadog.

DogStatsD is a metrics aggregation service bundled with the Datadog Agent (datadog-agent).
If you want to send metrics to Datadog, you may already have installed datadog-agent including DogStatsD,
but if you have not installed datadog-agent,

* install Datadog Agent (https://docs.datadoghq.com/agent/)
* or you can run Datadog Agent with Docker. (https://hub.docker.com/r/datadog/agent)
* or you can run only DogStatsD with Docker. (https://hub.docker.com/r/datadog/dogstatsd)

When using Docker, please make sure the UDP 8125 port is open by adding the option `-p 8125:8125/udp` to the `docker run` command.

# Install
Download the latest compiled binaries and put it anywhere in your executable path.

https://github.com/crowdworks/cyqldog/releases

or you can run with Docker. (described later)

# Usage
Create your configuration file, and pass it to `-C` flag to run cyqldog.

```bash
$ cyqldog -C /path/to/cyqldog.yml
```

# Configuration

An example for cyqldog.yml as follows:

```yaml
# An example for cyqldog.yml

# Data source is a configuration of the database to connect.
data_source:
  # Driver is type of the database to connect.
  # Currently suppoted databases are as follows:
  #  - postgres
  #  - mysql
  #
  # An example for Postgres (including Redshift)
  #
  driver: postgres
  # Options is a map of options to connect.
  # These options are passed to sql.Open.
  # The supported options are depend on the database driver.
  options:
    # For all supported options, see godoc in lib/pq.
    # https://godoc.org/github.com/lib/pq
    #
    # The host to connect to. Values that start with / are for unix domain sockets. (default is localhost)
    host: db.example.com
    # The port to bind to. (default is 5432)
    port: 5432
    # The user to sign in as
    user: cyqldog
    # The user's password
    # Sensitive values such as a password can be read from environment variables.
    # In this example, an environment variable named DB_PASSWORD is set.
    password: {{ .DB_PASSWORD }}
    # The name of the database to connect to
    dbname: cyqldogdb
    # Whether or not to use SSL (default is require, this is not the default for libpq)
    # * disable - No SSL
    # * require - Always SSL (skip verification)
    # * verify-ca - Always SSL (verify that the certificate presented by the server was signed by a trusted CA)
    # * verify-full - Always SSL (verify that the certification presented by the server was signed by a trusted CA and the server host name matches the one in the certificate)
    sslmode: disable

  # An example for MySQL
  #
  # Note that options other than host, port, user, password, dbname are passed as connection parameters.
  # For all supported connection parameters, see README in go-sql-driver/mysql.
  # https://github.com/go-sql-driver/mysql#parameters
  #
  # driver: mysql
  # options:
  #   host: db.example.com
  #   port: 3306
  #   user: cyqldog
  #   password: {{ .DB_PASSWORD }}
  #   dbname: cyqldogdb
  #   charset: "utf8"
  #   collation: "utf8_general_ci"

# Notifiers are configurations of output plugins.
notifiers:
  # Dogstatsd is a configuration of the dogstatsd to connect.
  dogstatsd:
    # Host is a hostname or IP address of the dogstatsd.
    host: {{ .DD_HOST }}
    # Port is a port number of the dogstatsd.
    port: 8125
    # Namespace to prepend to all statsd calls
    namespace: playground.cyqldog
    # Tags are global tags to be added to every statsd call
    tags:
      - "env:local"
      - "source:db.example.com"

# Rules are a list of rules to monitor.
rules:
  # Name of the rule.
  # The metic name sent to the dogstatsd is:
  #   Dogstatsd.Namespace + Rule.Name + Rule.ValueCols[*]
  - name: test1
    # Interval of the monitoring.
    # Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
    interval: 5s
    # Query to the database.
    query: "SELECT COUNT(*) AS count FROM table1"
    # Notifier is a name of notifier to send metrics.
    notifier: dogstatsd
    # ValueCols is a list of names of the columns used as metric values.
    # In this example, the following metrics are sent.
    # * playground.cyqldog.test1.count (with tags ["env:local", "source:db.example.com"])
    value_cols:
      - count
  - name: test2
    interval: 10s
    query: "SELECT tag1, val1, tag2, val2 FROM table1"
    notifier: dogstatsd
    # TagCols is a list of names of the columns used as metric tags.
    # In this example, the following metrics are sent.
    # * playground.cyqldog.test2.value1 (with tags ["env:local", "source:db.example.com", "tag1:(value of column tag1)", "tag2:(value of column tag2)"])
    # * playground.cyqldog.test2.value2 (with tags ["env:local", "source:db.example.com", "tag1:(value of column tag1)", "tag2:(value of column tag2)"])
    tag_cols:
      - tag1
      - tag2
    value_cols:
      - val1
      - val2
```

# Run with Docker
You can run cyqldog with Docker.

Start test database and dogstatsd with docker compose for demo.

## Postgres

```
$ docker compose up -d

$ docker run -it --rm \
    -v $(pwd)/cyqldog/test-fixtures/postgres:/app/config \
    -e PGPASSWORD=password \
    postgres:14.1-alpine \
    psql -h host.docker.internal -U cyqldog -d cyqldogdb -f /app/config/setup_dev.sql
```

The Docker images are available in the GitHub Container Registry.

https://github.com/crowdworks/cyqldog/pkgs/container/cyqldog

Mount a sample configuration file, set environment variables, and run cyqldog.

```
$ docker run -it --rm \
    -v $(pwd)/cyqldog/test-fixtures/postgres:/app/config \
    -e DB_HOST=host.docker.internal \
    -e DB_PASSWORD=password \
    -e DD_HOST=host.docker.internal \
    -e DD_API_KEY=[YOUR_DD_API_KEY] \
    ghcr.io/crowdworks/cyqldog:latest -C config/cyqldog.yml
```

## MySQL

```
$ docker compose up -d

$ docker run -it --rm \
    -v $(pwd)/cyqldog/test-fixtures/mysql:/app/config \
    mysql:8.0 \
    /bin/bash -c "mysql -h host.docker.internal -u cyqldog -ppassword cyqldogdb < /app/config/setup_dev.sql"
```

The Docker images are available in the GitHub Container Registry.

https://github.com/crowdworks/cyqldog/pkgs/container/cyqldog

Mount a sample configuration file, set environment variables, and run cyqldog.

```
$ docker run -it --rm \
    -v $(pwd)/cyqldog/test-fixtures/mysql:/app/config \
    -e DB_HOST \
    -e DB_PASSWORD \
    -e DD_HOST \
    -e DD_API_KEY \
    ghcr.io/crowdworks/cyqldog:latest -C config/cyqldog.yml
```

# Contributions

Any feedback and contributions are welcome. Feel free to open an issue and submit a pull request.

# License

The cyqldog is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
