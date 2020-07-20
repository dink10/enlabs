package processing

// Config keeps configuration of post processing.
type Config struct {
	CancellationTime uint64 `env:"CANCELLATION_TIME,required"`
}
