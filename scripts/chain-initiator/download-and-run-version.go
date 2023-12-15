// nolint: nakedret
package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func isURL(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

func downloadAndRunVersion(binaryPathOrURL string, skipDownload bool) (path string, version string, err error) {
	if !isURL(binaryPathOrURL) {
		// If the input is a local path
		path = binaryPathOrURL

		// Check if the path exists
		if _, err = os.Stat(path); os.IsNotExist(err) {
			err = errors.New(fmt.Sprintf("binary file does not exist at the specified path: %v", path))
			return
		}

		// Run the command 'binary version'
		cmd := exec.Command(path, "version")
		var versionOutput []byte
		versionOutput, err = cmd.CombinedOutput()
		if err != nil {
			return
		}
		version = strings.TrimSpace(string(versionOutput))

		return
	}

	if skipDownload {
		// Extract version from the URL
		re := regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?`)
		versionMatches := re.FindStringSubmatch(binaryPathOrURL)
		if len(versionMatches) == 0 {
			err = errors.New("no version found in URL")
			return
		}
		version = versionMatches[0]

		// Set the binary path based on the version
		path = "/tmp/sifnoded-" + version

		// Check if the path exists
		if _, err = os.Stat(path); os.IsNotExist(err) {
			err = errors.New(fmt.Sprintf("binary file does not exist at the specified path: %v", path))
		}

		return
	}

	// Download the binary
	resp, err := http.Get(binaryPathOrURL) // nolint: gosec
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Create a temporary file
	tmpFile, err := ioutil.TempFile("", "binary-*")
	if err != nil {
		return
	}
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath) // Clean up

	// Write the downloaded content to the file
	_, err = io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	if err != nil {
		return
	}

	// Make the file executable
	err = os.Chmod(tmpFilePath, 0755)
	if err != nil {
		return
	}

	// Run the command 'binary version'
	cmd := exec.Command(tmpFilePath, "version")
	versionOutput, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	version = strings.TrimSpace(string(versionOutput))

	// Rename the temporary file
	newFilePath := "/tmp/sifnoded-" + version
	err = os.Rename(tmpFilePath, newFilePath)
	if err != nil {
		return
	}
	path = newFilePath

	return
}
