package db

import "database/sql"

// DatabaseMapped is the interface that any objects passed into database mapped methods like Create, Update, Delete, Get, GetAll etc.
type DatabaseMapped interface{}

// TableNameProvider is a type that implements the TableName() function.
// The only required method is TableName() string that returns the name of the table in the database this type is mapped to.
//
//	type MyDatabaseMappedObject {
//		Mycolumn `db:"my_column"`
//	}
//	func (_ MyDatabaseMappedObject) TableName() string {
//		return "my_database_mapped_object"
//	}
// If you require different table names based on alias, create another type.
type TableNameProvider interface {
	TableName() string
}

// Populatable is an interface that you can implement if your object is read often and is performance critical.
type Populatable interface {
	Populate(rows *sql.Rows) error
}

// RowsConsumer is the function signature that is called from within Each().
type RowsConsumer func(r *sql.Rows) error
