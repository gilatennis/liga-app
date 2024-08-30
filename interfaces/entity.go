package interfaces

import "liga-app/entities"

type IEntity interface {
	SetCreatedBy(user *entities.User)
	SetUpdatedBy(user *entities.User)
	SetDeletedBy(user *entities.User)
}
