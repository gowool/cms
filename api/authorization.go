package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/gowool/cms"
)

var (
	ErrCallTargetNotFound   = errors.New("huma api: call target not found")
	ErrOperationMetadataNil = errors.New("huma api: operation metadata is nil")
)

type AuthorizationConfig struct {
	Skipper    func(huma.Context) bool
	Authorizer cms.Authorizer
	Logger     *zap.Logger
}

func Authorization(cfg AuthorizationConfig) func(huma.API) func(huma.Context, func(huma.Context)) {
	if cfg.Authorizer == nil {
		panic("authorizer is not specified")
	}
	if cfg.Logger == nil {
		panic("logger is not specified")
	}

	if cfg.Skipper == nil {
		cfg.Skipper = func(huma.Context) bool {
			return false
		}
	}

	return func(api huma.API) func(huma.Context, func(huma.Context)) {
		unauthorized := func(ctx huma.Context, errs ...error) {
			status := http.StatusUnauthorized
			message := http.StatusText(status)

			if err := huma.WriteErr(api, ctx, status, message, errs...); err != nil {
				cfg.Logger.Error("huma api: failed to write error", zap.Error(err))
			}
		}

		fn := func(c huma.Context) error {
			if cfg.Skipper(c) {
				return nil
			}

			o := c.Operation()
			if o.Metadata == nil {
				return ErrOperationMetadataNil
			}

			target, ok := o.Metadata["target"].(*cms.CallTarget)
			if !ok {
				return ErrCallTargetNotFound
			}

			claims := cms.CtxClaims(c.Context())

			decision, err := cfg.Authorizer.Authorize(c.Context(), claims, target)
			if err != nil {
				return err
			}

			if decision != cms.DecisionAllow {
				if claims.Scheme == cms.BasicScheme {
					c.SetHeader(echo.HeaderWWWAuthenticate, "basic realm=Restricted")
				}
				return fmt.Errorf("huma api: authorizer decision `%s`", decision)
			}
			return nil
		}

		return func(c huma.Context, next func(huma.Context)) {
			if err := fn(c); err != nil {
				unauthorized(c)
				cfg.Logger.Error("huma api: failed to authorize", zap.Error(err))
				return
			}
			next(c)
		}
	}
}
