package app

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/ini.v1"
)

var sampleConf = `# fuck
[global]
#log = /var/log/goKVM.log
switch = right    # move mose to top right to switch change
port = 1357
psk = Amneiht@12345
clipboard = yes    # share clibroad  default no
listen = 0.0.0.0  # listen on all interface
`

func TestConfig(t *testing.T) {
	// cfg, err := ini.Load([]byte(sampleConf))
	data, err := os.ReadFile("D:/Workspace/Go/goKVM/build/config.conf")
	if err != nil {
		fmt.Println("Read file", err)
	}
	cfg, err := ini.Load(data)
	if err != nil {
		fmt.Println(err)
	}
	gb := cfg.Section(DEFAULTSESSION)
	psk := gb.Key(PSK)
	fmt.Println("psk =", psk)
	fmt.Println("file input ", string(data))
}
