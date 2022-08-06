package rclone

import (
	"bufio"
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
	fmt.Println("searching os")
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell.exe", strings.Fields(fmt.Sprintf("powershell.exe Stop-Process -Id (Get-NetTCPConnection -LocalPort %d).OwningProcess -Force", config.ValueOf.RcloneListenPort))...)
		rcloneBin += ".exe"
	case "darwin", "linux":
		cmd = exec.Command("kill", strings.Fields(fmt.Sprintf("$(lsof -t -i:%d)", config.ValueOf.RcloneListenPort))...)
	default:
		fmt.Println("Unsupported OS:", runtime.GOOS)
		os.Exit(1)
	}
	cmd.Run()
	if !exists(rcloneBin) {
		cmd = exec.Command("python", "scripts/install_rclone.py")
		cmd.Run()
		if !exists(rcloneBin) {
			panic(`Couldn't find rclone binary
Please download a suitable executable of rclone from 'rclone.org' and move it to the 'bin' folder.`)
			//
		}
	}
	fmt.Println("running rclone")
	startRclone(rcloneBin)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func startRclone(rcloneBin string) {
	cmd := exec.Command(rcloneBin, strings.Fields(fmt.Sprintf("rcd --rc-no-auth --rc-serve --rc-addr localhost:%d --config rclone.conf", config.ValueOf.RcloneListenPort))...)
	cmd.Stderr = os.Stdout
	w, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(w)
	scanner.Split(bufio.ScanWords)
	for scanner != nil && scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	err := cmd.Wait()
	if err != nil {
		fmt.Println("running rclone with osperm")
		if errors.Is(err, os.ErrPermission) {
			exec.Command("chmod", "+x", rcloneBin).Run()
			cmd.Start()
		}
	}
}
