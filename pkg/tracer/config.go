package tracer

type Config struct {
	ServiceName        string
	Environment        string
	Endpoint           string
	TraceRatioFraction float64
}
