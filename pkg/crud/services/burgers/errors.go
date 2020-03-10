package burgers

import "fmt"

type DbError struct {
	err error
}

type QueryError struct {
	query string
	err error
}

func NewQueryError(query string, err error) *QueryError {
	return &QueryError{query: query, err: err}
}

func (e *QueryError) Unwrap() error {
	return e.err
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("can't execute query %s: %s", e.query, e.err.Error())
}

func (e *DbError) Error() string {
	return fmt.Sprintf("can't execute db operation: %s", e.err.Error())
}

func (e *DbError) Unwrap() error {
	return e.err
}

func NewDbError(err error) *DbError {
	return &DbError{err:err}
}

