package errors

type Layer int

const (
	LayerStore Layer = iota
	LayerService
	LayerPkg
)
