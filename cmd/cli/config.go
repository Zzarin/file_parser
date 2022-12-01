package main

import "time"

type Config struct {
	DbDsn             string        `long:"db-dsn" description:"Dsn to connect to the postgreSQL" env:"DB_DSN"`
	DbConnMaxLifetime time.Duration `long:"db-conn-max-lifetime" description:"ConnMaxLifetime " env:"DB_CONN_MAX_LIFETIME"`
	DbMaxOpenConns    int           `long:"db-max-open-conns" description:"Db proper configuration of max open connections" env:"DB_MAX_OPEN_CONNS"`
	DbMaxIdleConns    int           `long:"db-max-idle-conns" description:"Db proper configuration of max idle connections " env:"DB_MAX_IDLE_CONNS"`
	FileName          string        `short:"f" long:"os-file-name" description:"Name of the file to read" env:"FILE_NAME"`
	StrReadingFreq    time.Duration `short:"t" long:"str-reading-freq" description:"Time in sec between readings of each string" env:"READ_FREQ"`
}
