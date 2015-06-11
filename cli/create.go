package cli

import (
	"errors"
	"github.com/qadium/plumber/shell"
	"log"
	"os"
)

func (ctx *Context) Create(name string) error {
	// creates a pipeline by initializing a git repo at ~/.plumb/<NAME>
	log.Printf("==> Creating '%s' pipeline", name)
	defer log.Printf("<== Creation complete.")

	if name == "" {
		return errors.New("Cannot create a pipeline with no name.")
	}

	log.Printf(" |  Making directory")
	// note that we use PipelinePath instead of GetPipeline here; this
	// is because we only need the path to create it
	path := ctx.PipelinePath(name)

	// if the path already exists, give an error
	if _, err := os.Stat(path); err == nil {
		return errors.New("Pipeline already exists.")
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	log.Printf("    Created pipeline directory at '%s'", path)

	log.Printf(" |  Initializing pipeline with git")
	if err := shell.RunAndLog("git", "init", path); err != nil {
		return err
	}
	log.Printf("    Done.")

	return nil
}
