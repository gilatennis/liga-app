package entities

type Venue struct {
	Base
	Name    string
	Slug    string
	Address string
	Lat     float64 `gorm:"type:decimal(17,15)"`
	Lon     float64 `gorm:"type:decimal(18,15)"`
	Courts  []Court `gorm:"ForeignKey:VenueUUID"`
	Image   string
}
