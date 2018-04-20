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
			Name:     "Apricot Wheat Beer",
			Status:   "upcoming",
			Featured: true,
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
					FirstName: "annabelle",
					LastName:  "davidson",
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
			Status:   "brewing",
			Featured: false,
			Brewers: []Brewer{
				Brewer{
					FirstName: "sawyer",
					LastName:  "davidson",
					Featured:  true,
					Rank: &Rank{
						Name:  "level 3 brewmaster",
						Level: 3,
					},
				},
				Brewer{
					FirstName: "brooks",
					LastName:  "davidson",
					Featured:  true,
					Rank: &Rank{
						Name:  "level 4 brewmaster",
						Level: 4,
					},
				},
			},
		},
		Beer{
			Name:     "Rice to Meet You",
			Status:   "active-full",
			Featured: false,
			Brewers: []Brewer{
				Brewer{
					FirstName: "pete",
					LastName:  "smith",
					Featured:  false,
					Rank: &Rank{
						Name:  "level 5 brewmaster",
						Level: 5,
					},
				},
				Brewer{
					FirstName: "jennie",
					LastName:  "morton",
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
			Status:   "past",
			Featured: false,
			Brewers: []Brewer{
				Brewer{
					FirstName: "michael",
					LastName:  "davidson",
					Featured:  false,
					Rank: &Rank{
						Name:  "level 7 brewmaster",
						Level: 7,
					},
				},
				Brewer{
					FirstName: "jenny",
					LastName:  "davidson",
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

func seedBrewers() {
	brewers := []Brewer{
		Brewer{
			FirstName: "justin",
		},
		Brewer{
			FirstName: "annabelle",
		},
		Brewer{
			FirstName: "sawyer",
		},
		Brewer{
			FirstName: "brooks",
		},
	}
	for _, b := range brewers {
		db.Save(&b)
	}
}

func seedRanks() {
	ranks := []Rank{
		Rank{
			Name:  "Level 1 Brewmaster",
			Level: 1,
		},
		Rank{
			Name:  "Level 2 Brewmaster",
			Level: 2,
		},
		Rank{
			Name:  "Level 3 Brewmaster",
			Level: 3,
		},
		Rank{
			Name:  "Level 4 Brewmaster",
			Level: 4,
		},
		Rank{
			Name:  "Level 5 Brewmaster",
			Level: 5,
		},
	}
	for _, r := range ranks {
		db.Save(&r)
	}
}
