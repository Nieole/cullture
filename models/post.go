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
	ID           uuid.UUID    `json:"id" db:"id"`
	ProjectID    uuid.UUID    `json:"-" db:"project_id"`
	Project      *Project     `json:"project,omitempty" belongs_to:"projects"`
	Image        nulls.String `json:"image" db:"image"`
	UserPhone    nulls.String `json:"user_phone,omitempty" db:"user_phone"`
	UserID       nulls.UUID   `json:"-" db:"user_id"`
	User         *User        `json:"user,omitempty" belongs_to:"users"`
	Content      nulls.String `json:"content" db:"content"`
	IsDelete     bool         `json:"is_delete" db:"is_delete"`
	Tags         Tags         `json:"tags" many_to_many:"post_tags" order_by:"created_at desc"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	IsLike       bool         `json:"like" db:"-"`
	LikeCount    int64        `json:"like_count" db:"-"`
	CommentCount int          `json:"comment_count" db:"-"`
	Comments     Comments     `json:"comments" has_many:"comments" order_by:"created_at desc"`
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

//MarshalBinary MarshalBinary
func (p *Posts) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

//UnmarshalBinary UnmarshalBinary
func (p *Posts) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
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
	REDIS.SAdd(fmt.Sprintf("%v:%v:like", (&pop.Model{Value: p}).TableName(), p.ID), user.ID.String())
	REDIS.SAdd(fmt.Sprintf("%v:%v:like", "user", user.ID), p.ID.String())
}

//UnLike UnLike
func (p *Post) UnLike(user *User) {
	REDIS.SRem(fmt.Sprintf("%v:%v:like", (&pop.Model{Value: p}).TableName(), p.ID), user.ID.String())
	REDIS.SRem(fmt.Sprintf("%v:%v:like", "user", user.ID), p.ID.String())
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

//ChangeLike ChangeLike
func (p *Posts) ChangeLike(user *User, tx *pop.Connection, phone string) {
	REDIS.SUnionStore(fmt.Sprintf("%v:%v:like", "user", user.ID), fmt.Sprintf("%v:%v:like", "user", phone))
	REDIS.Del(fmt.Sprintf("%v:%v:like", "user", phone))
	for _, post := range *p {
		post.ChangeLike(user, tx)
	}
	posts := &Posts{}
	if err := tx.All(posts); err == nil {
		for _, post := range *posts {
			key := fmt.Sprintf("%v:%v:like", (&pop.Model{Value: p}).TableName(), post.ID)
			if result, err := REDIS.SIsMember(key, phone).Result(); err == nil {
				if result {
					REDIS.SRem(key, phone)
					REDIS.SAdd(key, user.ID.String())
				}
			}
		}
	}
}

//ChangeLike ChangeLike
func (p *Post) ChangeLike(user *User, tx *pop.Connection) {
	p.UserID = nulls.NewUUID(user.ID)
	err := tx.Update(p)
	if err != nil {
		log.Println(err)
	}
}

//CheckLike CheckLike
func (p *Post) CheckLike(user *User) bool {
	if user == nil {
		return false
	}
	result, err := REDIS.SIsMember(fmt.Sprintf("%v:%v:like", (&pop.Model{Value: p}).TableName(), p.ID), user.ID.String()).Result()
	if err != nil {
		log.Println(fmt.Sprintf("failed SIsMember like %s : %v", user.ID, err))
		return false
	}
	return result
}

//FillLike FillLike
func (p *Posts) FillLike(user *User) Posts {
	out := make(Posts, 0, len(*p))
	for _, post := range *p {
		post.FillLike(user)
		out = append(out, post)
	}
	return out
}

//FillCount FillCount
func (p *Posts) FillCount(tx *pop.Connection) Posts {
	out := make(Posts, 0, len(*p))
	for _, post := range *p {
		post.FillCount(tx)
		out = append(out, post)
	}
	return out
}

//FillCount FillCount
func (p *Post) FillCount(tx *pop.Connection) {
	comments := &Comments{}
	if count, err := tx.Where("post_id = ?", p.ID).Where("is_delete = ?", false).Count(comments); err == nil {
		p.CommentCount = count
	}
}

//FillLike FillLike
func (p *Post) FillLike(user *User) {
	p.IsLike = p.CheckLike(user)
	p.LikeCount = p.CountLike()
}
