package piscine

type Availability struct {
	Day          string    `json:"day"`
	OpeningHours []Opening `json:"opening-hours"`
}

type Opening struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}
