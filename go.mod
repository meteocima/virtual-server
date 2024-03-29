module github.com/meteocima/virtual-server

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mikkeloscar/sshconfig v0.1.0
	github.com/pkg/sftp v1.12.0
	github.com/stretchr/testify v1.6.1
	github.com/tevino/abool v1.2.0
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

// replace github.com/mikkeloscar/sshconfig => ../sshconfig
