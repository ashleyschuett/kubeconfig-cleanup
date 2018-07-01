package main

import (
	"github.com/ashleyschuett/kubeconfig-cleanup/pkg/config"
)

func main() {
	m := config.NewManager()
	m.Run()
}
