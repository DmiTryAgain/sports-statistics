package statistic

import "sports-statistics/internal/service/entity/user"

type User struct {
	id *user.Id
}

func (u *User) Construct(id *user.Id) *User {
	u.id = id

	return u
}

func (u *User) GetId() *user.Id {
	return u.id
}
