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
	name        string

	Host     string
	Port     int
	User     string
	KeyPath  string
	hostName string
	config   *ssh.ClientConfig
	client   *ssh.Client
}

// Name ...
func (conn *SSHConnection) Name() string {
	return conn.name
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
				conn.Name(),
				conn.Host,
				conn.Host,
			)
		}
		fmt.Printf(
			"VPE: successfully connected to server `%s` using hostname %s. Subsequents requests will use hostname %s\n",
			conn.Name(),
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
				conn.Name(),
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

// OpenAppendWriter ...
func (conn *SSHConnection) OpenAppendWriter(file vpath.VirtualPath) (io.WriteCloser, error) {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return nil, err
	}
	writer, err := client.OpenFile(file.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY)
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
		return nil, fmt.Errorf("cannot create new sftp client: %w", err)
	}
	defer client.Close()

	return client.Stat(path.Path)
}

// Glob ...
func (conn *SSHConnection) Glob(pattern vpath.VirtualPath) (vpath.VirtualPathList, error) {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	files, err := client.Glob(pattern.Path)
	if err != nil {
		return nil, err
	}
	result := make(vpath.VirtualPathList, len(files))
	for idx, file := range files {
		result[idx] = vpath.New(pattern.Host, file)
	}
	return result, nil
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
func (conn *SSHConnection) Run(command vpath.VirtualPath, args []string, options RunOptions) (Process, error) {
	client, err := sftp.NewClient(conn.client)
	if err != nil {
		return nil, fmt.Errorf("Run `%s`: sftp.NewClient: %w", command.String(), err)
	}
	defer client.Close()

	cmd, err := conn.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Run `%s`: client.NewSession: %w", command.String(), err)
	}

	if options.Stderr == nil {
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = options.Stderr
	}

	if options.Stdout == nil {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = options.Stdout
	}

	if options.Stdin == nil {
		cmd.Stdin = os.Stdin
	} else {
		cmd.Stdin = options.Stdin
	}

	process := &SSHProcess{
		cmd:       cmd,
		completed: make(chan struct{}),
	}

	if options.OutFromLog != nil {
		go copyLines(process, cmd.Stdout, *options.OutFromLog)
	}

	if options.ErrFromLog != nil {
		go copyLines(process, cmd.Stderr, *options.ErrFromLog)
	}

	cmdStr := command.Path
	cmdStr = fmt.Sprintf("%s %s", cmdStr, strings.Join(args, " "))

	if options.Cwd.Path != "" {
		cmdStr = fmt.Sprintf("cd %s && %s", options.Cwd.Path, cmdStr)
	}

	err = cmd.Start(cmdStr)

	if err != nil {
		return nil, fmt.Errorf("Run `%s`: session.Start error: %w", command, err)
	}

	go func() {
		err := cmd.Wait()

		if err == nil {
			process.state = 0
		} else {
			if exerr, ok := err.(*ssh.ExitError); ok {
				process.state = exerr.ExitStatus()
			} else {
				panic(err)
			}
		}
		close(process.completed)
	}()

	return process, nil
}

// SSHProcess ...
type SSHProcess struct {
	cmd       *ssh.Session
	completed chan struct{}
	state     int
}

// Kill ...
func (proc *SSHProcess) Kill() error {
	return nil
}

// Wait ...
func (proc *SSHProcess) Wait() (int, error) {
	<-proc.completed
	return proc.state, nil

}
