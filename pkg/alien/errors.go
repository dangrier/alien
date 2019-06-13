package alien

// Error is a string which satisfies the error interface
type Error string

// Error implements the error interface
func (e Error) Error() string {
	return string(e)
}

// Define error constants
const (
	ErrNotInitialised = Error("alien not initialised")
	ErrProbeNotFound  = Error("probe not found")
)
