package repository

type NoSQL interface {
	ImportTable(string, interface{}) error
	DescribeImport(string, interface{}) error
}
