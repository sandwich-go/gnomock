package main

// all known presets should go right here so that they are available when
// requested over HTTP:
import (
	_ "github.com/sandwich-go/gnomock/preset/cockroachdb"
	_ "github.com/sandwich-go/gnomock/preset/elastic"
	_ "github.com/sandwich-go/gnomock/preset/influxdb"
	_ "github.com/sandwich-go/gnomock/preset/k3s"
	_ "github.com/sandwich-go/gnomock/preset/kafka"
	_ "github.com/sandwich-go/gnomock/preset/localstack"
	_ "github.com/sandwich-go/gnomock/preset/mariadb"
	_ "github.com/sandwich-go/gnomock/preset/memcached"
	_ "github.com/sandwich-go/gnomock/preset/mongo"
	_ "github.com/sandwich-go/gnomock/preset/mssql"
	_ "github.com/sandwich-go/gnomock/preset/mysql"
	_ "github.com/sandwich-go/gnomock/preset/postgres"
	_ "github.com/sandwich-go/gnomock/preset/rabbitmq"
	_ "github.com/sandwich-go/gnomock/preset/redis"
	_ "github.com/sandwich-go/gnomock/preset/splunk"
	// new presets go here
)
