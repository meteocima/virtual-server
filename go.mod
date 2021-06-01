module github.com/meteocima/virtual-server

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bramvdbogaerde/go-scp v0.0.0-20210527193300-acf430e39785
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mikkeloscar/sshconfig v0.1.0
	github.com/pkg/sftp v1.12.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

// replace github.com/mikkeloscar/sshconfig => ../sshconfig
