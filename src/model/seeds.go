package model

func migrateDB() {
	db.DropTableIfExists(&Beer{}, &Brewer{}, &Rank{}, "beer_brewers")
	db.AutoMigrate(&Beer{}, &Brewer{}, &Rank{})
}

func seedAll() {
	seedBeersBrewersRanks()
}

func seedBeersBrewersRanks() {
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
}
