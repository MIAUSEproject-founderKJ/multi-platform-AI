//boot/probe/type_struct.go


type Probe interface {
	Name() string
	Run(ctx context.Context) (any, error)
}



type ProbeResult[T any] struct {
	Value    T
	Error    error
	Duration time.Duration
	Source   string
}

