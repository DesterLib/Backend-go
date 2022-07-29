package rclone

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/desterlib/backend-go/config"
)

func Restart() {
	var rcloneBin string = "bin/rclone"
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell.exe", strings.Fields(fmt.Sprintf("powershell.exe Stop-Process -Id (Get-NetTCPConnection -LocalPort %d).OwningProcess -Force", config.ValueOf.RcloneListenPort))...)
		rcloneBin += ".exe"
	case "linux":
		cmd = exec.Command("bash", strings.Fields(fmt.Sprintf("kill $(lsof -t -i:%d)", config.ValueOf.Port))...)
	case "darwin":
		cmd = exec.Command("kill", strings.Fields(fmt.Sprintf("$(lsof -t -i:%d)", config.ValueOf.RcloneListenPort))...)
	default:
		fmt.Println("Unsupported OS:", runtime.GOOS)
		os.Exit(1)
	}
	err := cmd.Run()
	if err != nil {
		panic(err.Error())
	}
	if !exists(rcloneBin) {
		cmd = exec.Command("python", "scripts/install_rclone.py")
		cmd.Run()
		if !exists(rcloneBin) {
			panic(`Couldn't find rclone binary
Please download a suitable executable of rclone from 'rclone.org' and move it to the 'bin' folder.`)
			//
		}
	}
	cmd = exec.Command(rcloneBin, strings.Fields(fmt.Sprintf("rcd --rc-no-auth --rc-serve --rc-addr localhost:%d --config rclone.conf", config.ValueOf.RcloneListenPort))...)
	cmd.Run()
	// os.StartProcess()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func createProcess(name string, arg []string) {}
