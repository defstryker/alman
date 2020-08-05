package main

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/aerth/playwav"
)

func screenClear() {
	goos := runtime.GOOS
	var cmd string
	switch goos {
	case "windows":
		cmd = "cls"
	default:
		cmd = "clear"
	}
	c := exec.Command(cmd)
	c.Stdout = os.Stdout
	c.Run()
}

func beep() {
	_, err := playwav.FromFile("assets/bell.wav")
	if err != nil {
		panic(err)
	}
	// c := exec.Command(``)
	// c.Stdout = os.Stdout
	// log.Fatalln(c.Run())
}
