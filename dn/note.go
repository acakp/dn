package dn

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"os/exec"
	"time"
)

func Enter(filepath, editor, format string) {
	if filepath[0] == '~' {
		// replace '~' with actual user path
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		filepath = home + filepath[1:]
	}
	if filepath[len(filepath)-1] != '/' {
		filepath += "/"
	}

	if editor == "" {
		var ok bool
		editor, ok = os.LookupEnv("EDITOR")
		if ok == false {
			fmt.Println("[warning]: $EDITOR environment variable is empty, will use vim instead...")
			editor = "vim"
		}
	}

	filename := genName(format)
	fullpath := filepath + filename
	cmd := exec.Command(editor, fullpath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

// # available variants:
// # year: $YYYY or $YY
// # month: $MM or $M
// # day: $D
// # weekday: $W or $WW
func genName(format string) string {
	now := time.Now()
	var name string
	if format == "" {
		name = fmt.Sprintf("%v-%v-%v %v", now.Year(), int(now.Month()), now.Day(), now.Weekday().String()[:3])
	} else {
	}
	fmt.Println(name)
	return name
}

type Config struct {
	Path   string
	Editor string
	Format string
}

func ReadConf() Config {
	confPath, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	confPath += "/dn/config.toml"
	fmt.Println(confPath)
	confRaw, err := os.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	var conf Config
	_, err = toml.Decode(string(confRaw), &conf)
	if err != nil {
		panic(err)
	}
	fmt.Println("path:", conf.Path)
	fmt.Println("editor:", conf.Editor)
	fmt.Println("format:", conf.Format)
	return conf
}
