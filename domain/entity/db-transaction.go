package entity

type DBTransaction interface {
	Commit() error
	Rollback() error
}
