package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
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
	c := exec.Command(`python -c 'from playsound import playsound; playsound("./assets/Ding.mp3")'`)
	c.Stdout = os.Stdout
	log.Println(c.Run())
}
