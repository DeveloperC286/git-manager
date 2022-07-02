package main

import (
	"encoding/csv"
	"flag"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Get where to read in the config from.
	config := flag.String("config", "~/.git-repositories-managing", "The data file containing all the Git repositories to manage.")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})

	// Getting user home directory so we can find replace ~ as Go does not handle it.
	usr, err := user.Current()
	home := usr.HomeDir

	if err != nil {
		log.Fatal(err)
	}

	// Read in the config.
	file, err := os.Open(toPath(home, *config))

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Try to read the config as CSV data.
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	// Used to wait upon all the Go routines for line of data.
	var channels []chan struct{}

	for _, line := range data {
		// TODO error handling if not remote,local format.
		if len(line) == 2 {
			// Channel to return the Go routine for this line has finished.
			c := make(chan struct{})
			channels = append(channels, c)

			// Go routine per line as the Git network operations are the bottleneck.
			go func(remote string, local string, c chan struct{}) {
				if !exists(local) {
					// Does not exist locally so just clone from remote.
					cmd := exec.Command("mkdir", "-p", local)
					err := cmd.Run()

					if err != nil {
						log.Error(err)
						c <- struct{}{}
						return
					}

					cmd = exec.Command("git", "clone", remote, local)
					err = cmd.Run()

					if err != nil {
						log.Error(err)
						c <- struct{}{}
						return
					}

					log.Info("Cloned " + remote + " to " + local)
				}

				c <- struct{}{}
			}(line[0], toPath(home, line[1]), c)
		}
	}

	// Wait on the Go routine for each line of data to finish.
	for _, waitingOn := range channels {
		<-waitingOn
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func toPath(home string, path string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}

	return path
}
