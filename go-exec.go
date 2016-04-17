package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

var env []string

func init() {
	cp, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "GOPATH=") {
			env = append(env, v+":"+cp)
		} else {
			env = append(env, v)
		}
	}
}

// GoExec execute go commang with GOPATH=$GOPATH:`pwd`
func GoExec(params ...string) {
	var err error
	fmt.Printf("run: go %s", strings.Join(params, " "))
	cmd := exec.Command("go", params...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	if err = cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
	}
	fmt.Println()
}
