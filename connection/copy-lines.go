package connection

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/meteocima/virtual-server/tailor"
	"github.com/meteocima/virtual-server/vpath"
)

func copyLines(proc Process, w io.Writer, outLogFile vpath.VirtualPath) {
	var logFile *os.File
	var err error = errors.New("empty")
	for err != nil {
		logFile, err = os.Open(outLogFile.Path)
		if os.IsNotExist(err) {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: copyLines error: (os.Open `%s`\n): %s", outLogFile.Path, err.Error())
			return
		}
	}

	go func() {
		tailProc := tailor.New(logFile, w, 1024)
		errs := tailProc.Start()
		proc.Wait()
		tailProc.Stop()
		err := <-errs
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: copyLines error (reading lines from `%s`): %s\n", outLogFile.Path, err.Error())
		}
	}()

}
