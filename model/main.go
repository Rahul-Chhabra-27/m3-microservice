package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var OrderCollection *mongo.Collection

type Order struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Price     int64              `bson:"price,omitempty"`
	Quantity  int64              `bson:"quantity,omitempty"`
	CreatedAt primitive.DateTime `bson:"createdAt,omitempty"`
	UpdatedAt primitive.DateTime `bson:"updatedAt,omitempty"`
}
