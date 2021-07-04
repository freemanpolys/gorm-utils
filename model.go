package gorm

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
// It may be embedded into your model or you may build your own model without it
//    type User struct {
//      gorm.Base
//    }
// https://gorm.io/docs/conventions.html
type Base struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	UpdatedAt time.Time `gorm:"default:now()"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New()
	return
}
