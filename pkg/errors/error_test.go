package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrUndefined(t *testing.T) {
	err := Wrap(errors.New("undefined error"))
	a := assert.New(t)
	a.Equal(err.code, ErrUndefined)
}

func TestErrBusiness(t *testing.T) {
	err := New(ErrBusinessLogicFailure)
	a := assert.New(t)
	a.Equal(err.code, ErrBusinessLogicFailure)
}
