package api

import (
	"context"
	"net/http"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"github.com/Syncano/orion/app/settings"
	"github.com/Syncano/pkg-go/v2/database/manager"
	"github.com/Syncano/pkg-go/v2/util"
)

const validationRetries = 2

var unsafeMethods = map[string]empty{
	"POST":   {},
	"PUT":    {},
	"PATCH":  {},
	"DELETE": {},
}

// IsSafeMethod returns true if http request method is unsafe.
func IsSafeMethod(meth string) bool {
	_, ok := unsafeMethods[meth]
	return !ok
}

// BindAndValidate binds and validates object against echo framework.
func BindAndValidate(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}

// BindValidateAndExec binds, validates object and executes function with retry.
func BindValidateAndExec(c echo.Context, i interface{}, fn func() error) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	_, err := util.RetryWithCritical(validationRetries, 0, func() (bool, error) {
		if err := c.Validate(i); err != nil {
			return true, err
		}

		return false, fn()
	})

	return err
}

// SimpleDelete selects for update and deletes object, returning 201 if everything went fine.
func SimpleDelete(c echo.Context, mgr Deleter, q *orm.Query, v Verboser) error {
	if err := mgr.RunInTransaction(func(tx *pg.Tx) error {
		if err := manager.Lock(q); err != nil {
			if err == pg.ErrNoRows {
				return NewNotFoundError(v)
			}

			return err
		}
		return mgr.Delete(v)
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func AddRequestID(ctx context.Context, c echo.Context) (outCtx context.Context, reqID string) {
	return util.AddRequestID(ctx, func() string { return c.Get(settings.ContextRequestID).(string) })
}
