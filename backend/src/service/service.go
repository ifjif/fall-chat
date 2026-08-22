package service

type Service interface {
	Service()
	Name() string
	Inject()
}
