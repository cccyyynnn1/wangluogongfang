package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func main() {

	cmd := exec.Command("D:\\Anaconda3\\python.exe", "C:\\Users\\Admin\\Desktop\\gogogogo\\code\\main.py")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Output:", out.String())
	}
}
