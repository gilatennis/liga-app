package court

import (
	"liga-app/entities"
)

type VModelPageCourtList struct {
	User   *entities.User
	Courts []entities.Court
}
