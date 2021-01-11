package tasks

import (
	"fmt"
	"io"
)

var tasks = map[string]Task{}

// Task ...
type Task interface {
	fmt.Stringer
	Run() error
	ID() string
	InfoLogFilePath() string
	DetailedLogFilePath() string
	DetailedLog() io.WriteCloser
	InfoLog() io.WriteCloser
}

// Run ...
func Run(tsk Task) {
	go func() {
		err := tsk.Run()
		if err != nil {
			fmt.Fprintf(tsk.InfoLog(), "\nError running task `%s`\n: %s\n\n", tsk.ID(), err.Error())
		} else {
			fmt.Fprintf(tsk.InfoLog(), "\nCompleted task `%s`\n\n", tsk.String())
		}
		delete(tasks, tsk.ID())
	}()
}

// List ...
func List(w io.Writer) {
	for _, task := range tasks {
		fmt.Fprintf(w, "%s\n", task)
	}
}
