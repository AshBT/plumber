package cli

import (
	"fmt"
	"github.com/qadium/plumb/shell"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func addOne(pipeline string, bundle string) error {
	log.Printf(" |  Adding '%s' to '%s'.", bundle, pipeline)
	defer log.Printf("    Added '%s'.", bundle)

	path, err := pipelinePath(pipeline)
	if err != nil {
		return err
	}

	log.Printf(" |  Parsing bundle config.")
	ctx, err := parseConfigFromDir(bundle)
	log.Printf("    Done.")

	log.Printf(" |  Copying `.plumb.yml` config to `%s.yml`.", ctx.Name)
	config := fmt.Sprintf("%s/%s.yml", path, ctx.Name)
	bytes, err := yaml.Marshal(&ctx)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(config, bytes, 0644); err != nil {
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Adding `%s.yml` to version control.", ctx.Name)

	if err := shell.RunAndLog("git", "-C", path, "add", config); err != nil {
		return err
	}

	message := fmt.Sprintf("Updated '%s' config.", ctx.Name)
	if err := shell.RunAndLog("git", "-C", path, "commit", "-m", message, "--author", "\"Plumb Bot <plumb@qadium.com>\""); err != nil {
		return err
	}

	return nil
}

func Add(pipeline string, bundles ...string) error {
	log.Printf("==> Adding '%v' to '%s' pipeline", bundles, pipeline)
	defer log.Printf("<== Adding complete.")

	for _, bundle := range bundles {
		addOne(pipeline, bundle)
	}

	return nil
}
