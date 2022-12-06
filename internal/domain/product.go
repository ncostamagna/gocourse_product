package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID        string         `json:"id" gorm:"type:char(36);not null;primary_key;unique_index"`
	Name      string         `json:"name" gorm:"type:char(50);not null"`
	Price     float64        `json:"price"`
	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	Deleted   gorm.DeletedAt `json:"-"`
}

func (c *Product) BeforeCreate(tx *gorm.DB) (err error) {

	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
