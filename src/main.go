package main

import (
	"bytes"
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
	// Disable the checking of number of fields per record.
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	// Used to wait upon all the Go routines for line of data.
	var channels []chan struct{}

	for lineNumber, line := range data {
		if len(line) == 2 {
			// Channel to return the Go routine for this line has finished.
			statusChannel := make(chan struct{})
			channels = append(channels, statusChannel)

			// Go routine per line as the Git network operations are the bottleneck.
			go func(remote string, local string, statusChannel chan struct{}) {
				if !exists(local) {
					// Does not exist locally so just clone from remote.
					cmd := exec.Command("git", "clone", remote, local)
					// Setup reading the command output
					stderr := new(bytes.Buffer)
					cmd.Stderr = stderr

					err = cmd.Run()

					if err != nil {
						log.WithFields(log.Fields{
							"remote": remote,
							"local":  local,
						}).WithError(err).Error(stderr.String())
						statusChannel <- struct{}{}
						return
					}

					log.WithFields(log.Fields{
						"remote": remote,
						"local":  local,
					}).Info("Successfully cloned to local location from remote repository.")
				} else {
					// Exists locally so just pull from origin/HEAD and rebase.
					cmd := exec.Command("git", "name-rev", "--name-only", "origin/HEAD")
					// Execute inside the local Git repo.
					cmd.Dir = local
					// Setup reading the command output
					stdout := new(bytes.Buffer)
					stderr := new(bytes.Buffer)
					cmd.Stdout = stdout
					cmd.Stderr = stderr

					err = cmd.Run()

					if err != nil {
						log.WithFields(log.Fields{
							"remote": remote,
							"local":  local,
						}).WithError(err).Error(stderr.String())
						statusChannel <- struct{}{}
						return
					}

					head := strings.TrimSpace(stdout.String())

					// Get the branch currently checked out.
					cmd = exec.Command("git", "branch", "--show-current")
					// Execute inside the local Git repo.
					cmd.Dir = local
					// Setup reading the command output
					stdout = new(bytes.Buffer)
					stderr = new(bytes.Buffer)
					cmd.Stdout = stdout
					cmd.Stderr = stderr

					err = cmd.Run()

					if err != nil {
						log.WithFields(log.Fields{
							"remote": remote,
							"local":  local,
						}).WithError(err).Error(stderr.String())
						statusChannel <- struct{}{}
						return
					}

					current := strings.TrimSpace(stdout.String())

					if head == current {
						cmd = exec.Command("git", "pull", "--rebase", "--autostash")
						// Execute inside the local Git repo.
						cmd.Dir = local
						// Setup reading the command output
						stderr = new(bytes.Buffer)
						cmd.Stderr = stderr

						err = cmd.Run()

						if err != nil {
							log.WithFields(log.Fields{
								"remote": remote,
								"local":  local,
							}).WithError(err).Error(stderr.String())
							statusChannel <- struct{}{}
							return
						}

						log.WithFields(log.Fields{
							"remote": remote,
							"local":  local,
						}).Info("Updated local head branch from remote head branch.")
					}
				}

				statusChannel <- struct{}{}
			}(line[0], toPath(home, line[1]), statusChannel)
		} else {
			log.Warnf("Do not know how to parse the line number %d with the content %s.", lineNumber+1, printArray(line))
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

func printArray(array []string) string {
	var buffer strings.Builder

	buffer.WriteString("[")

	for index, item := range array {
		if index != 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString("\"" + item + "\"")
	}

	buffer.WriteString("]")

	return buffer.String()
}
