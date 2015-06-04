package cli

import (
	"github.com/qadium/plumber/shell"
	"log"
	"os"
)

func Create(name string) error {
	// creates a pipeline by initializing a git repo at ~/.plumb/<NAME>
	log.Printf("==> Creating '%s' pipeline", name)
	defer log.Printf("<== Creation complete.")

	log.Printf(" |  Making directory")
	// note that we use PipelinePath instead of GetPipeline here; this
	// is because we only need the path to create it
	path, err := PipelinePath(name)
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
