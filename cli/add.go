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
	"os"
	"fmt"
	"github.com/qadium/plumber/shell"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func addOne(ctx *Context, pipeline string, bundle string) error {
	log.Printf(" |  Adding '%s' to '%s'.", bundle, pipeline)
	defer log.Printf("    Processed '%s'.", bundle)

	path, err := ctx.GetPipeline(pipeline)
	if err != nil {
		return err
	}

	log.Printf(" |  Parsing bundle config.")
	bundleCfg, err := ParseBundleFromDir(bundle)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("    Could not find '%s' in path '%s'.", bundleConfig, bundle)
			log.Printf("    Skipping '%s'.", bundle)
			return nil
		}
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Copying `.plumber.yml` config to `%s.yml`.", bundleCfg.Name)
	config := fmt.Sprintf("%s/%s.yml", path, bundleCfg.Name)
	bytes, err := yaml.Marshal(&bundleCfg)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(config, bytes, 0644); err != nil {
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Adding `%s.yml` to version control.", bundleCfg.Name)

	if err := shell.RunAndLog("git", "-C", path, "add", config); err != nil {
		return err
	}

	message := fmt.Sprintf("Updated '%s' config.", bundleCfg.Name)
	if err := shell.RunAndLog("git", "-C", path, "commit", "-m", message, "--author", "\"Plumber Bot <plumber@qadium.com>\""); err != nil {
		return err
	}

	return nil
}

func (ctx *Context) Add(pipeline string, bundles ...string) error {
	log.Printf("==> Adding '%v' to '%s' pipeline", bundles, pipeline)
	defer log.Printf("<== Adding complete.")

	for _, bundle := range bundles {
		if err := addOne(ctx, pipeline, bundle); err != nil {
			return err
		}
	}

	return nil
}
