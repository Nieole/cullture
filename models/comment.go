package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// Comment is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Comment struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	IsDelete  bool         `json:"is_delete" db:"is_delete"`
	Content   nulls.String `json:"content" db:"content"`
	CommentID nulls.UUID   `json:"-" db:"comment_id"`
	Comment   *Comment     `json:"parent" belongs_to:"comments"`
	Comments  Comments     `json:"children" has_many:"comments" order_by:"created_at desc"`
	UserID    uuid.UUID    `json:"-" db:"user_id"`
	User      *User        `json:"user" belongs_to:"users"`
	PostID    uuid.UUID    `json:"-" db:"post_id"`
	Post      *Post        `json:"post" belongs_to:"posts"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (c Comment) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Comments is not required by pop and may be deleted
type Comments []Comment

// String is not required by pop and may be deleted
func (c *Comments) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

//FromString FromString
func (c *Comment) FromString(data string) error {
	err := json.Unmarshal([]byte(data), c)
	if err != nil {
		return err
	}
	return nil
}

//FromString FromString
func (c *Comments) FromString(data string) error {
	err := json.Unmarshal([]byte(data), c)
	if err != nil {
		return err
	}
	return nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Comment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Comment) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Comment) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
