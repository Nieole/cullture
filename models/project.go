package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// Project is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Project struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	Name           string        `json:"name" db:"name"`
	RemoteID       nulls.String  `json:"remote_id" db:"remote_id"`
	IsDelete       bool          `json:"is_delete" db:"is_delete"`
	Posts          Posts         `json:"posts,omitempty" has_many:"posts" order_by:"created_at desc"`
	OrganizationID uuid.UUID     `json:"-" db:"organization_id"`
	Organization   *Organization `json:"organization,omitempty" belongs_to:"organizations"`
	Introduction   nulls.String  `json:"introduction" db:"introduction"`
	RegionCode     nulls.String  `json:"region_code" db:"region_code"`
	Address        nulls.String  `json:"address" db:"address"`
	Longitude      string        `json:"longitude" db:"-"`
	Latitude       string        `json:"latitude" db:"-"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (p Project) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Projects is not required by pop and may be deleted
type Projects []Project

// String is not required by pop and may be deleted
func (p Projects) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

//FromString FromString
func (p *Project) FromString(data string) error {
	err := json.Unmarshal([]byte(data), p)
	if err != nil {
		return err
	}
	return nil
}

//FromString FromString
func (p *Projects) FromString(data string) error {
	err := json.Unmarshal([]byte(data), p)
	if err != nil {
		return err
	}
	return nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Project) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Project) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Project) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
