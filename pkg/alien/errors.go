package alien

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotInitialised = Error("alien not initialised")
	ErrProbeNotFound  = Error("probe not found")
)
