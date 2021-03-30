# Bigboy

[![Actions Status](https://github.com/igor-starostenko/bigboy/workflows/build/badge.svg)](https://github.com/igor-starostenko/bigboy/actions)
[![codecov](https://codecov.io/gh/igor-starostenko/bigboy/branch/master/graph/badge.svg)](https://codecov.io/gh/igor-starostenko/bigboy)

High data-rate SQL-to-JSON extraction from SQL Server, PostgreSQL, or MySQL.

Written by Dave Templin
Maintained by Igor Starostenko

# Overview

Bigboy is basically a **SQL-TO-JSON** tool that extracts data from SQL Server, PostgreSQL, or MySQL databases.
It is designed to handle extremely high data extraction rates (multi-million rows / gigabyte-range) which is achieved by running multiple concurrent extraction queries from a configurable thread pool.
The tool provides a simple configuration model for managing any number of data extractions.
It also exposes a simple and minimal command-line interface (CLI) that works great for adhoc or batch/cron use-cases.

## Features

* Extract data from SQL Server, PostgreSQL, or MySQL
* Perform simple SQL-to-JSON transformations
* Nest rows to form complex hierarchical (or document-oriented) data
* Define configuration parameters to manage dynamic queries
* Run queries in parallel from a configurable thread pool for high data rates
* Combine data from multiple different database sources
* Apply timezone to dates stored without a timezone

## Quickstart

* Download the [latest release](https://github.com/igor-starostenko/bigboy/releases/latest) of bigboy for your operating system
* Create a `bigboy.json` with the list of `connections` (See [configuration details](#configuration) for more information)
* Create a target directory with `target.json` file with `connection` name (See [target](#targets) for more information)
* Create `extract.sql` file in the newly created target directory with the SQL query to use to extract the data (See [prefetch](#prefetching) section to extract data in parallel)
* Consider using available [transforms](#transforms) to modify the data on the fly
* Extract the data from the target: `./bigboy ${target}` (You may need to change access permissions to the executable by running `chmod +x bigboy`)

# Concepts

## Connections

## Targets

## Prefetching

fetch, prefetch

## Transforms

nest, script, split, timezone
special field names: _parent, _value

# Reference

## Command Arguments

* `-c` Bigboy config file path *(default=\"bigboy.json\")*
* `-e` Maximum overall number of errors before aborting *(default=100)*
* `-n` Include nulls in output *(default=false)*
* `-o` Output directory *(creates \"out\" directory if not specified)*
* `-p` Number of rows per page extracted *(default=1000)*
* `-q` Supress informational output *(default=false)*
* `-r` Number of consecutive errors before aborting *(default=3)*
* `-v` Print version info about bigboy and exit
* `-w` Number of background workers *(default=4)*

> Above flags take priority over the configuration in `bigboy.json` file.

## Configuration

This section describes the `bigboy.json` file format.

| Name | Description |
| --- | --- |
| `connections` | ... |
| `errors` | ... |
| `nulls` | ... |
| `page` | ... |
| `quiet` | ... |
| `retries` | ... |
| `workers` | ... |

### connections

| Name | Description |
| --- | --- |
| `driver` | ... |
| `server` | ... |
| `database` | ... |
| `dsn` | ... |
| `port` | ... |
| `user` | ... |
| `password` | ... |
| `max` | ... |
| `timezone` | ... |

## target.json

This section describes the `target.json` file format.

| Name | Description |
| --- | --- |
| `connection` | ... |
| `fetch` | ... |
| `params` | ... |
| `prefetch` | ... |
| `nest` | ... |
| `script` | ... |
| `split` | ... |
| `timezone` | ... |

### nest

| Name | Description |
| --- | --- |
| `connection` | ... |
| `childKey` | ... |
| `parentKey` | ... |
| `fetch` | ... |
| `timezone` | ... |

### param

| Name | Description |
| --- | --- |
| `name` | ... |
| `type` | ... |
| `default` | ... |

### split

| Name | Description |
| --- | --- |
| `by` | ... |
| `value` | ... |

## Date Format

All dates are assumed to be in GMT unless a timezone is specified.
If a time is not specified then midnight GMT is assumed.
Examples below illustrate various scenarios of specifying a date or date-range.

The following examples assume there is a target named `log` with a single paramter of type `date` representing a start date for the extraction.

| Example                               | Comments
| ------------------------------------- | ------------------------------------------------------- |
| `bigboy log 2017-07-21`               | 7/21/2017 at midnight GMT
| `bigboy log "2017-07-21 15:00:00"`    | 7/21/2017 at 3pm GMT
| `bigboy log today`                    | Midnight GMT of the current day
| `bigboy log yesterday`                | Midnight GMT of the previous day

The following examples assume there is a target named `sales` with two paramters of type `date` representing a date range for the extraction.

| Example                               | Comments
| ------------------------------------- | ------------------------------------------------------- |
| `bigboy sales 2017-07-21 2017-07-23`  | From 7/21/2017 to 7/23/2017 midnight-to-midnight GMT
| `bigboy sales 2017-07-21 2d`          | Midnight GMT of the previous day.

> The time-zone database needed by LoadLocation may not be present on all systems, especially non-Unix systems. LoadLocation looks in the directory for an uncompressed zip file named by the ZONEINFO environment variable, if any, then looks in known installation locations on Unix systems, and finally looks in $GOROOT/lib/time/zoneinfo.zip.

# Build

Install [golang](https://golang.org/dl/)

```bash
$ git clone https://github.com/davetemplin/bigboy.git
$ go build
```

## Test

```bash
$ go test .
```

## Cross compile

On Windows:
```bash
$ build windows
$ build linux
$ build darwin
```

On unix:
```bash
bash build.sh
```

## Release

Upload build executables is automated using `release` workflow in GitHub actions. It gets triggered by a push of a new git tag to origin.
Make sure to keep [CHANGELOG.md](./CHANGELOG.md) updated and [version](./constants.go) incremented.

# References

There are lots of ways to approach ETL, and lots of vendors that want to sell you a solution!
Here are some additional references that may be helpful...

* [Wikipedia article on ETL](https://en.wikipedia.org/wiki/Extract,_transform,_load)
* [Performing ETL from a Relational Database into BigQuery](https://cloud.google.com/solutions/performing-etl-from-relational-database-into-bigquery)
* [ETL Software: Top 63](https://www.predictiveanalyticstoday.com/top-free-extract-transform-load-etl-software/)
