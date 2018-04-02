package ssettings

const (
	WINDOWSIZE          = 1024
	MAX_CLIENT_DEADTIME = 10000000
	KEEP_ALIVE          = false
)

// Settings is a structure used for global server and client settings.
type Settings struct {
	MaxWindowSize    int
	MaxClientTimeout int64
	KeepAlive        bool
}

// GServerSettings is used by client initialization and various functions.
var GServerSettings *Settings = nil

// GetWindowSize returns the max window for the packet buffer.
func (s *Settings) GetWindowSize() int {
	if GServerSettings != nil {
		return s.MaxWindowSize
	}
	return WINDOWSIZE
}

// KeptAlive is used to decide to close the socket on recv or keep it alive.
func (s *Settings) KeptAlive() bool {
	if GServerSettings != nil {
		return s.KeepAlive
	}
	return KEEP_ALIVE
}

// GetClientTimeout returns the client's max timeout time in milliseconds.
func (s *Settings) GetClientTimeout() int64 {
	if GServerSettings != nil {
		return s.MaxClientTimeout
	}
	return MAX_CLIENT_DEADTIME
}

// Initialize the global structure for settings.
func Initialize(windowSize int, clientTimeout int64, keepAlive bool) {
	GServerSettings = &Settings{
		MaxWindowSize:    windowSize,
		MaxClientTimeout: clientTimeout,
		KeepAlive:        keepAlive,
	}
}
