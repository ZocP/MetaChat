package app

type APP interface {
	OnStart() error
	OnStop() error
}
