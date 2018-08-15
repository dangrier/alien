package probe

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotInitialised            = Error("probe not initialised")
	ErrNotRunning                = Error("probe not running")
	ErrNotStopped                = Error("probe not stopped")
	ErrInvalidEndpoint           = Error("probe invalid: endpoint")
	ErrInvalidMethod             = Error("probe invalid: method")
	ErrInvalidFrequencyZero      = Error("probe invalid: frequency is zero")
	ErrInvalidSuccessFilterEmpty = Error("probe invalid: no success filter")
	ErrFilterAlreadySet          = Error("probe with success filter: filter already set")
)
