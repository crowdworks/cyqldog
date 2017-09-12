# cyqldog
[![CircleCI](https://circleci.com/gh/crowdworks/cyqldog.svg?style=svg)](https://circleci.com/gh/crowdworks/cyqldog)
[![GitHub release](http://img.shields.io/github/release/crowdworks/cyqldog.svg?style=flat-square)](https://github.com/crowdworks/cyqldog/releases)
[![GoDoc](https://godoc.org/github.com/crowdworks/cyqldog/deploy?status.svg)](https://godoc.org/github.com/crowdworks/cyqldog/deploy)

A monitoring tool for periodically executing SQLs and sending metrics to Datadog.

# Features

* Execute multiple SQLs at different intervals and send metrics.
* Supported data sources are Postgres (including Redshift) and MySQL.
* Supported notifier to send metrics is Datadog (using DogStatsD).

# Requirements
## DogStatsD

The cyqldog uses [DogStatsD](https://docs.datadoghq.com/guides/dogstatsd/) to send metrics to Datadog.

DogStatsD is a metrics aggregation service bundled with the Datadog Agent (dd-agent).
If you want to send metrics to Datadog, you may alreadly have installed dd-agent including DogStatsD,
but if you have not installed dd-agent,

* install dd-agent (https://docs.datadoghq.com/guides/basic_agent_usage/)
* or you can run dd-agent with Docker. (https://github.com/DataDog/docker-dd-agent)
* or you can run only DogStatsD with Docker. (https://github.com/DataDog/docker-dd-agent#dogstatsd)

When using Docker, please make sure the UDP 8125 port is open by adding the option `-p 8125:8125/udp` to the `docker run` command.

# Install
## cyqldog

Download the latest compiled binaries and put it anywhere in your executable path.

https://github.com/crowdworks/cyqldog/releases

# Usage
Create your configuration file, and pass it to `-C` flag to run cysqldog.

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
  # Currently supported options are as follows.
  # All options of go-sql-driver/mysql should be supported, but it is not implemented yet.
  #
  # driver: mysql
  # options:
  #   host: db.example.com
  #   port: 3306
  #   user: cyqldog
  #   password: {{ .DB_PASSWORD }}
  #   dbname: cyqldogdb

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
    # * playground.cyqldog.test2.value1 (with tags ["env:local", "source:db.example.com", (value of column tag1), (value of column tag2)])
    # * playground.cyqldog.test2.value2 (with tags ["env:local", "source:db.example.com", (value of column tag1), (value of column tag2)])
    tag_cols:
      - tag1
      - tag2
    value_cols:
      - val1
      - val2
```

# Contributions

Any feedback and contributions are welcome. Feel free to report an issue and submit a pull request.

# License

The cyqldog is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
