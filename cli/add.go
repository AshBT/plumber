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
	"fmt"
	"github.com/qadium/plumber/shell"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func addOne(ctx *Context, pipelinePath string, bundle string) error {
	log.Printf(" |  Parsing bundle config.")
	bundleCfg, err := ParseBundleFromDir(bundle)
	if err != nil {
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Copying `.plumber.yml` config to `%s.yml`.", bundleCfg.Name)
	config := fmt.Sprintf("%s/%s.yml", pipelinePath, bundleCfg.Name)
	bytes, err := yaml.Marshal(&bundleCfg)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(config, bytes, 0644); err != nil {
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Adding `%s.yml` to version control.", bundleCfg.Name)

	if err := shell.RunAndLog("git", "-C", pipelinePath, "add", config); err != nil {
		return err
	}

	message := fmt.Sprintf("Updated '%s' config.", bundleCfg.Name)
	if err := shell.RunAndLog("git", "-C", pipelinePath, "commit", "-m", message, "--author", "\"Plumber Bot <plumber@qadium.com>\""); err != nil {
		return err
	}

	return nil
}

func (ctx *Context) Add(pipeline string, bundles ...string) error {
	log.Printf("==> Adding '%v' to '%s' pipeline", bundles, pipeline)
	defer log.Printf("<== Adding complete.")

	path, err := ctx.GetPipeline(pipeline)
	if err != nil {
		return err
	}

	skipped := []string{}
	for _, bundle := range bundles {
		log.Printf(" |  Adding '%s' to '%s'.", bundle, pipeline)
		if err := addOne(ctx, path, bundle); err != nil {
			if os.IsNotExist(err) {
				log.Printf("    Could not find '%s' in path '%s'.", bundleConfig, bundle)
				log.Printf("    Skipped '%s'.", bundle)
				skipped = append(skipped, bundle)
			} else {
				return err
			}
		} else {
			log.Printf("    Added '%s'.", bundle)
		}
	}

	if len(skipped) > 0 {
		log.Printf(" *   Skipped '%v'", skipped)
	}
	return nil
}
