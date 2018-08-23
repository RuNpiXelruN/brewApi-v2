package db

import (
	"time"
)

// BasicBeer struct
type BasicBeer struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// BasicBrewer struct
type BasicBrewer struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Beer struct
type Beer struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name" gorm:"not null;" sql:"unique"`
	Description    string     `json:"description" sql:"default:'A default beer description here'"`
	Status         string     `json:"status" sql:"default:'upcoming'"`
	ImageURL       string     `json:"image_url" sql:"default:'https://placeimg.com/180/400/any'"`
	AlcoholContent float64    `json:"alcohol_content" sql:"default:4.44"`
	Featured       bool       `json:"featured" sql:"default:false"`
	Brewers        []Brewer   `json:"brewers" gorm:"many2many:beer_brewers"`
	CreatedAt      time.Time  `json:"-"`
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `json:"-"`
}

// Brewer struct
type Brewer struct {
	ID        uint       `json:"id"`
	FirstName string     `json:"first_name" gorm:"not null"`
	LastName  string     `json:"last_name"`
	Username  *string    `json:"username" sql:"unique"`
	Featured  bool       `json:"featured" sql:"default:false"`
	Rank      *Rank      `json:"rank"`
	RankID    *uint      `json:"rank_id"`
	Beers     []Beer     `json:"beers" gorm:"many2many:beer_brewers"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

// Rank struct
type Rank struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name" gorm:"not null;" sql:"unique"`
	Level     int        `json:"level" gorm:"not null;" sql:"index:idx_rank_level; unique"`
	Brewers   []Brewer   `json:"brewers"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}
