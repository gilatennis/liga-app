package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
)

type WilayahFilter struct {
	IDProvince    string
	IDCity        string
	IDSubdistrict string
	IDVillage     string
	Type          int
	Name          string
}

type IWilayahRepository interface {
	Get(ctx context.Context, p *common.Pagination, f *WilayahFilter) ([]entities.Wilayah, error)
}
