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
