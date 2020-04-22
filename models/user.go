package models

import (
	"culture/cache"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type User struct {
	ID                   uuid.UUID    `json:"id" db:"id"`
	IsDelete             bool         `json:"is_delete" db:"is_delete"`
	Name                 string       `json:"name" db:"name"`
	LoginName            string       `json:"login_name" db:"login_name"`
	PasswordHash         nulls.String `json:"-" db:"password_hash"`
	Password             nulls.String `json:"password" db:"-"`
	PasswordConfirmation nulls.String `json:"password_confirmation" db:"-"`
	Avatar               nulls.String `json:"avatar" db:"avatar"`
	Sex                  int          `json:"sex" db:"sex"`
	Type                 int          `json:"-" db:"type"`
	Birthday             nulls.Time   `json:"birthday" db:"birthday"`
	Introduction         nulls.String `json:"introduction" db:"introduction"`
	Background           nulls.String `json:"background" db:"background"`
	Role                 nulls.String `json:"role" db:"role"`
	IsActive             bool         `json:"is_active" db:"is_active"`
	CreatedAt            time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at" db:"updated_at"`
	Posts                Posts        `json:"posts,omitempty" has_many:"posts" order_by:"created_at desc"`
	Comments             Comments     `json:"comments" has_many:"comments" order_by:"create_at desc"`
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
		&validators.StringLengthInRange{Field: u.LoginName, Min: 3, Max: 100, Name: "LoginName", Message: "账号长度不符合要求"},
		&validators.FuncValidator{
			Fn: func() bool {
				if u.Password.Valid {
					return u.Password.String == u.PasswordConfirmation.String
				}
				return true
			},
			Field:   "Password",
			Name:    "Password",
			Message: "密码不一致",
		},
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

//BeforeUpdate BeforeUpdate
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

//AfterUpdate AfterUpdate
func (u *User) AfterUpdate(tx *pop.Connection) error {
	err := cache.Clean(fmt.Sprintf("cache:user:%v", u.ID))
	if err != nil {
		log.Printf("clean cache failed : %v", err)
	}
	return nil
}

//Load Load
func (u *User) Load(tx *pop.Connection, expiration time.Duration) error {
	return cache.Once(fmt.Sprintf("cache:user:%v", u.ID), u, func() (interface{}, error) {
		err := tx.Find(u, u.ID)
		if err != nil {
			return nil, err
		}
		return u, err
	}, expiration)
}
