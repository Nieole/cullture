package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// Post is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Post struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	ProjectID uuid.UUID    `json:"-" db:"project_id"`
	Project   *Project     `json:"project,omitempty" belongs_to:"projects"`
	Image     nulls.String `json:"image" db:"image"`
	UserID    uuid.UUID    `json:"-" db:"user_id"`
	User      *User        `json:"user,omitempty" belongs_to:"users"`
	Content   nulls.String `json:"content" db:"content"`
	IsDelete  bool         `json:"is_delete" db:"is_delete"`
	Tags      Tags         `json:"tags" many_to_many:"post_tags" order_by:"created_at desc"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	IsLike    bool         `json:"like" db:"-"`
	LikeCount int64        `json:"like_count" db:"-"`
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

//FromString FromString
func (p *Posts) FromString(data string) error {
	err := json.Unmarshal([]byte(data), p)
	if err != nil {
		return err
	}
	return nil
}

//FromString FromString
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
	return validate.Validate(), nil
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

//Like Like
func (p *Post) Like(user *User) {
	REDIS.SAdd(fmt.Sprintf("%v:%v:like", (&pop.Model{Value: p}).TableName(), p.ID), user.ID)
	REDIS.SAdd(fmt.Sprintf("%v:%v:like", "user", user.ID), p.ID)
}

//UnLike UnLike
func (p *Post) UnLike(user *User) {
	REDIS.SRem(fmt.Sprintf("%v:%v:like", (&pop.Model{Value: p}).TableName(), p.ID), user.ID)
	REDIS.SRem(fmt.Sprintf("%v:%v:like", "user", user.ID), p.ID)
}

//CountLike CountLike
func (p *Post) CountLike() int64 {
	result, err := REDIS.SCard(fmt.Sprintf("%v:%v:like", (&pop.Model{Value: p}).TableName(), p.ID)).Result()
	if err != nil {
		log.Println(fmt.Sprintf("failed scard like : %v", err))
		return 0
	}
	return result
}

//CheckLike CheckLike
func (p *Post) CheckLike(user *User) bool {
	result, err := REDIS.SIsMember(fmt.Sprintf("%v:%v:like", (&pop.Model{Value: p}).TableName(), p.ID), user.ID).Result()
	if err != nil {
		log.Println(fmt.Sprintf("failed SIsMember like %s : %v", user.ID, err))
		return false
	}
	return result
}

//Fill Fill
func (p *Posts) Fill(user *User) Posts {
	out := make(Posts, 0, len(*p))
	for _, post := range *p {
		post.Fill(user)
		out = append(out, post)
	}
	return out
}

//Fill Fill
func (p *Post) Fill(user *User) {
	p.IsLike = p.CheckLike(user)
	p.LikeCount = p.CountLike()
}
