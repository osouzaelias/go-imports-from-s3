package main

import (
	"go-import-from-s3/infrastructure"
	"go-import-from-s3/infrastructure/database"
)

func main() {
	var app = infrastructure.NewConfig().DbNoSQL(database.InstanceDynamoDB)
	app.Start()
}
