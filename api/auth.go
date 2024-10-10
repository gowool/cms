package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/gowool/cms"
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

const prefixRefreshToken = "token:refresh:"

type SignIn struct {
	Body struct {
		Email    string `json:"email,omitempty" required:"true" minLength:"3" maxLength:"254" format:"email"`
		Password string `json:"password,omitempty" required:"true" minLength:"8" maxLength:"64"`
	}
}

type OTP struct {
	Body struct {
		Password string `json:"password,omitempty" required:"true" minLength:"6" maxLength:"6" pattern:"[0-9]+"`
	}
}

type RefreshToken struct {
	Body struct {
		RefreshToken string `json:"refresh_token,omitempty" required:"true"`
	}
}

type Session struct {
	AccessToken  string `json:"access_token" required:"true"`
	RefreshToken string `json:"refresh_token" required:"true"`
}

type cacheItem struct {
	ID    int64 `json:"id"`
	TwoFA bool  `json:"two_fa"`
}

type Auth struct {
	logger      *zap.Logger
	repo        repository.Admin
	cache       cms.Cache
	tokenExpiry time.Duration
	secret      string
	tags        []string
}

func NewAuth(
	repo repository.Admin,
	cache cms.Cache,
	secret string,
	tokenExpiry time.Duration,
	logger *zap.Logger,
) Auth {
	return Auth{
		logger:      logger.Named("auth"),
		repo:        repo,
		cache:       cache,
		secret:      secret,
		tokenExpiry: tokenExpiry,
		tags:        []string{"Auth"},
	}
}

func (r Auth) Register(_ *echo.Echo, humaAPI huma.API) {
	Register(humaAPI, r.signIn, huma.Operation{
		Summary:  "Sign In",
		Method:   http.MethodPost,
		Path:     "/auth/sign-in",
		Tags:     r.tags,
		Security: []map[string][]string{},
		Metadata: map[string]any{
			"target": &cms.CallTarget{
				Access: map[cms.AuthScheme]cms.Decider{
					cms.UnknownScheme: cms.NewDecider(cms.AccessPublic, false),
				},
			},
		},
	})
	Register(humaAPI, r.otp, huma.Operation{
		Summary: "OTP",
		Method:  http.MethodPost,
		Path:    "/auth/otp",
		Tags:    r.tags,
		Metadata: map[string]any{
			"target": &cms.CallTarget{
				Access: map[cms.AuthScheme]cms.Decider{
					cms.BasicScheme: cms.NewDecider(cms.AccessPrivate, false),
					cms.JWTScheme:   cms.NewDecider(cms.AccessPrivate, false),
				},
			},
		},
	})
	Register(humaAPI, r.refreshToken, huma.Operation{
		Summary:  "Refresh Token",
		Method:   http.MethodPost,
		Path:     "/auth/refresh-token",
		Tags:     r.tags,
		Security: []map[string][]string{},
		Metadata: map[string]any{
			"target": &cms.CallTarget{
				Access: map[cms.AuthScheme]cms.Decider{
					cms.UnknownScheme: cms.NewDecider(cms.AccessPublic, false),
				},
			},
		},
	})
}

func (r Auth) signIn(ctx context.Context, in *SignIn) (*Response[Session], error) {
	admin, err := r.repo.FindByEmail(ctx, in.Body.Email)
	if err != nil {
		return nil, r.error(err)
	}

	if err = admin.ValidatePassword(in.Body.Password); err != nil {
		return nil, r.error(err)
	}

	return r.session(ctx, admin, false)
}

func (r Auth) otp(ctx context.Context, in *OTP) (*Response[Session], error) {
	admin := cms.CtxAdmin(ctx)
	if admin == nil {
		return nil, r.error(errors.New("invalid context, admin not found"))
	}

	if err := admin.ValidateOTP(in.Body.Password); err != nil {
		return nil, r.error(err)
	}

	return r.session(ctx, *admin, true)
}

func (r Auth) refreshToken(ctx context.Context, in *RefreshToken) (*Response[Session], error) {
	key := prefixRefreshToken + in.Body.RefreshToken
	var item cacheItem
	if err := r.cache.Get(ctx, key, &item); err != nil {
		return nil, r.error(err)
	}

	admin, err := r.repo.FindByID(ctx, item.ID)
	if err != nil {
		return nil, r.error(err)
	}

	return r.session(ctx, admin, item.TwoFA)
}

func (r Auth) session(ctx context.Context, admin model.Admin, twoFA bool) (*Response[Session], error) {
	tag := fmt.Sprintf("admin:tag:%d", admin.ID)
	_ = r.cache.DelByTag(ctx, tag)

	accessToken, err := cms.NewJWT(
		jwt.MapClaims{"sub": admin.Email, "model": reflect.TypeOf(admin).Name(), "2fa": twoFA},
		admin.Salt+r.secret,
		r.tokenExpiry,
	)
	if err != nil {
		return nil, r.error(err)
	}

	refreshToken := uuid.NewString()
	if err = r.cache.Set(ctx, prefixRefreshToken+refreshToken, &cacheItem{ID: admin.ID, TwoFA: twoFA}, tag); err != nil {
		return nil, r.error(err)
	}

	return &Response[Session]{
		Body: Session{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (r Auth) error(err error) error {
	r.logger.Error("login failed", zap.Error(err))
	return huma.Error400BadRequest("Login failed, please try again")
}
