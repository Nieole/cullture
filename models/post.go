package models

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"time"
)

// Post is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Post struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	ProjectID uuid.UUID    `json:"-" db:"project_id"`
	Project   *Project     `json:"project,omitempty" belongs_to:"projects"`
	Image     nulls.String `json:"image" db:"image"`
	UserPhone string       `json:"user_phone" db:"user_phone"`
	Content   nulls.String `json:"content" db:"content"`
	IsDelete  bool         `json:"is_delete" db:"is_delete"`
	Tags      Tags         `json:"tags" many_to_many:"post_tags" order_by:"created_at desc"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	IsLike    bool         `json:"like,omitempty"`
	IsHate    bool         `json:"hate,omitempty"`
}

// String is not required by pop and may be deleted
func (p Post) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Posts is not required by pop and may be deleted
type Posts []Post

// String is not required by pop and may be deleted
func (p Posts) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

func (p *Posts) FromString(data string) error {
	err := json.Unmarshal([]byte(data), p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) FromString(data string) error {
	err := json.Unmarshal([]byte(data), p)
	if err != nil {
		return err
	}
	return nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Post) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.UserPhone, Name: "UserPhone"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Post) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Post) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (p *Post) Like(phone string) {
	REDIS.SAdd(fmt.Sprintf("%v:%v:like", (&pop.Model{Value: p}).TableName(), p.ID), phone)
	REDIS.SAdd(fmt.Sprintf("%v:%v:like", "user", phone), p.ID.String())
}

func (p *Post) Hate(phone string) {
	REDIS.SAdd(fmt.Sprintf("%v:%v:hate", p.ID, (&pop.Model{Value: p}).TableName()), phone)
	REDIS.SAdd(fmt.Sprintf("%v:%v:hate", "user", phone), p.ID.String())
}
