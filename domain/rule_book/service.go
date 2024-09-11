package rule_book

import (
	"context"
	domainErrors "liga-app/domain/errors"
	"liga-app/entities"
	iCommon "liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"
)

type Service struct {
	gormStore iRepo.IGormRepo
}

func NewService(gormStore iRepo.IGormRepo) *Service {
	return &Service{
		gormStore: gormStore,
	}
}

func (s *Service) GetByName(ctx context.Context, name string) (*entities.RuleBook, error) {
	p := iCommon.NewPagination(1, 1, []iCommon.OrderBy{})
	where := make(map[string]interface{})
	where["name"] = name
	rows, _, err := s.gormStore.Get(ctx, p, []entities.RuleBook{}, &where)
	if err != nil {
		return nil, err
	}

	ruleBooks, _ := rows.([]entities.RuleBook)

	if len(ruleBooks) == 0 {
		return nil, &domainErrors.ErrNotFound{Message: "rulebook not found"}
	}

	return &ruleBooks[0], nil
}
