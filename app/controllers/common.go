package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/Syncano/orion/app/api"
	"github.com/Syncano/orion/app/validators"
	"github.com/Syncano/orion/pkg/cache"
	"github.com/Syncano/orion/pkg/settings"
)

// CacheInvalidate ...
func CacheInvalidate(c echo.Context) error {
	v := &validators.CacheInvalidateForm{}

	if err := api.BindValidateAndExec(c, v, func() error {
		key := fmt.Sprintf("%s:%s", v.VersionKey, settings.Common.SecretKey)
		hash := hmac.New(sha256.New, []byte(key)).Sum(nil)
		if v.Signature != hex.EncodeToString(hash) {
			return api.NewGenericError(http.StatusBadRequest, "Invalid signature.")
		}

		cache.InvalidateVersion(v.VersionKey, settings.Common.CacheTimeout)
		return nil
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
