package models

import (
	"time"
)

type PageSchema struct {
	Page      string    `json:"page" bson:"page"`
	Key       string    `json:"key" bson:"key"`
	Label     string    `json:"label" bson:"label"`
	Sortable  bool      `json:"sortable" bson:"sortable"`
	Selected  bool      `json:"selected" bson:"selected"`
	CreatedAt time.Time `json:"created_at, omitempty"`
	UpdatedAt time.Time `json:"updated_at, omitempty"`
}
