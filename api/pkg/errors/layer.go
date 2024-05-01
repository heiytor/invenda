package errors

type Layer int

const (
	LayerStore Layer = iota + 1
	LayerService
	LayerRoute
	LayerPkg
)
