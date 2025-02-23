package factory

type Factory struct{}

func Default() *Factory {
	return &Factory{}
}
