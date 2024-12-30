package ports

type Worker interface {
	Start(id, payload string) error
	Stop(id string) error
}
