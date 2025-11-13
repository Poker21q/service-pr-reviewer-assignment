package service

type repository interface{}

type Service struct {
	repo repository
}

func New(repository repository) *Service {
	return &Service{
		repo: repository,
	}
}
