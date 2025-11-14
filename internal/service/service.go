package service

type Service struct {
	storage   storage
	txManager txManager
}

func MustNew(storage storage, txManager txManager) *Service {
	return &Service{
		storage:   storage,
		txManager: txManager,
	}
}
