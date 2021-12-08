package cmd

import (
	"os"
	user2 "os/user"
	"path/filepath"
)

func RemoveNodeDir() {
	user, err := user2.Current()
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(filepath.Join(user.HomeDir, ".sifnoded"))
	if err != nil {
		panic(err)
	}
}
