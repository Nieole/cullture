package models

import (
	"culture/cache"
	"culture/work"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/buffalo/worker"

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
	PostsCount     int           `json:"posts_count,omitempty"`
	OrganizationID uuid.UUID     `json:"-" db:"organization_id"`
	Organization   *Organization `json:"organization,omitempty" belongs_to:"organizations"`
	Introduction   nulls.String  `json:"introduction" db:"introduction"`
	RegionCode     nulls.String  `json:"region_code" db:"region_code"`
	Address        nulls.String  `json:"address" db:"address"`
	Country        nulls.String  `json:"country" db:"country"`
	Province       nulls.String  `json:"province" db:"province"`
	City           nulls.String  `json:"city" db:"city"`
	District       nulls.String  `json:"district" db:"district"`
	Type           string        `json:"type" db:"type"`
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

//AfterSave AfterSave
func (p *Project) AfterSave(tx *pop.Connection) error {
	err := cache.Clean(fmt.Sprintf("cache:project:%v", p.ID))
	if err != nil {
		log.Printf("clean cache failed : %v", err)
	}
	err = cache.Clean("cache:project_geo")
	if err != nil {
		log.Printf("clean cache failed : %v", err)
	}
	work.W.Perform(worker.Job{
		Handler: "update_project",
	})
	return nil
}

//Count 统计项目数量
func (p *Projects) Count() (int, error) {
	return DB.Where("is_delete = ?", false).Count(p)
}
