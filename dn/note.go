package dn

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
)

func Enter(filepath, editor, format, ext string) {
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
			if os.PathSeparator == '\\' { // is Windows?
				editor = "notepad.exe"
			} else {
				editor = "vim"
			}
		}
	}

	filename := genName(format, ext)
	fullpath := filepath + filename
	cmd := exec.Command(editor, fullpath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func genName(format string, ext string) string {
	now := time.Now()
	var name string
	if format == "" {
		name = fmt.Sprintf("%v-%v-%v %v.%s", now.Year(), int(now.Month()), now.Day(), now.Weekday().String()[:3], ext)
	} else {
		name = formatName(format, ext, now)
	}
	return name
}

func formatName(format, ext string, now time.Time) string {
	re := regexp.MustCompile(`%YYYY|%YY|%MM|%M|%D|%WW|%W|%w`)
	matches := re.FindAllString(format, -1)

	for _, match := range matches {
		switch match {
		case "%YYYY":
			format = regexp.MustCompile(`%YYYY`).ReplaceAllString(format, fmt.Sprintf("%d", now.Year()))
		case "%YY":
			format = regexp.MustCompile(`%YY`).ReplaceAllString(format, fmt.Sprintf("%02d", now.Year()%100))
		case "%MM":
			format = regexp.MustCompile(`%MM`).ReplaceAllString(format, fmt.Sprintf("%02d", now.Month()))
		case "%M":
			format = regexp.MustCompile(`%M`).ReplaceAllString(format, now.Month().String())
		case "%D":
			format = regexp.MustCompile(`%D`).ReplaceAllString(format, fmt.Sprintf("%02d", now.Day()))
		case "%WW":
			format = regexp.MustCompile(`%WW`).ReplaceAllString(format, fmt.Sprintf("%02d", int(now.Weekday())))
		case "%W":
			format = regexp.MustCompile(`%W`).ReplaceAllString(format, now.Weekday().String())
		case "%w":
			format = regexp.MustCompile(`%w`).ReplaceAllString(format, now.Weekday().String()[:3])
		}
	}
	return format + "." + ext
}

type Config struct {
	Path      string
	Editor    string
	Format    string
	Extension string
}

func ReadConf() Config {
	confPath, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	pathSep := string(os.PathSeparator)
	confPath += fmt.Sprintf("%vdn%vconfig.toml", pathSep, pathSep)
	if _, err := os.Stat(confPath); err != nil {
		if os.IsNotExist(err) {
			// fmt.Printf("[info]: config file (%v) not found\n", confPath)
			return Config{}
		} else {
			panic(err)
		}
	}
	confRaw, err := os.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	var conf Config
	_, err = toml.Decode(string(confRaw), &conf)
	if err != nil {
		panic(err)
	}
	return conf
}
