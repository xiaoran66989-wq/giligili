package main

import (
	"giligili/conf"
	"giligili/server"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	preparePackagedWorkingDirectory()
	conf.Init()
	r := server.NewRouter()

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":3000"
		if runtime.GOOS == "windows" {
			addr = "127.0.0.1:3000"
		}
	}

	if runtime.GOOS == "windows" && os.Getenv("AUTO_OPEN_BROWSER") != "0" {
		go func() {
			time.Sleep(time.Second)
			_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:3000/").Start()
		}()
	}

	if err := r.Run(addr); err != nil {
		panic(err)
	}
}

func preparePackagedWorkingDirectory() {
	if runtime.GOOS != "windows" {
		return
	}
	executable, err := os.Executable()
	if err != nil {
		return
	}
	appDir := filepath.Dir(executable)
	if _, err := os.Stat(filepath.Join(appDir, "web", "index.html")); err == nil {
		_ = os.Chdir(appDir)
	}
}
