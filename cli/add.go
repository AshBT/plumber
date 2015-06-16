package cli

import (
	"fmt"
	"github.com/qadium/plumber/shell"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func addOne(ctx *Context, pipeline string, bundle string) error {
	log.Printf(" |  Adding '%s' to '%s'.", bundle, pipeline)
	defer log.Printf("    Added '%s'.", bundle)

	path, err := ctx.GetPipeline(pipeline)
	if err != nil {
		return err
	}

	log.Printf(" |  Parsing bundle config.")
	bundleConfig, err := ParseBundleFromDir(bundle)
	if err != nil {
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Copying `.plumber.yml` config to `%s.yml`.", bundleConfig.Name)
	config := fmt.Sprintf("%s/%s.yml", path, bundleConfig.Name)
	bytes, err := yaml.Marshal(&bundleConfig)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(config, bytes, 0644); err != nil {
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Adding `%s.yml` to version control.", bundleConfig.Name)

	if err := shell.RunAndLog("git", "-C", path, "add", config); err != nil {
		return err
	}

	message := fmt.Sprintf("Updated '%s' config.", bundleConfig.Name)
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
