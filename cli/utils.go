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
	"net"
	"net/url"
	"os"
	"os/user"
	"strings"
)

// Useful information for the plumber CLI tool goes here
type Context struct {
	PipeDir       string // the directory to store plumber pipelines
	KubeSubdir    string // the suffix to use to store kubernetes files
	GitCommit     string // the current git commit
	Version       string // the current version
	ManagerImage  string // the desired image name for bootstrapping
	BootstrapDir  string // the directory to use for bootstrapping
	ImageRepo     string // the prefix to use for images
	DockerCmd     string // the command for running `docker` in a shell
	DockerIface   string // the docker network interface name
	DockerHostEnv string // the DOCKER_HOST environment variable
	GcloudCmd     string // the command for running `gcloud` in a shell
	KubectlCmd    string // the command for running `kubectl` in a shell
}

const (
	plumberDir = ".plumber"
	bootstrapDir = ".plumber-bootstrap"
	k8sDir = "k8s"
)

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
		"manager",
		fmt.Sprintf("%s/%s", usr.HomeDir, bootstrapDir),
		"plumber",
		"docker",
		"docker0",
		"DOCKER_HOST",
		"gcloud",
		"kubectl",
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
	k8s := fmt.Sprintf("%s/%s", path, d.KubeSubdir)
	return k8s
}

// Get the manager's image name
func (d *Context) GetManagerImage() string {
	return d.GetImage(d.ManagerImage)
}

// Get an image name
func (d *Context) GetImage(name string) string {
	if d.ImageRepo == "" {
		return name
	} else {
		return fmt.Sprintf("%s/%s", d.ImageRepo, name)
	}
}

// Gets the IP address used for local containers. This uses the
// DOCKER_HOST environment variable; if not available, we'll look for
// the ip address associated with the `docker0` device. If that, too,
// is missing, we will give an error.
func (d *Context) GetDockerHost() (string, error) {
	hostIp := os.Getenv(d.DockerHostEnv)
	if hostIp == "" {
		// find interface named d.DockerIface (usually "docker0")
		iface, err := net.InterfaceByName(d.DockerIface)
		if err != nil {
			return "", err
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return "", err
			}
			ipv4 := ip.To4()
			if ipv4 != nil {
				return ipv4.String(), nil
			}
		}
		return "", fmt.Errorf("Could not get IPv4 gateway for device '%s'", d.DockerIface)
	} else {
		// docker host is usually in the form of [http:// | unix://]IP:PORT
		hostUrl, err := url.Parse(hostIp)
		if err != nil {
			return "", err
		}
		hostIp = strings.Split(hostUrl.Host, ":")[0]
		return hostIp, nil
	}
}

// Replaces underscores with "-" and lowercases the string
func EscapeUnderscore(a string) string {
	return strings.ToLower(strings.Replace(a, "_", "-", -1))
}

func versionString() string {
	versionString := version
	if versionPrerelease != "" {
		versionString += "-" + versionPrerelease
	}
	return versionString
}
