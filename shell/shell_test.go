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
package shell_test

import (
	"github.com/qadium/plumber/shell"
	"io/ioutil"
	"log"
	"syscall"
	"testing"
	"time"
)

func TestRunAndLog(t *testing.T) {
	err := shell.RunAndLog("true")
	if err != nil {
		t.Error(err)
	}
}

func TestRunAndLogFails(t *testing.T) {
	err := shell.RunAndLog("-notlikely-to-be-a*-cmd")
	if err == nil {
		t.Error("Expected an error but never got one!")
	}
}

func TestInterrupt(t *testing.T) {
	// set the interrupt handler to go off after 50 milliseconds
	go func() {
		time.Sleep(50 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	err := shell.RunAndLog("/bin/bash", "-c", "while true; do true; done")
	if err == nil || err.Error() != "signal: interrupt" {
		t.Error("Should've received a SIGINT")
	}
}

func BenchmarkRunAndLog(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		shell.RunAndLog("echo", "true")
	}
}
