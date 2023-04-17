package domain

import (
	"errors"
)

var ErrResourceInUse = errors.New("a operação entra em conflito com a disponibilidade do recurso. " +
	"Por exemplo, você tentou recriar uma tabela existente, " +
	"ou tentou excluir uma tabela atualmente no estado CREATING")

var ErrLimitExceeded = errors.New("Não há limite para o número de backups diários sob " +
	"demanda que podem ser feitos.")

var ErrImportConflict = errors.New("houve um conflito ao importar da fonte S3 especificada. " +
	"Isso pode ocorrer quando a importação atual entra em conflito com uma solicitação de " +
	"importação anterior que tinha o mesmo token de cliente.")

var ErrImportNotFound = errors.New("a importação especificada não foi encontrada")

type ImportTableRepository interface {
	Import(ImportTable) (ImportTable, error)
	Describe(ImportTable) (ImportTable, error)
}

type ImportTable struct {
	ID       string
	FileName string
	Status   string
}

//func NewImportTable(ID string, fileName string, Status string) ImportTable {
//	return ImportTable{
//		ID:       ID,
//		fileName: fileName,
//		Status:   Status,
//	}
//}
//
//func (a ImportTable) ID() string {
//	return a.ID
//}
//
//func (a ImportTable) FileName() string {
//	return a.fileName
//}
//
//func (a ImportTable) Status() string {
//	return a.Status
//}
