package connection

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/meteocima/virtual-server/tailor"
	"github.com/meteocima/virtual-server/vpath"
	"golang.org/x/crypto/ssh"
)

func copyLines(proc Process, w io.Writer, outLogFile vpath.VirtualPath) {
	if outLogFile.Host == "localhost" {
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

		return
	}

	cn, err := FindHost(outLogFile.Host)
	if err != nil {
		panic(fmt.Errorf("copyLines to log from %s: FindHost: %w", outLogFile.String(), err))
	}

	conn := cn.(*SSHConnection).client

	cmd, err := conn.NewSession()
	if err != nil {
		panic(fmt.Errorf("copyLines to log from %s: conn.client.NewSession: %w", outLogFile.String(), err))
	}
	defer cmd.Close()

	cmdStr := fmt.Sprintf("tail -F '%s'", outLogFile.Path)

	err = cmd.Start(cmdStr)
	if err != nil {
		panic(fmt.Errorf("copyLines to log from %s: cmd.Start: %w", outLogFile.String(), err))
	}

	go func() {
		_, err := proc.Wait()
		if err != nil {
			panic(fmt.Errorf("copyLines to log from %s: proc.Wait: %w", outLogFile.String(), err))
		}
		err = cmd.Signal(ssh.SIGKILL)
		if err != nil {
			panic(fmt.Errorf("copyLines to log from %s: cmd.Signal(ssh.SIGKILL): %w", outLogFile.String(), err))
		}
	}()

}
