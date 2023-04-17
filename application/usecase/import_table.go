package usecase

import (
	"go-import-from-s3/domain"
	"log"
)

type ImportTableUseCase interface {
	Execute()
}

type ImportTableInput struct {
	FileName string
}

type ImportTablePresenter interface {
	Output(table domain.ImportTable) ImportTableOutput
}

type ImportTableOutput struct {
	ID     string
	Status string
}

type importTableInteractor struct {
	repo      domain.ImportTableRepository
	presenter ImportTablePresenter
}

func NewImportTableInteractor(repo domain.ImportTableRepository, presenter ImportTablePresenter) ImportTableUseCase {
	return importTableInteractor{
		repo:      repo,
		presenter: presenter,
	}
}

func (a importTableInteractor) Execute() {
	// todo: buscar no bucket o nome do arquivo
	//var table = domain.NewImportTable("", "filename.csv", "")

	table := domain.ImportTable{FileName: "filename.csv"}
	table, err := a.repo.Import(table)
	if err != nil {
		log.Fatalln(err)
	}

	// todo: aguardar a conclusão da importação
}
