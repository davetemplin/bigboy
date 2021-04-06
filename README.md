# Bigboy

[![Actions Status](https://github.com/igor-starostenko/bigboy/workflows/build/badge.svg)](https://github.com/igor-starostenko/bigboy/actions)
[![codecov](https://codecov.io/gh/igor-starostenko/bigboy/branch/master/graph/badge.svg)](https://codecov.io/gh/igor-starostenko/bigboy)

High data-rate SQL-to-JSON extraction from SQL Server, PostgreSQL, or MySQL.

Written by Dave Templin
Maintained by Igor Starostenko

## Table of Contents
- [Overview](#overview)
  - [Features](#features)
  - [Quickstart](#quickstart)
- [Concepts](#concepts)
  - [Connections](#connections)
  - [Targets](#targets)
  - [Prefetching](#prefetching)
  - [Parameterization](#parameterization)
  - [Transforms](#transforms)
    - [Nesting](#nesting)
    - [Scripting](#scripting)
    - [Split Output](#split-output)
    - [Timezones](#timezones)
- [Reference](#reference)
  - [Configuration](#configuration)
    - [connections](#connection)
  - [Target](#target)
    - [nest](#nest)
    - [params](#params)
    - [split](#split)
- [Development](#development)
  - [Build](#build)
  - [Test](#test)
  - [Cross compile](#cross-compile)
  - [Release](#release)
- [Helpful links](#helpful-links)

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
* The data is extracted in a JSONL (newline delimited JSON) format in `out/${target}/` directory (by default) with current date filename

[See basic example](/EXAMPLES.md#basic)

# Concepts

## Connections

Authentication credentials and connection configuration can be set for all database sources in a single config file (See [configuration schema](#configuration))

## Targets

A directory path that contains required files (instructions) on how to extract data from a source database. Target may contain the following files:

- `target.json` - **REQUIRED** - core configuration file of a target that specifiec which connection to use as well as other parameters (See [target.json schema](#target))
- `extract.sql` - **REQUIRED** - SQL query to be used to extract data from source database for a given connection. If prefetch is set up, needs to apply string interpolation (For example `WHERE id IN (%s)`)
- `prefetch.sql` - *OPTIONAL* - if `target.json` has `prefetch` set to `true`, the `prefetch.sql` loads ids to be passed into `extract.sql` and run in parallel
- `nest.sql` - *OPTIONAL* - if `target.json` has `nest` array of objects configured, the `nest.sql` loads ids to be added into the main output with new column for each object

## Prefetching

To extract data on large datasets it's recommended to utilize the prefetch feature that allows to split extract in multiple parallel processes. When [target](#target) has `prefetch` set to `true`, the `prefetch.sql` has to be provided that has an SQL query with a `SELECT` statement with one integer column.
The results of the prefetch query would be be passed into multiple processes each execute an `extract.sql` with a subset of the prefetched numbers. The extract query needs to apply string interpolation (For example `WHERE id IN (%s)`)

[See prefetch example](/EXAMPLES.md#prefetch)

## Parameterization

Queries can accept parameters to customize the request from outside the target on runtime. Params can be configured in [target](#target) with a given type and default value. On runtime the params are passed into a `prefetch.sql` or an `extract.sql` using syntax like `WHERE date >= ?1`, where `1` is the index of the param (starting from `1`).

[See params example](/EXAMPLES.md#params)

## Transforms

### Nesting

Allows for the output JSON records to have nested arrays. To enable, the [target](#target) has to have `nest` params configured and `nest.sql` defined.
The nest query runs before the main extract and then inserts the results in an array format into a column defined in target `childKey` param. For records to be inserted into the result of the main extract query, the `nest.sql` needs to have a column with the name `_parent` which value would match the value of the column in the `extract.sql` which name is defined in target `parentKey` param.
If `nest.sql` has a `_value` column only the values would be inserted in the array. Otherwise the entire record object (except the `_parent` column) would be added inside the array.

[See nest example](/EXAMPLES.md#nest)

### Scripting

*Not yet implemented*

### Split output

Produces multiple files instead of one. To enable, the [target](#target) has to have `split` param set.
When configured, the output argument set with `-o` flag has to be a directory without extension.
Currently supports only split by `date`. Every output file would contain records for each day returned from the extract query. Use `layout` param for MySQL or when the date column is stored as STRING.

[See split example](/EXAMPLES.md#split)

### Timezone format

All dates are assumed to be in GMT unless a timezone is specified.
If a time is not specified then midnight GMT is assumed.
Examples below illustrate various scenarios of specifying a date or date-range.

The following examples assume there is a target named `log` with a single paramter of type `date` representing a start date for the extraction.

| Example                               | Comments |
| ------------------------------------- | ------------------------------------------------------- |
| `bigboy log 2017-07-21`               | 7/21/2017 at midnight GMT |
| `bigboy log "2017-07-21 15:00:00"`    | 7/21/2017 at 3pm GMT |
| `bigboy log today`                    | Midnight GMT of the current day |
| `bigboy log yesterday`                | Midnight GMT of the previous day |

The following examples assume there is a target named `sales` with two parameters of type `date` representing a date range for the extraction.

| Example                               | Comments |
| ------------------------------------- | ------------------------------------------------------- |
| `bigboy sales 2017-07-21 2017-07-23`  | From 7/21/2017 to 7/23/2017 midnight-to-midnight GMT |
| `bigboy sales 2017-07-21 2d`          | Midnight GMT of the previous day |

> The time-zone database needed by LoadLocation may not be present on all systems, especially non-Unix systems. LoadLocation looks in the directory for an uncompressed zip file named by the ZONEINFO environment variable, if any, then looks in known installation locations on Unix systems, and finally looks in $GOROOT/lib/time/zoneinfo.zip.

# Reference

## Command Line Arguments

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

This section describes the `bigboy.json` (default, unless `-c` flag is used) file format.

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `connections` | [connection](#connection)[] | + | Array of connection configurations for each database source |
| `errors` | INTEGER | - | Total number of errors to ignore before aborting |
| `nulls` | BOOLEAN | - | If nulls are allowed to be included in the output |
| `page` | INTEGER | - | Number of rows per page |
| `quiet` | BOOLEAN | - | If the terminal output should be supressed |
| `retries` | INTEGER | - | Number of retries to perform in case of an error |
| `workers` | INTEGER | - | Number of workers to run extracts in parallel |

### connection

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `driver` | STRING | + | Database driver (`mysql`, `mssql`, `postgres`) |
| `server` | STRING | + | Database connection string |
| `database` | STRING | + | Database name |
| `dsn` | STRING | - | DB dsn |
| `port` | INTEGER | - | DB port |
| `user` | STRING | - | DB username for authentication |
| `password` | STRING | - | DB password for authentication |
| `max` | INTEGER | - | Maximum number of open database connections |
| `timezone` | STRING | - | Can be `UTC` or `Local` or IANA Time Zone database format (For example `America/Los_Angeles`) |

## target

This section describes the `target.json` file format.

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `connection` | STRING | + | Connection name |
| `fetch` | STRING | - | File name for the main extract SQL query. Default `extract.sql` |
| `params` | [param](#param) | - | Allows to pass params into a query. For example to filter the data |
| `prefetch` | BOOLEAN | - | If `prefetch.sql` should be used for parallel extraction |
| `nest` | [nest](#nest)[] | - | Array of columns to be added for each record |
| `script` | STRING | - | *Not yet implemented* |
| `split` | [split](#split) | - | Produces multiple files instead of one |
| `timezone` | STRING | - | Defaults to connection timezone |

### nest

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `connection` | STRING | - | Allows to use different database source. Uses target connection by default |
| `childKey` | STRING | + | New array property that would contain all matched records |
| `parentKey` | STRING | + | Has to be an integer |
| `fetch` | STRING | - | File name for the nest SQL query. Default `nest.sql` for the first nest column |
| `timezone` | STRING | - | Defaults to connection timezone |

### param

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | STRING | - | Name of the param. Only use for logging and readability |
| `type` | STRING | - | Converts the text argument to the given type. One of `integer`, `float`, `string`, `boolean` or `date`. **Required** unless param is `null` |
| `default` | STRING | - | Setting default param if no command line params are passed |

### split

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `by` | STRING | + | Has to be set to `date` |
| `layout` | STRING | - | Converts to date from string for MySQL if the date column is STRING (See [golang layout format](https://yourbasic.org/golang/format-parse-string-time-date-example/) |
| `value` | STRING | + | Column name which contains the value by which the files to split |

# Development

## Build

Install [golang](https://golang.org/dl/).

Clone the repo and `cd` into it.

```bash
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

# Helpful links

There are lots of ways to approach ETL, and lots of vendors that want to sell you a solution!
Here are some additional references that may be helpful...

* [Wikipedia article on ETL](https://en.wikipedia.org/wiki/Extract,_transform,_load)
* [Performing ETL from a Relational Database into BigQuery](https://cloud.google.com/solutions/performing-etl-from-relational-database-into-bigquery)
* [ETL Software: Top 63](https://www.predictiveanalyticstoday.com/top-free-extract-transform-load-etl-software/)
