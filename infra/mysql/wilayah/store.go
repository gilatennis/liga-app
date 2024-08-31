package wilayah

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db}
}

func (s *Store) Get(ctx context.Context, p *common.Pagination, f *repositories.WilayahFilter) ([]entities.Wilayah, error) {
	var wilayahs []entities.Wilayah

	dbHandle := s.db.WithContext(ctx).Where("type = ?", f.Type)

	if f != nil {
		if f.IDProvince != "" {
			dbHandle = dbHandle.Where("id_province = ?", f.IDProvince)
		}
		if f.IDCity != "" {
			dbHandle = dbHandle.Where("id_city = ?", f.IDCity)
		}
		if f.IDSubdistrict != "" {
			dbHandle = dbHandle.Where("id_subdistrict = ?", f.IDSubdistrict)
		}
		if f.IDVillage != "" {
			dbHandle = dbHandle.Where("id_village = ?", f.IDVillage)
		}
		if f.Name != "" {
			dbHandle = dbHandle.Where("name LIKE ?", "%"+f.Name+"%")
		}
	}

	err := dbHandle.Find(&wilayahs).Error
	if err != nil {
		return nil, err
	}

	return wilayahs, nil
}
