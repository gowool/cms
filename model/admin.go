package model

import (
	"math"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/gowool/cms/internal"
)

type Role int8

const (
	RoleGuest Role = iota
	RoleReader
	RoleWriter
	RoleAdmin = Role(math.MaxInt8)
)

func (b Role) IsValid() bool {
	return b&^(RoleReader|RoleWriter|RoleAdmin) == 0
}

func NewRole(r string) Role {
	switch r {
	case "reader":
		return RoleReader
	case "writer":
		return RoleWriter
	case "admin":
		return RoleAdmin
	default:
		return RoleGuest
	}
}

func (b Role) String() string {
	switch b {
	case RoleReader:
		return "reader"
	case RoleWriter:
		return "writer"
	case RoleAdmin:
		return "admin"
	default:
		return "guest"
	}
}

type Admin struct {
	ID       int64     `json:"id,omitempty" required:"true"`
	Avatar   string    `json:"avatar,omitempty" required:"true"`
	Email    string    `json:"email,omitempty" required:"true" format:"email"`
	Role     Role      `json:"role,omitempty" required:"true"`
	Salt     string    `json:"_" hidden:"true"`
	Password Password  `json:"-" hidden:"true"`
	OTP      OTP       `json:"-" hidden:"true"`
	Created  time.Time `json:"created,omitempty" required:"true"`
	Updated  time.Time `json:"updated,omitempty" required:"true"`
}

func (a Admin) GetID() int64 {
	return a.ID
}

func (a Admin) WithRandomSalt() Admin {
	a.Salt = internal.RandomString(50)
	return a
}

func (a Admin) ValidatePassword(password string) error {
	return a.Password.Validate(password)
}

func (a Admin) ValidateOTP(password string) error {
	return a.OTP.Validate(password)
}

func (a Admin) OTPKey(issuer string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: a.Email,
		Period:      30,
		SecretSize:  otpSize,
		Secret:      a.OTP[:],
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", err
	}
	return key.String(), nil
}
