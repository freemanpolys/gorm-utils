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
type Model struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (model *Model) BeforeCreate(tx *gorm.DB) (err error) {
	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return
}

// BeforeUpdate will update UpdatedAt.
func (model *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	model.UpdatedAt = time.Now()
	return
}
