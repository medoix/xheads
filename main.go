package main

import (
	"fmt"
	"os/exec"
	"strings"
)


func currentOutputs() []string {
	return scanConnectedOutputs()
}

func makeXrandrScanCommand() string {
	return fmt.Sprintf("xrandr | grep ' connected' | awk '{print $1}'")
}

func scanConnectedOutputs() (outputs []string) {
        cmd := makeXrandrScanCommand()
        out, err := exec.Command("sh", "-c", cmd).Output()
        if err != nil {
                fmt.Println("Cannot read connected displays:", err)
                return
        }
        lines := strings.Trim(string(out), "\n")
        outputs = strings.Split(lines, "\n")
        return
}

func main() {
	outputs := currentOutputs()
	fmt.Println("outputs: ", outputs)
}
