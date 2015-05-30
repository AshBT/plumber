package plumb

import (
	"os/user"
	"os"
	"fmt"
	"log"
	"github.com/qadium/plumb/shell"
)

// (TODO) break out into separate file?
func pipelinePath(name string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/.plumb/%s", usr.HomeDir, name)
	return path, nil
}

func Create(name string) error {
	// creates a pipeline by initializing a git repo at ~/.plumb/<NAME>
	log.Printf("==> Creating '%s' pipeline", name)
	defer log.Printf("<== Creation complete.")


	log.Printf(" |  Making directory")
	path, err := pipelinePath(name)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	log.Printf("    Created pipeline directory at '%s'", path)

	log.Printf(" |  Initializing pipeline with git")
	shell.RunAndLog("git", "init", path)
	log.Printf("    Done.")

	return nil
}
