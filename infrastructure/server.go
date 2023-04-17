package infrastructure

import (
	"go-import-from-s3/application/adapter/repository"
	"go-import-from-s3/application/presenter"
	"go-import-from-s3/application/usecase"
	"go-import-from-s3/infrastructure/database"
	"log"
)

type config struct {
	dbNoSQL repository.NoSQL
}

func NewConfig() *config {
	return &config{}
}

func (c *config) DbNoSQL(instance int) *config {
	db, err := database.NewDatabaseNoSQLFactory(instance)
	if err != nil {
		log.Fatalln(err, "Could not make a connection to the database")
	}

	log.Println("Successfully connected to the NoSQL database")

	c.dbNoSQL = db
	return c
}

func (c *config) Start() {
	uc := usecase.NewImportTableInteractor(
		repository.NewImportTableNoSQL(c.dbNoSQL),
		presenter.NewImportTablePresenter(),
	)

	uc.Execute()
}
