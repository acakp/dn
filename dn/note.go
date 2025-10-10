package dn

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Enter(filepath string) {
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

	editor, ok := os.LookupEnv("EDITOR")
	if ok == false {
		fmt.Println("[warning]: $EDITOR environment variable is empty, will use vim instead...")
		editor = "vim"
	}

	filename := genName()
	fullpath := filepath + filename
	cmd := exec.Command(editor, fullpath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func genName() string {
	now := time.Now()
	name := fmt.Sprintf("%v-%v-%v %v", now.Year(), int(now.Month()), now.Day(), now.Weekday().String()[:3])
	fmt.Println(name)
	return name
}
