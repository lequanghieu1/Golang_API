package models

type EventCode struct {
	Code int    `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}
