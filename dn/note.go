package dn

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
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
// # year: %YYYY (output e.g. 2025) or %YY (25)
// # month: %MM (output e.g. 09) or %M (September)
// # day: %D (e.g. 03)
// # weekday: %W (e.g. Monday) or %WW (01)
func genName(format string) string {
	now := time.Now()
	var name string
	if format == "" {
		name = fmt.Sprintf("%v-%v-%v %v", now.Year(), int(now.Month()), now.Day(), now.Weekday().String()[:3])
	} else {
		name = formatName(format, now)
	}
	fmt.Println(name)
	return name
}

func formatName(format string, now time.Time) string {
	re := regexp.MustCompile(`%YY|%YYYY|%MM|%M|%D|%W|%WW`)
	matches := re.FindAllString(format, -1)

	for _, match := range matches {
		switch match {
		case "%YY":
			format = regexp.MustCompile(`%YY`).ReplaceAllString(format, fmt.Sprintf("%02d", now.Year()%100))
		case "%YYYY":
			format = regexp.MustCompile(`%YYYY`).ReplaceAllString(format, fmt.Sprintf("%d", now.Year()))
		case "%MM":
			format = regexp.MustCompile(`%MM`).ReplaceAllString(format, fmt.Sprintf("%02d", now.Month()))
		case "%M":
			format = regexp.MustCompile(`%M`).ReplaceAllString(format, now.Month().String())
		case "%D":
			format = regexp.MustCompile(`%D`).ReplaceAllString(format, fmt.Sprintf("%02d", now.Day()))
		case "%W":
			format = regexp.MustCompile(`%W`).ReplaceAllString(format, now.Weekday().String())
		case "%WW":
			format = regexp.MustCompile(`%WW`).ReplaceAllString(format, fmt.Sprintf("%02d", int(now.Weekday())))
		}
	}
	return format
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
	if _, err := os.Stat(confPath); err != nil {
		if os.IsNotExist(err) {
			return Config{}
		} else {
			panic(err)
		}
	}
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
