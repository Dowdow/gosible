package runner

type Event any

type ActionStarted struct {
	Name string
}

type ActionCompleted struct {
	Stdout  string
	Stderr  string
	Success bool
}

type Done struct{}

type Failed struct {
	Err error
}
