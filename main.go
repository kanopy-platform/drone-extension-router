package main

import (
	"github.com/kanopy-platform/drone-extension-router/internal/cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		log.Fatalln(err)
	}
}
