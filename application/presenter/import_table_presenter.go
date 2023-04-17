package presenter

import (
	"go-import-from-s3/application/usecase"
	"go-import-from-s3/domain"
)

type importTablePresenter struct{}

func NewImportTablePresenter() usecase.ImportTablePresenter {
	return importTablePresenter{}
}

func (a importTablePresenter) Output(table domain.ImportTable) usecase.ImportTableOutput {
	return usecase.ImportTableOutput{
		ID:     table.ID,
		Status: table.Status,
	}
}
