package plumb

import (
	"os/user"
	"os"
	"fmt"
	"log"
	"github.com/qadium/plumb/shell"
)

func Create(name string) error {
	// creates a pipeline by initializing a git repo at ~/.plumb/<NAME>
	log.Printf("==> Creating '%s' pipeline", name)
	defer log.Printf("<== Creation complete.")


	log.Printf(" |  Making directory")
	usr, err := user.Current()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/.plumb/%s", usr.HomeDir, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	log.Printf("    Created pipeline directory at '%s'", path)

	log.Printf(" |  Initializing pipeline with git")
	shell.RunAndLog("git", "init", path)
	log.Printf("    Done.")
	// log.Printf("    Created '%s'", wrapper.Name())
	//
	// log.Printf(" |  Writing wrapper.")
	// log.Printf("    Done.")
	//
	// log.Printf(" |  Making temp file for Dockerfile")
	// log.Printf("    Created '%s'", dockerfile.Name())
	//
	// log.Printf(" |  Writing Dockerfile.")
	// log.Printf("    Done.")
	//
	// log.Printf(" |  Building container.")
	// log.Printf("    Container 'plumb/%s' built.", ctx.Name)
	return nil
}
