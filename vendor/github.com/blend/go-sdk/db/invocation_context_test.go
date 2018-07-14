package db

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCtxInTxUsesArguments(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	withTx := NewInvocationContext(Default()).InTx(tx)
	assert.NotNil(withTx.tx)
}

func TestCtxInTx(t *testing.T) {
	assert := assert.New(t)

	withTx := NewInvocationContext(Default()).InTx()
	assert.Nil(withTx.tx)
}

func TestCtxInvoke(t *testing.T) {
	assert := assert.New(t)

	inv := NewInvocationContext(Default()).Invoke()
	assert.Nil(inv.Validate())
}

func TestCtxInvokeError(t *testing.T) {
	assert := assert.New(t)

	inv := NewInvocationContext(nil).Invoke()
	assert.NotNil(inv.Validate(), "should fail the connection not nil check")
}
