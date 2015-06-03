package shell

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func logBuffer(prefix string, pipe io.ReadCloser) {
	in := bufio.NewScanner(pipe)

	for in.Scan() {
		log.Printf("    %s %s", prefix, in.Text())
	}
	if err := in.Err(); err != nil {
		log.Printf("    error %s", err)
	}
}

func RunAndLog(name string, args ...string) error {
	log.Printf("    Exec shell: '%s'", name+" "+strings.Join(args, " "))
	cmd := exec.Command(name, args...)

	// install signal handler to forward interrupts to subprocess
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer signal.Stop(sig)

	go func() {
		// wait for signal from CTRL-C
		<-sig
		log.Printf("    Received CTRL-C; terminating '%s' process.", name)
		cmd.Process.Signal(os.Interrupt)
	}()

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

	go logBuffer(">", stdout)
	go logBuffer("!", stderr)

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
