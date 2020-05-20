package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// Banner is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Banner struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	Content   nulls.String `json:"content" db:"content"`
	Title     nulls.String `json:"title" db:"title"`
	Sort      int          `json:"sort" db:"sort"`
	IsDelete  bool         `json:"is_delete" db:"is_delete"`
	Target    nulls.String `json:"target" db:"target"`
	Source    nulls.String `json:"source" db:"source"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (b Banner) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Banners is not required by pop and may be deleted
type Banners []Banner

// String is not required by pop and may be deleted
func (b Banners) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *Banner) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *Banner) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *Banner) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
