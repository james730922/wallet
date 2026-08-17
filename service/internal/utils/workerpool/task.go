package workerpool

type Task interface {
	Exec() error
}
