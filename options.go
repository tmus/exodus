package exodus

type direction int

const (
	Unknown direction = iota
	Up
	Down
)

func directionFromString(dir string) direction {
	if dir == "up" {
		return Up
	}

	if dir == "down" {
		return Down
	}

	return Unknown
}

type Options struct {
	direction direction
}

func (o Options) Direction() direction {
	return o.direction
}
