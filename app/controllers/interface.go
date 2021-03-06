package controllers

import (
	"reflect"

	"github.com/labstack/echo/v4"

	"github.com/Syncano/orion/app/api"
	"github.com/Syncano/orion/app/serializers"
)

//go:generate go run github.com/vektra/mockery/cmd/mockery -inpkg -testonly -name Paginator
type Paginator interface {
	FilterObjects(cursor Cursorer) error
	ProcessObjects(c echo.Context, cursor Cursorer, typ reflect.Type, serializer serializers.Serializer, responseLimit *int) ([]api.RawMessage, error)
	CreateCursor(c echo.Context, defaultOrderAsc bool) Cursorer
}

// Assert interface compatibility.
var (
	_ Paginator = (*PaginatorDB)(nil)
	_ Paginator = (*PaginatorOrderedDB)(nil)
	_ Paginator = (*PaginatorRedis)(nil)
)

//go:generate go run github.com/vektra/mockery/cmd/mockery -inpkg -testonly -name Cursorer
type Cursorer interface {
	NextURL(path string) string
	PrevURL(path string) string

	Limit() int
	LastPK() int
	IsOrderAsc() bool
	IsForward() bool
	SetFirst(interface{})
	SetLast(interface{})
}

// Assert interface compatibility.
var (
	_ Cursorer = (*cursor)(nil)
	_ Cursorer = (*keysetcursor)(nil)
)
