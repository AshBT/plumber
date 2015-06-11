package cli

import (
	"fmt"
	"os"
	"os/user"
)

// Useful information for the plumber CLI tool goes here
type Context struct {
	PipeDir      string // the directory to store plumber pipelines
	KubeSuffix   string // the suffix to use to store kubernetes files
	GitCommit    string // the current git commit
	Version      string // the current version
	Image        string // the desired image name
	BootstrapDir string // the directory to use for bootstrapping
}

const plumberDir = ".plumber"

const bootstrapDir = ".plumber-bootstrap"

const k8sDir = "k8s"

// The default context stores all plumber pipelines in the user's
// home directory at ~/.plumber; all kubernetes files are stored at
// ~/.plumber/$PIPELINE/k8s
//
// It also includes some basic versioning information
func NewDefaultContext() (*Context, error) {
	// use the current user to store the plumb data
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	d := &Context{
		fmt.Sprintf("%s/%s", usr.HomeDir, plumberDir),
		k8sDir,
		GitCommit,
		versionString(),
		"plumber/manager",
		fmt.Sprintf("%s/%s", usr.HomeDir, bootstrapDir),
	}
	return d, nil
}

// Given the `name` of a pipeline, return the path where we should store
// information about it.
func (d *Context) PipelinePath(name string) string {
	return fmt.Sprintf("%s/%s", d.PipeDir, name)
}

// Get a pipeline; this differs from PipelinePath in that it also checks
// if the file / path exists.
func (d *Context) GetPipeline(name string) (string, error) {
	path := d.PipelinePath(name)

	// make sure file exists and we have permissions, etc.
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	return path, nil
}

// Given the `name` of a pipeline, return the path where we should store
// kubernetes YAML files for pods, services, and replication controllers
func (d *Context) KubernetesPath(name string) string {
	path := d.PipelinePath(name)
	k8s := fmt.Sprintf("%s/%s", path, d.KubeSuffix)
	return k8s
}

func versionString() string {
	versionString := version
	if versionPrerelease != "" {
		versionString += "-" + versionPrerelease
	}
	return versionString
}
