package main

type Car struct {
	ID string `json:"_id,omitempty" bson:"_id,omitempty"`
	Make string `json:"make,omitempty" bson:"make,omitempty"`
	Model string `json:"model,omitempty" bson:"model,omitempty"`
}
