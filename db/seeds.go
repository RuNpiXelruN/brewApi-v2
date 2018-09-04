package db

import (
	"go_apps/go_api_apps/brewApi-v2/utils"
	"time"

	"github.com/satori/go.uuid"
)

func dropWithSeed() {
	dropDatabase()
	seedDatabase()
}

func dropDatabase() {
	db.DropTableIfExists(&Beer{}, &Brewer{}, &Rank{}, "beer_brewers", &Session{}, &User{})
	db.AutoMigrate(&Beer{}, &Brewer{}, &Rank{}, &Session{}, &User{})
}

func migrateDatabase() {
	db.AutoMigrate(&Beer{}, &Brewer{}, &Rank{}, &Session{}, &User{})
}

func seedDatabase() {
	migrateDatabase()

	beers := []Beer{
		Beer{
			Name:     "Rice to Meet You",
			Status:   "active-empty",
			Featured: false,
			Brewers: []Brewer{
				Brewer{
					FirstName: "justin",
					LastName:  "davidson",
					Featured:  false,
					Rank: &Rank{
						Name:  "level 1 brewmaster",
						Level: 1,
					},
				},
				Brewer{
					FirstName: "luis",
					LastName:  "ramos",
					Featured:  false,
					Rank: &Rank{
						Name:  "level 2 brewmaster",
						Level: 2,
					},
				},
			},
		},
		Beer{
			Name:     "Pineapple Pale Ale",
			Status:   "active-full",
			Featured: false,
			Brewers: []Brewer{
				Brewer{
					FirstName: "thom",
					LastName:  "vincent",
					Featured:  false,
					Rank: &Rank{
						Name:  "level 3 brewmaster",
						Level: 3,
					},
				},
				Brewer{
					FirstName: "alex",
					LastName:  "rapley",
					Featured:  false,
					Rank: &Rank{
						Name:  "level 4 brewmaster",
						Level: 4,
					},
				},
			},
		},
		Beer{
			Name:     "Apricot Wheat Beer",
			Status:   "active-full",
			Featured: false,
			Brewers: []Brewer{
				Brewer{
					FirstName: "shealan",
					LastName:  "forshaw",
					Featured:  false,
					Rank: &Rank{
						Name:  "level 5 brewmaster",
						Level: 5,
					},
				},
				Brewer{
					FirstName: "chris",
					LastName:  "ellis",
					Featured:  false,
					Rank: &Rank{
						Name:  "level 6 brewmaster",
						Level: 6,
					},
				},
			},
		},
		Beer{
			Name:     "Redfern Imperial IPA",
			Status:   "brewing",
			Featured: false,
			Brewers: []Brewer{
				Brewer{
					FirstName: "ronnie",
					LastName:  "pyne",
					Featured:  false,
					Rank: &Rank{
						Name:  "level 7 brewmaster",
						Level: 7,
					},
				},
				Brewer{
					FirstName: "liam",
					LastName:  "fiddler",
					Featured:  false,
					Rank: &Rank{
						Name:  "level 8 brewmaster",
						Level: 8,
					},
				},
			},
		},
	}

	for _, b := range beers {
		db.Save(&b)
	}

	user := User{
		Email:    "justin@mentallyfriendly.com",
		Password: utils.StringPointer("password"),
		Session: &Session{
			Value:   uuid.NewV4().String(),
			Expires: time.Now().Add(1 * time.Hour),
		},
	}

	db.Save(&user)

	// roles := []Role{
	// 	Role{
	// 		Name: "admin",
	// 		Users: []User{
	// 			User{
	// 				Email:    "justin@mentallyfriendly.com",
	// 				Password: utils.StringPointer("password"),
	// 				Session: &Session{
	// 					Value:   uuid.NewV4().String(),
	// 					Expires: time.Now().Add(1 * time.Hour),
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Role{
	// 		Users: []User{
	// 			User{
	// 				Email: "justin@socialplayground.com",
	// 				Session: &Session{
	// 					Value:   uuid.NewV4().String(),
	// 					Expires: time.Now().Add(1 * time.Hour),
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// for _, role := range roles {
	// 	db.Save(&role)
	// }
}
