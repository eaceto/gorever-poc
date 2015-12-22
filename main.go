package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/eaceto/gorever-poc/gorever"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
	"time"
)

var (
	mode = kingpin.Flag("mode", "Execution mode").Default("updater").String()
)

func main() {
	kingpin.Parse()
	if *mode == "" {
		*mode = "updater"
	}

	logrus.SetLevel(logrus.DebugLevel)
	log := logrus.WithFields(logrus.Fields{
		"mode": *mode,
	})

	log.Debugf(">> Starting as %s - pid: %d", *mode, os.Getpid())

	if *mode == "updater" {
		updateChannel := make(chan bool, 16)
		updater, err := gorever.NewUpdater(updateChannel)

		var configProcess *os.Process = nil
		for {
			if configProcess != nil {
				err = configProcess.Kill()
				if err != nil {
					log.Debugf("Could not kill process %d", configProcess.Pid)
				}
			}

			configProcess, err = launchConfigurator(log)
			if err != nil {
				panic(err)
			} else if configProcess != nil {
				log.Debugf("Configurator at pid: %d", configProcess.Pid)
			}

		UpdateLoop:
			for {
				switch {
				case <-updateChannel:
					if updater.HasNewVersion() {
						log.Debugf("A new version is available!")
						err = updater.Update()
						if err != nil {
							panic(err)
						}
						log.Debugf("Application updated.")
						break UpdateLoop
					}
					break
				}
			}
		}
	} else if *mode == "configurator" {
		for {
			time.Sleep(2 * time.Second)
			log.Debugf("Configurator running with PID %d", os.Getpid())
		}
	}

	os.Exit(0)
}

func launchConfigurator(log *logrus.Entry) (*os.Process, error) {
	configArgs := make([]string, 1)
	configArgs = append(configArgs, "--mode")
	configArgs = append(configArgs, "configurator")
	mode := false
	for _, v := range os.Args[1:] {
		if strings.Contains(v, "--mode") {
			mode = true
		} else if mode {
			mode = false
		} else {
			configArgs = append(configArgs, v)
		}
	}

	log.Debugf("Will run %s with %s", os.Args[0], configArgs)

	cmdToRun := os.Args[0]
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	process, err := os.StartProcess(cmdToRun, configArgs, procAttr)
	if err != nil {
		log.Debugf("ERROR Unable to run %s: %s", cmdToRun, err.Error())
		return nil, err
	}

	return process, nil
}
