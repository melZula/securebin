package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

// Data ...
type Data struct {
	ID                int    `json:"id"`
	Img               string `json:"img"`
	Password          string `json:"password,omitempty"` // omitempty: if empty - don't return
	EncryptedPassword string `json:"-"`                  // no render
}

// Validate ...
func (u *Data) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Password, validation.Required.When(u.EncryptedPassword == ""), validation.Length(6, 100)),
	)
}

// BeforeCreate ...
func (u *Data) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}
	return nil
}

// Sanitize ...
func (u *Data) Sanitize() {
	u.Password = ""
}

// ComparePassword ...
func (u *Data) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
