package models

import (
	"encoding/json"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type User struct {
	ID                   uuid.UUID    `json:"id" db:"id"`
	IsDelete             bool         `json:"is_delete" db:"is_delete"`
	Name                 string       `json:"name" db:"name"`
	LoginName            string       `json:"login_name" db:"login_name"`
	PasswordHash         nulls.String `json:"-" db:"password_hash"`
	Password             nulls.String `json:"password" db:"-"`
	PasswordConfirmation string       `json:"passwordConfirmation" db:"-"`
	Avatar               nulls.String `json:"avatar" db:"avatar"`
	Sex                  int          `json:"sex" db:"sex"`
	Birthday             nulls.Time   `json:"birthday" db:"birthday"`
	Introduction         nulls.String `json:"introduction" db:"introduction"`
	Background           nulls.String `json:"background" db:"background"`
	IsActive             bool         `json:"is_active" db:"is_active"`
	CreatedAt            time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at" db:"updated_at"`
	Posts                Posts        `json:"posts,omitempty" has_many:"posts" order_by:"created_at desc"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.StringIsPresent{Field: u.LoginName, Name: "LoginName"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (u *User) BeforeUpdate(tx *pop.Connection) error {
	if u.Password.Valid {
		ph, err := bcrypt.GenerateFromPassword([]byte(u.Password.String), bcrypt.DefaultCost)
		if err != nil {
			return errors.WithStack(err)
		}
		u.PasswordHash = nulls.NewString(string(ph))
	}
	return nil
}
