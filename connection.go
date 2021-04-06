package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Connection ...
type Connection struct {
	Driver   string `json:"driver"`
	Server   string `json:"server"`
	Database string `json:"database"`
	Dsn      string `json:"dsn"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Timezone string `json:"timezone"`
	Max      int    `json:"max"`
	db       *sql.DB
	location *time.Location
}

func connect(key string) *Connection {
	connection, ok := config.Connections[key]

	if !ok {
		stop(fmt.Sprintf("Invalid connection key: '%s'", key), 1)
	}

	dsn := formatDsn(connection)
	if dsn == "" {
		stop(fmt.Sprintf("Invalid driver specified for connection: '%s'", key), 1)
	}

	if connection.db == nil {
		var err error
		connection.db, err = sql.Open(connection.Driver, dsn)
		check(err)

		err = connection.db.Ping()
		if err != nil {
			stop(fmt.Sprintf("Unable to establish connection to server \"%s\"", connection.Server), 3)
		}

		if connection.Max != 0 {
			connection.db.SetMaxOpenConns(connection.Max)
		}

		if connection.Timezone != "" {
			connection.location, err = time.LoadLocation(connection.Timezone)
			check(err)
		}
	}

	return connection
}

func disconnect() {
	for _, connection := range config.Connections {
		if connection.db != nil {
			err := connection.db.Close()
			check(err)
			connection.db = nil
		}
	}
}

func formatDsn(connection *Connection) string {
	if connection.Dsn != "" && !strings.HasPrefix(connection.Dsn, "...") {
		return connection.Dsn
	}

	var dsn string
	if connection.Driver == "mssql" {
		dsn = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d",
			connection.Server,
			connection.User,
			connection.Password,
			connection.Port)
	} else if connection.Driver == "postgres" {
		dsn = fmt.Sprintf("host=%s user=%s password='%s' port=%d dbname=%s",
			strings.Replace(connection.Server, " ", "\\ ", -1),
			strings.Replace(connection.User, " ", "\\ ", -1),
			strings.Replace(connection.Password, "'", "\\'", -1),
			connection.Port,
			strings.Replace(connection.Database, " ", "\\ ", -1))
	} else if connection.Driver == "mysql" {
		dsn = fmt.Sprintf("%s:%s@%s/%s",
			connection.User,
			connection.Password,
			connection.Server,
			connection.Database)
	}

	if strings.HasPrefix(connection.Dsn, "...") {
		dsn += connection.Dsn[3:]
	}

	return dsn
}
