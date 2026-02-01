package domain

import "time"

type Bucket struct {
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creation_date"`
}
