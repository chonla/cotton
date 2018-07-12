package testsuite

// Task is task for setup and teardown
type Task struct{}

// NewTask to create a new task
func NewTask(title, link string) *Task {
	return &Task{}
}
