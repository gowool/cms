package cms

import (
	"reflect"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cast"

	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

func BasicAuthValidator(repo repository.Admin) middleware.BasicAuthValidator {
	return func(user, password string, c echo.Context) (bool, error) {
		r := c.Request()
		ctx := r.Context()

		admin, err := repo.FindByEmail(ctx, user)
		if err != nil {
			return false, err
		}
		if err = admin.Password.Validate(password); err != nil {
			return false, err
		}

		ctx = WithAdmin(ctx, &admin)
		ctx = WithClaims(ctx, &Claims{
			Subject: &admin,
			Scheme:  BasicScheme,
		})

		c.SetRequest(r.WithContext(ctx))

		return true, nil
	}
}

func JWTAuthValidator(repo repository.Admin, secret string) func(string, echo.Context) (bool, error) {
	return func(token string, c echo.Context) (bool, error) {
		r := c.Request()
		ctx := r.Context()

		claims, _ := ParseUnverifiedJWT(token)
		switch cast.ToString(claimsValue(claims, "model")) {
		case reflect.TypeOf((*model.Admin)(nil)).Elem().Name():
			subject, err := claims.GetSubject()
			if err != nil {
				return false, err
			}

			admin, err := repo.FindByEmail(ctx, subject)
			if err != nil {
				return false, err
			}

			if _, err = ParseJWT(token, admin.Salt+secret); err != nil {
				return false, err
			}

			ctx = WithAdmin(ctx, &admin)
			ctx = WithClaims(ctx, &Claims{
				Subject: &admin,
				Scheme:  JWTScheme,
				TwoFA:   cast.ToBool(claimsValue(claims, "2fa")),
			})

			c.SetRequest(r.WithContext(ctx))

			return true, nil
		}
		return false, nil
	}
}

func claimsValue(claims jwt.MapClaims, key string) any {
	if v, ok := claims[key]; ok {
		return v
	}
	return nil
}
