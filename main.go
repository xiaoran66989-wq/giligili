package main

import (
	"giligili/conf"
	"giligili/server"
	"net"
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

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	if err := r.RunListener(listener); err != nil {
		panic(err)
	}
}
