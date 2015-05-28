package shell

import (
	"os/exec"
	"log"
	"bufio"
	"io"
	"strings"
)

func RunAndLog(name string, args ...string) error {
	log.Printf("    Exec shell: '%s'", name + " " + strings.Join(args, " "))
	cmd := exec.Command(name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	multi := io.MultiReader(stdout, stderr)
	in := bufio.NewScanner(multi)

	for in.Scan() {
		log.Printf("    > %s", in.Text())
	}
	if err := in.Err(); err != nil {
		log.Printf("    > error: %s", err)
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
