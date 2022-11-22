package statistic

type Statistic struct {
	id       *Id
	training *Training
	user     *User
	created  *Created
	count    *Count
	sets     *Sets
}

func (s *Statistic) Construct(
	id *Id,
	training *Training,
	user *User,
	created *Created,
	count *Count,
	sets *Sets,
) *Statistic {
	s.id = id
	s.training = training
	s.user = user
	s.created = created
	s.count = count
	s.sets = sets

	return s
}

func (s *Statistic) GetId() *Id {
	return s.id
}

func (s *Statistic) GetTraining() *Training {
	return s.training
}

func (s *Statistic) GetUser() *User {
	return s.user
}

func (s *Statistic) GetCreated() *Created {
	return s.created
}

func (s *Statistic) GetCount() *Count {
	return s.count
}

func (s *Statistic) GetsSets() *Sets {
	return s.sets
}
