package models

type ModelDevice struct {
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
}
