package dn

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	s "strings"
	"time"

	"github.com/BurntSushi/toml"
)

func Enter(conf Config) {
	if conf.Path[0] == '~' {
		// replace '~' with actual user path
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		conf.Path = home + conf.Path[1:]
	}
	if conf.Path[len(conf.Path)-1] != '/' {
		conf.Path += "/"
	}

	if conf.Editor == "" {
		var ok bool
		conf.Editor, ok = os.LookupEnv("EDITOR")
		if ok == false {
			if os.PathSeparator == '\\' { // is Windows?
				conf.Editor = "notepad.exe"
			} else {
				conf.Editor = "vim"
			}
		}
	}

	filename := genName(conf.Format, conf.Extension, conf.IsPrevious)
	fullpath := conf.Path + filename
	cmd := exec.Command(conf.Editor, fullpath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func genName(format, ext string, isPrevious bool) string {
	now := time.Now()
	yyyy := fmt.Sprintf("%d", now.Year())
	yy := fmt.Sprintf("%02d", now.Year()%100)
	mm := fmt.Sprintf("%02d", now.Month())
	m := now.Month().String()
	d := fmt.Sprintf("%02d", now.Day())
	ww := fmt.Sprintf("%02d", int(now.Weekday()))
	w := now.Weekday().String()
	wlower := now.Weekday().String()[:3]

	re := regexp.MustCompile(`%YYYY|%YY|%MM|%M|%D|%WW|%W|%w`)
	matches := re.FindAllString(format, -1)

	if isPrevious {
		switch {
		case s.Contains(format, `%D`) || s.Contains(format, `%W`) || s.Contains(format, `%w`):
			wprev := now.Weekday() - 1
			ww = fmt.Sprintf("%02d", int(wprev))
			w = wprev.String()
			wlower = wprev.String()[:3]
			d = fmt.Sprintf("%02d", now.Day()-1)
		case s.Contains(format, `%M`):
			mprev := now.Month() - 1
			mm = fmt.Sprintf("%02d", mprev)
			m = mprev.String()
		case s.Contains(format, `%YY`):
			yyyy = fmt.Sprintf("%d", now.Year()-1)
			yy = fmt.Sprintf("%02d", now.Year()%100-1)
		}
	}

	for _, match := range matches {
		switch match {
		case "%YYYY":
			format = regexp.MustCompile(`%YYYY`).ReplaceAllString(format, yyyy)
		case "%YY":
			format = regexp.MustCompile(`%YY`).ReplaceAllString(format, yy)
		case "%MM":
			format = regexp.MustCompile(`%MM`).ReplaceAllString(format, mm)
		case "%M":
			format = regexp.MustCompile(`%M`).ReplaceAllString(format, m)
		case "%D":
			format = regexp.MustCompile(`%D`).ReplaceAllString(format, d)
		case "%WW":
			format = regexp.MustCompile(`%WW`).ReplaceAllString(format, ww)
		case "%W":
			format = regexp.MustCompile(`%W`).ReplaceAllString(format, w)
		case "%w":
			format = regexp.MustCompile(`%w`).ReplaceAllString(format, wlower)
		}
	}
	return format + "." + ext
}

type Config struct {
	Path       string
	Editor     string
	Format     string
	Extension  string
	IsPrevious bool
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
