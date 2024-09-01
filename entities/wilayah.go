package entities

/**
* Type 1 = Province
* Type 2 = City
* Type 3 = Subdistrict
* Type 4 = Village
 */
type Wilayah struct {
	IDProvince    string `gorm:"type:varchar(24);primaryKey;uniqueIndex:idx_wilayah"`
	IDCity        string `gorm:"type:varchar(24);primaryKey;uniqueIndex:idx_wilayah"`
	IDSubdistrict string `gorm:"type:varchar(24);primaryKey;uniqueIndex:idx_wilayah"`
	IDVillage     string `gorm:"type:varchar(24);primaryKey;uniqueIndex:idx_wilayah"`
	Name          string `gorm:"type:varchar(256)"`
	Type          int
}
