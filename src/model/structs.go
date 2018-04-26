package model

import "time"

// Beer struct
type Beer struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name" gorm:"not null;" sql:"unique"`
	Description    string     `json:"description" sql:"default:'A default beer description here'"`
	Status         string     `json:"status" sql:"default:'upcoming'"`
	ImageURL       string     `json:"image_url;" sql:"default:'https://placeimg.com/180/400/any'"`
	AlcoholContent float64    `json:"alcohol_content;" sql:"default:4.44"`
	Featured       bool       `json:"featured;" sql:"default:false"`
	Brewers        []Brewer   `json:"brewers" gorm:"many2many:beer_brewers"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

// Brewer struct
type Brewer struct {
	ID        uint       `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Featured  bool       `json:"featured" sql:"default:false"`
	Rank      *Rank      `json:"rank"`
	RankID    *uint      `json:"rank_id"`
	Beers     []Beer     `json:"beers" gorm:"many2many:beer_brewers"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Rank struct
type Rank struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name" gorm:"not null;" sql:"unique"`
	Level     int        `json:"level" gorm:"not null;" sql:"index:idx_rank_level; unique"`
	Brewers   []Brewer   `json:"brewers"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
