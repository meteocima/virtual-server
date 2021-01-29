package connection

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/meteocima/virtual-server/vpath"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SSHConnection ...
type SSHConnection struct {
	BackupHosts []string
	Name        string

	Host     string
	Port     int
	User     string
	KeyPath  string
	hostName string
	config   *ssh.ClientConfig
	client   *ssh.Client
}

// HostName ...
func (conn *SSHConnection) HostName() string {
	return conn.hostName
}

func privateSSHKey(path string) (ssh.AuthMethod, error) {
	privateKey, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)

	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}

type sshReader struct {
	client *sftp.Client
	reader io.ReadCloser
}

func (r sshReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}
func (r sshReader) Close() error {
	err := r.reader.Close()
	if err != nil {
		r.client.Close()
		return err
	}

	return r.client.Close()
}

// OpenReader ...
func (conn *SSHConnection) OpenReader(file vpath.VirtualPath) (io.ReadCloser, error) {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return nil, err
	}

	reader, err := client.Open(file.Path)
	if err != nil {
		return nil, err
	}
	return sshReader{client, reader}, nil
}

const maxRetries = 5

// Open ...
func (conn *SSHConnection) Open() error {
	fmt.Println("OPEN CN")
	conn.config = &ssh.ClientConfig{
		User:            conn.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 5,
	}

	key, err := privateSSHKey(conn.KeyPath)
	if err != nil {
		return fmt.Errorf("Open error: cannot read ssh key %s: %w", conn.KeyPath, err)
	}
	conn.config.Auth = []ssh.AuthMethod{key}

	retryCount := 0
	failed := true

	for failed && retryCount < maxRetries {
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", conn.Host, conn.Port), conn.config)
		if err == nil && retryCount > 0 && len(conn.BackupHosts) > 0 {
			fmt.Printf(
				"VPE: successfully connected to server `%s` using hostname %s. Subsequents requests will use hostname %s\n",
				conn.Name,
				conn.Host,
				conn.Host,
			)
		}
		fmt.Printf(
			"VPE: successfully connected to server `%s` using hostname %s. Subsequents requests will use hostname %s\n",
			conn.Name,
			conn.Host,
			conn.Host,
		)
		if failed = err != nil; failed {
			retryHost := conn.Host
			failedHost := conn.Host
			if conn.BackupHosts != nil {
				retryHost = conn.BackupHosts[0]
				conn.Host = conn.BackupHosts[0]
				conn.BackupHosts = append(conn.BackupHosts[1:], failedHost)
			}
			fmt.Printf(
				"VPE: cannot connect to server `%s` using hostname %s: %v\nThe operation will be retried in 10 seconds on %s\n",
				conn.Name,
				failedHost,
				err,
				retryHost,
			)
			time.Sleep(10 * time.Second)
			retryCount++
		}

		conn.client = client
	}

	if err != nil {
		return fmt.Errorf("cannot dial ssh server %s: %w", conn.Host, err)
	}
	fmt.Println("CONNECTED")
	return nil
}

// Close ...
func (conn *SSHConnection) Close() error {
	return conn.client.Close()
}

type sshWriter struct {
	client *sftp.Client
	writer io.WriteCloser
}

func (r sshWriter) Write(p []byte) (n int, err error) {
	return r.writer.Write(p)
}

func (r sshWriter) Close() error {
	err := r.writer.Close()
	if err != nil {
		r.client.Close()
		return err
	}

	return r.client.Close()
}

// OpenWriter ...
func (conn *SSHConnection) OpenWriter(file vpath.VirtualPath) (io.WriteCloser, error) {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return nil, err
	}

	writer, err := client.OpenFile(file.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY)
	if err != nil {
		return nil, err
	}
	return sshWriter{client, writer}, nil
}

// ReadDir ...
func (conn *SSHConnection) ReadDir(dir vpath.VirtualPath) (vpath.VirtualPathList, error) {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	files, err := client.ReadDir(dir.Path)
	if err != nil {
		return nil, fmt.Errorf("ReadDir `%s`: sftp.ReadDir: %w", dir.String(), err)
	}
	filenames := make(vpath.VirtualPathList, len(files))
	for i, f := range files {
		filenames[i] = dir.Join(f.Name())
	}
	sort.Sort(filenames)
	return filenames, nil
}

// Stat ...
func (conn *SSHConnection) Stat(path vpath.VirtualPath) (os.FileInfo, error) {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.Stat(path.Path)
}

// Link ...
func (conn *SSHConnection) Link(source, target vpath.VirtualPath) error {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Symlink(source.Path, target.Path)
}

// MkDir ...
func (conn *SSHConnection) MkDir(dir vpath.VirtualPath) error {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return fmt.Errorf("MkDir `%s`: sftp.NewClient: %w", dir.String(), err)
	}
	defer client.Close()

	err = client.MkdirAll(dir.Path)
	if err != nil {
		return fmt.Errorf("MkDir `%s`: sftp.MkdirAll: %w", dir.String(), err)
	}

	return nil
}

// RmDir ...
func (conn *SSHConnection) RmDir(dir vpath.VirtualPath) error {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return fmt.Errorf("RmDir `%s`: sftp.NewClient: %w", dir.String(), err)
	}
	defer client.Close()

	err = client.RemoveDirectory(dir.Path)
	if err != nil {
		return fmt.Errorf("RmDir `%s`: sftp.RemoveDirectory: %w", dir.String(), err)
	}

	return nil
}

// RmFile ...
func (conn *SSHConnection) RmFile(file vpath.VirtualPath) error {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return fmt.Errorf("RmFile `%s`: sftp.NewClient: %w", file.String(), err)
	}
	defer client.Close()
	err = client.Remove(file.Path)
	if err != nil {
		return fmt.Errorf("RmFile `%s`: sftp.Remove: %w", file.String(), err)
	}
	return nil
}

/*
type singleWriter struct {
	b  bytes.Buffer
	mu sync.Mutex
}

func (w *singleWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}
*/

// Run ...
func (conn *SSHConnection) Run(command vpath.VirtualPath, args []string, options ...RunOptions) (Process, error) {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return nil, fmt.Errorf("Run `%s`: sftp.NewClient: %w", command.String(), err)
	}
	defer client.Close()
	sess, err := conn.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Run `%s`: client.NewSession: %w", command.String(), err)
	}

	combinedOutput, writer := io.Pipe()

	sess.Stdout = writer
	sess.Stderr = writer

	process := &SSHProcess{
		combinedOutput: combinedOutput,
		cmd:            sess,
		completed:      make(chan struct{}),
	}

	cmd := command.Path
	cmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))

	if len(options) > 0 && options[0].Cwd.Path != "" {
		cmd = fmt.Sprintf("cd %s && %s", options[0].Cwd.Path, cmd)
	}
	fmt.Println("ssh running command is " + cmd)
	err = sess.Start(cmd)

	if err != nil {
		return nil, fmt.Errorf("Run `%s`: session.Start error: %w", command, err)
	}

	go func() {
		err := sess.Wait()
		if err == nil {
			process.state = 0
		} else {
			if exerr, ok := err.(*ssh.ExitError); ok {
				process.state = exerr.ExitStatus()
			} else {
				panic(err)
			}

		}

		writer.Close()
		close(process.completed)
	}()

	return process, nil
}

// SSHProcess ...
type SSHProcess struct {
	cmd *ssh.Session
	//stdout io.Reader
	//stderr io.Reader
	combinedOutput io.Reader
	completed      chan struct{}
	state          int

	//streamsCompleted *sync.WaitGroup
}

// CombinedOutput ...
func (proc *SSHProcess) CombinedOutput() io.Reader {
	/*
		proc.streamsCompleted.Add(1)
		combined, combinedWriter := io.Pipe()

		done := sync.WaitGroup{}
		done.Add(2)

		go func() {
			io.Copy(combinedWriter, proc.stdout)
			done.Done()
		}()

		go func() {
			io.Copy(combinedWriter, proc.stderr)
			done.Done()
		}()

		go func() {
			done.Wait()
			combinedWriter.Close()
			proc.streamsCompleted.Done()
		}()
	*/
	return proc.combinedOutput
}

// Kill ...
func (proc *SSHProcess) Kill() error {
	return nil
}

/*

// Stdin ...
func (proc *SSHProcess) Stdin() io.Writer {
	return nil
}

// Stdout ...
func (proc *SSHProcess) Stdout() io.Reader {
	proc.streamsCompleted.Add(1)
	processStdout, processStdoutWriter := io.Pipe()

	go func() {
		io.Copy(processStdoutWriter, proc.stdout)
		processStdoutWriter.Close()
		proc.streamsCompleted.Done()
	}()

	return processStdout
}

// Stderr ...
func (proc *SSHProcess) Stderr() io.Reader {
	proc.streamsCompleted.Add(1)
	processStderr, processStderrWriter := io.Pipe()

	go func() {
		io.Copy(processStderrWriter, proc.stderr)
		processStderrWriter.Close()
		proc.streamsCompleted.Done()
	}()

	return processStderr
}
*/

// Wait ...
func (proc *SSHProcess) Wait() (int, error) {
	<-proc.completed
	return proc.state, nil

}
