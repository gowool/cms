package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

var _ repository.Admin = (*AdminRepository)(nil)

type AdminRepository struct {
	Repository[model.Admin, int64]
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{
		Repository[model.Admin, int64]{
			DB:    db,
			Table: "admins",
			SelectColumns: []string{
				"id", "avatar", "email", "role", "salt", "password", "otp", "created", "updated",
			},
			RowScan: func(row interface{ Scan(...any) error }, m *model.Admin) error {
				var (
					role     Role
					otp      OTP
					password Password
				)
				if err := row.Scan(&m.ID, &m.Avatar, &m.Email, &role, &m.Salt, &password, &otp, &m.Created, &m.Updated); err != nil {
					return err
				}

				m.Role = model.Role(role)
				m.Password = model.Password(password)
				m.OTP = model.OTP(otp)
				return nil
			},
			InsertValues: func(m *model.Admin) map[string]any {
				now := time.Now()
				role := Role(m.Role)
				otp := OTP(m.OTP)
				return map[string]any{
					"avatar":   m.Avatar,
					"email":    m.Email,
					"role":     &role,
					"salt":     m.Salt,
					"password": Password(m.Password),
					"otp":      &otp,
					"created":  now,
					"updated":  now,
				}
			},
			UpdateValues: func(m *model.Admin) map[string]any {
				return map[string]any{
					"avatar":   m.Avatar,
					"email":    m.Email,
					"role":     Role(m.Role),
					"salt":     m.Salt,
					"password": Password(m.Password),
					"otp":      OTP(m.OTP),
					"updated":  time.Now(),
				}
			},
		},
	}
}

func (r *AdminRepository) FindByEmail(ctx context.Context, email string) (model.Admin, error) {
	return r.FindBy(ctx, "email", email)
}
