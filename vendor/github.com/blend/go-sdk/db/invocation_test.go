package db

import (
	"database/sql"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestInvocationLabels(t *testing.T) {
	assert := assert.New(t)

	inv := &Invocation{}
	inv = inv.WithLabel("test")
	assert.NotEmpty(inv.Label())
}

type jsonTestChild struct {
	Label string `json:"label"`
}

type jsonTest struct {
	ID   int    `db:"id,pk,auto"`
	Name string `db:"name"`

	NotNull  jsonTestChild `db:"not_null,json"`
	Nullable []string      `db:"nullable,json"`
}

func (jt jsonTest) TableName() string {
	return "json_test"
}

func createJSONTestTable(tx *sql.Tx) error {
	return Default().InTx(tx).Exec("create table json_test (id serial primary key, name varchar(255), not_null json, nullable json)")
}

func dropJSONTextTable(tx *sql.Tx) error {
	return Default().InTx(tx).Exec("drop table if exists json_test")
}

func TestInvocationJSONNulls(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()
	defer dropJSONTextTable(tx)

	assert.Nil(createJSONTestTable(tx))

	// try creating fully set object and reading it out
	obj0 := jsonTest{Name: uuid.V4().String(), NotNull: jsonTestChild{Label: uuid.V4().String()}, Nullable: []string{uuid.V4().String()}}
	assert.Nil(Default().InTx(tx).Create(&obj0))

	var verify0 jsonTest
	assert.Nil(Default().InTx(tx).Get(&verify0, obj0.ID))

	assert.Equal(obj0.ID, verify0.ID)
	assert.Equal(obj0.Name, verify0.Name)
	assert.Equal(obj0.Nullable, verify0.Nullable)
	assert.Equal(obj0.NotNull.Label, verify0.NotNull.Label)

	// try creating partially set object and reading it out
	obj1 := jsonTest{Name: uuid.V4().String(), NotNull: jsonTestChild{Label: uuid.V4().String()}} //note `Nullable` isn't set

	columns := Columns(obj1)
	values := columns.ColumnValues(obj1)
	assert.Len(values, 4)
	assert.Nil(values[3], "we shouldn't emit a literal 'null' here")
	assert.NotEqual("null", values[3], "we shouldn't emit a literal 'null' here")

	assert.Nil(Default().InTx(tx).Create(&obj1))

	var verify1 jsonTest
	assert.Nil(Default().InTx(tx).Get(&verify1, obj1.ID))

	assert.Equal(obj1.ID, verify1.ID)
	assert.Equal(obj1.Name, verify1.Name)
	assert.Nil(verify1.Nullable)
	assert.Equal(obj1.NotNull.Label, verify1.NotNull.Label)

	any, err := Default().InTx(tx).Query("select 1 from json_test where id = $1 and nullable is null", obj1.ID).Any()
	assert.Nil(err)
	assert.True(any, "we should have written a sql null, not a literal string 'null'")

	// set it to literal 'null' to test this is backward compatible
	assert.Nil(Default().InTx(tx).Exec("update json_test set nullable = 'null' where id = $1", obj1.ID))

	var verify2 jsonTest
	assert.Nil(Default().InTx(tx).Get(&verify2, obj1.ID))
	assert.Equal(obj1.ID, verify2.ID)
	assert.Equal(obj1.Name, verify2.Name)
	assert.Nil(verify2.Nullable, "even if we set it to literal 'null' it should come out golang nil")
	assert.Equal(obj1.NotNull.Label, verify2.NotNull.Label)
}
