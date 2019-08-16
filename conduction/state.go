package conduction

// State represent the temperature (T) and heat flux (Q) inside an element
type State struct {
	temperature *float64
}

func NewState() *State {
	return &State{temperature: nil}
}

func (s *State) SetT(value float64) {
	*s.temperature = value
}

func (s *State) T() float64 {
	return *s.temperature
}
