package motion

// Intent describes the movement intent of a living thing.
type Intent uint8

// The set of valid Intents.
const (
	Stay Intent = iota
	StepUp
	StepDown
	StepLeft
	StepRight
)
