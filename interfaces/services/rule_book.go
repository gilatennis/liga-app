package services

import (
	"context"
	"liga-app/entities"
)

type IRuleBookService interface {
	GetByName(ctx context.Context, name string) (*entities.RuleBook, error)
}
