/**
 * Copyright 2015 Qadium, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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
