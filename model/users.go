package model

import "go.mongodb.org/mongo-driver/bson/primitive"


type Users struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Role        string             `json:"role,omitempty" bson:"role,omitempty"`
	Username    string             `json:"username,omitempty" bson:"username,omitempty" gorm:"unique;not null"`
	Password    string             `json:"password,omitempty" bson:"password,omitempty"`
	PhoneNumber string             `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
}
