package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	MenuName string             `bson:"menu_name,omitempty" json:"menu_name,omitempty"`
	// Image          string             `bson:"image,omitempty" json:"image,omitempty"`
	Price          float64   `bson:"price,omitempty" json:"price,omitempty"`
	Description    string    `bson:"description,omitempty" json:"description,omitempty"`
	Stock          int       `bson:"stock,omitempty" json:"stock,omitempty"`
	MenuCategories string    `bson:"menu_categories,omitempty" json:"menu_categories,omitempty"`
	CreatedAt      time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"` // Field baru untuk waktu pembuatan menu
}
