package cli

import (
	"fmt"
	"log"
	//	"gopkg.in/yaml.v2"
	//	"io/ioutil"
	//	"text/template"
	"github.com/qadium/plumber/graph"
	"github.com/qadium/plumber/shell"
	"os/exec"
	"path/filepath"
	// "golang.org/x/oauth2/google"
	// "golang.org/x/oauth2"
	// "google.golang.org/cloud"
	// "google.golang.org/cloud/container"
	// kubectl "github.com/GoogleCloudPlatform/kubernetes/pkg/kubectl/cmd"
	// cmdutil "github.com/GoogleCloudPlatform/kubernetes/pkg/kubectl/cmd/util"
	// "os"
)

func contextsToGraph(ctxs []*Context) []*graph.Node {
	nodes := make([]*graph.Node, len(ctxs))
	m := make(map[string]int) // this map maps inputs to the index of the
	// node that uses it
	// build a map to create the DAG
	for i, ctx := range ctxs {
		nodes[i] = graph.NewNode(ctx.Name)
		for _, input := range ctx.Inputs {
			m[input.Name] = i
		}
	}

	for i, ctx := range ctxs {
		for _, output := range ctx.Outputs {
			if v, ok := m[output.Name]; ok {
				nodes[i].Children = append(nodes[i].Children, nodes[v])
			}
		}
	}
	return nodes
}

func Start(pipeline string, gce string) error {
	log.Printf("==> Starting '%s' pipeline", pipeline)
	defer log.Printf("<== '%s' finished.", pipeline)

	path, err := pipelinePath(pipeline)
	if err != nil {
		return err
	}

	if gce != "" {
		// we can probably get the project name with google cloud SDK

		// step 1. re-tag local containers to gcr.io/$GCE/$pipeline-$bundlename
		// step 2. push them to gce
		// step 3. generate k8s files in pipelinePath
		// step 4. launch all the services
		err := shell.RunAndLog("kubectl", "create", "-f", fmt.Sprintf("%s/k8s", path))
		if err != nil {
			return err
		}
		// step 5: open up the firewall?
		return nil
	}
	log.Printf("NO GCE ID provided. running locally")
	// start GOOGLE experiments?
	// when start is invoked with --gce PROJECT_ID, this piece of code
	// should be run
	// client, err := google.DefaultClient(oauth2.NoContext, "https://www.googleapis.com/auth/compute")
	// if err != nil {
	// 	return err
	// }
	// cloudCtx := cloud.NewContext("kubernetes-fun", client)
	//
	// resources, err := container.Clusters(cloudCtx, "")
	// if err != nil {
	// 	return err
	// }
	// for _, op := range resources {
	// 	log.Printf("%v", op)
	// }
	//
	// loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// log.Printf("loading rules: %v", *loadingRules)
	// configOverrides := &clientcmd.ConfigOverrides{}
	// kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	// cfg, err := kubeConfig.ClientConfig()
	// if err != nil {
	// 	return err
	// }

	// well, we just shell out!
	// f := cmdutil.NewFactory(nil)
	// cmd := kubectl.NewCmdCreate(f, os.Stdout)
	// f.BindFlags(cmd.PersistentFlags())
	// cmd.Flags().Set("filename", "/Users/echu/.plumber/foo/k8s")
	// cmd.Run(cmd, []string{})

	// end GOOGLE experiments

	log.Printf(" |  Building dependency graph.")
	configs, err := filepath.Glob(fmt.Sprintf("%s/*.yml", path))
	if err != nil {
		return err
	}

	ctxs := make([]*Context, len(configs))
	for i, config := range configs {
		ctxs[i], err = parseConfig(config)
		if err != nil {
			return err
		}
	}

	// graph with diamond (test case)
	// n1 := graph.NewNode("foo")
	// n2 := graph.NewNode("bar")
	// n3 := graph.NewNode("joe")
	// n4 := graph.NewNode("bob")
	// n1.Children = append(n1.Children, n2, n3)
	// n2.Children = append(n2.Children, n4)
	// n3.Children = append(n3.Children, n4)

	g := contextsToGraph(ctxs)
	sorted, err := graph.ReverseTopoSort(g)
	if err != nil {
		return err
	}
	log.Printf("    Reverse sorted: %v", sorted)
	log.Printf("    Completed.")

	log.Printf(" |  Starting bundles...")
	managerDockerArgs := []string{"run", "-p", "9800:9800", "--rm", "plumber/manager"}
	// walk through the reverse sorted bundles and start them up
	for i := len(sorted) - 1; i >= 0; i-- {
		bundleName := sorted[i]
		log.Printf("    Starting: '%s'", bundleName)
		cmd := exec.Command("docker", "run", "-d", "-P", fmt.Sprintf("plumber/%s", bundleName))
		containerId, err := cmd.Output()
		if err != nil {
			return err
		}

		defer func() {
			log.Printf("    Stopping: '%s'", bundleName)
			cmd := exec.Command("docker", "rm", "-f", string(containerId)[0:4])
			_, err := cmd.Output()
			if err != nil {
				panic(err)
			}
			log.Printf("    Stopped.")
		}()

		log.Printf("    Started: %s", string(containerId))
		cmd = exec.Command("docker", "inspect", "--format='{{(index (index .NetworkSettings.Ports \"9800/tcp\") 0).HostPort}}'", string(containerId)[0:4])
		portNum, err := cmd.Output()
		if err != nil {
			return err
		}
		// should use docker host IP (for local deploy)
		// should use "names" for kubernetes deploy
		managerDockerArgs = append(managerDockerArgs, fmt.Sprintf("http://172.17.42.1:%s", string(portNum[:len(portNum)-1])))
	}
	log.Printf("    Args passed to 'docker': %v", managerDockerArgs)
	log.Printf("    Done.")

	log.Printf(" |  Running manager. CTRL-C to quit.")
	err = shell.RunAndLog("docker", managerDockerArgs...)
	if err != nil {
		return err
	}
	log.Printf("    Done.")
	return nil
}
