package pg

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/gowool/cms/internal"
	"github.com/gowool/cms/model"
)

type Metas []model.Meta

func (m Metas) Scan(src any) error {
	switch src := src.(type) {
	case string:
		return json.Unmarshal(internal.Bytes(src), &m)
	case []byte:
		return json.Unmarshal(src, &m)
	default:
		return errors.New("invalid src type for Metas")
	}
}

func (m Metas) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "[]", nil
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return internal.String(raw), nil
}

type StrMap map[string]string

func (m StrMap) Scan(src any) error {
	switch src := src.(type) {
	case string:
		return json.Unmarshal(internal.Bytes(src), &m)
	case []byte:
		return json.Unmarshal(src, &m)
	default:
		return errors.New("invalid src type for StrMap")
	}
}

func (m StrMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return internal.String(raw), nil
}

type Role model.Role

func (r *Role) Scan(src any) error {
	switch src := src.(type) {
	case string:
		*r = Role(model.NewRole(src))
		return nil
	default:
		return errors.New("invalid src type for Role")
	}
}

func (r *Role) Value() (driver.Value, error) {
	if r == nil {
		return model.RoleGuest.String(), nil
	}
	return model.Role(*r).String(), nil
}

type Password model.Password

func (p *Password) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		*p = make(Password, len(src))
		copy(*p, src)
		return nil
	default:
		return errors.New("invalid src type for Password")
	}
}

func (p *Password) Value() (driver.Value, error) {
	if p == nil {
		return nil, errors.New("password is nil")
	}
	dst := make([]byte, len(*p))
	copy(dst, *p)
	return dst, nil
}

type OTP model.OTP

func (otp *OTP) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		copy(otp[:], src)
		return nil
	default:
		return errors.New("invalid src type for OTP")
	}
}

func (otp *OTP) Value() (driver.Value, error) {
	dst := make([]byte, len(otp))
	copy(dst, otp[:])
	return dst, nil
}
