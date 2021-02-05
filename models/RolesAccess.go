package models

import (
	"time"
)

type RolesAccess struct {
	Add       bool      `json:"add" bson:"add"`
	Read      bool      `json:"read" bson:"read"`
	Update    bool      `json:"update" bson:"update"`
	Delete    bool      `json:"detele" bson:"detele"`
	IDExtra   string    `json:"id_extra" bson:"id_extra"`
	NameModel string    `json:"name_model" bson:"name_model"`
	CreatedAt time.Time `json:"created_at, omitempty"`
	UpdatedAt time.Time `json:"updated_at, omitempty"`
}
