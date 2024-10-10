package cms

import (
	"context"

	"github.com/gomig/avatar"

	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type AdminService struct {
	repo repository.Admin
}

func NewAdminService(repo repository.Admin) *AdminService {
	return &AdminService{
		repo: repo,
	}
}

func (s *AdminService) Create(ctx context.Context, email, password, issuer string) (model.Admin, string, error) {
	pswd, err := model.NewPassword(password)
	if err != nil {
		return model.Admin{}, "", err
	}

	otp, err := model.NewOTP()
	if err != nil {
		return model.Admin{}, "", err
	}

	male := avatar.NewPersonAvatar(true)
	male.RandomizeShape(avatar.Circle)

	admin := model.Admin{
		Avatar:   male.SVG(),
		Email:    email,
		Password: pswd,
		OTP:      otp,
		Role:     model.RoleAdmin,
	}
	admin = admin.WithRandomSalt()

	key, err := admin.OTPKey(issuer)
	if err != nil {
		return model.Admin{}, "", err
	}

	if err = s.repo.Create(ctx, &admin); err != nil {
		return model.Admin{}, "", err
	}
	return admin, key, nil
}

func (s *AdminService) ChangePassword(ctx context.Context, email, password string) error {
	admin, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	admin.Password, err = model.NewPassword(password)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, &admin)
}

func (s *AdminService) ChangeRole(ctx context.Context, email string, role model.Role) error {
	admin, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	admin.Role = role

	return s.repo.Update(ctx, &admin)
}

func (s *AdminService) GetOTPKey(ctx context.Context, email, issuer string, newOTP bool) (string, error) {
	admin, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if newOTP {
		admin.OTP, err = model.NewOTP()
		if err != nil {
			return "", err
		}

		if err = s.repo.Update(ctx, &admin); err != nil {
			return "", err
		}
	}

	return admin.OTPKey(issuer)
}
