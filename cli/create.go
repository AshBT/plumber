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
