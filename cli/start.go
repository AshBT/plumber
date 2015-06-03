package cli

import (
	"fmt"
	"log"
	//	"gopkg.in/yaml.v2"
	//	"io/ioutil"
	//	"text/template"
	"github.com/qadium/plumb/graph"
	"github.com/qadium/plumb/shell"
	"os/exec"
	"path/filepath"
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

func Start(pipeline string) error {
	log.Printf("==> Starting '%s' pipeline", pipeline)
	defer log.Printf("<== '%s' finished.", pipeline)

	log.Printf(" |  Building dependency graph.")
	path, err := pipelinePath(pipeline)
	if err != nil {
		return err
	}

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
	managerDockerArgs := []string{"run", "-p", "9800:9800", "--rm", "plumb/manager"}
	// walk through the reverse sorted bundles and start them up
	for i := len(sorted) - 1; i >= 0; i-- {
		bundleName := sorted[i]
		log.Printf("    Starting: '%s'", bundleName)
		cmd := exec.Command("docker", "run", "-d", "-P", fmt.Sprintf("plumb/%s", bundleName))
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
	log.Printf("    %v", managerDockerArgs)
	log.Printf("    Done.")

	log.Printf(" |  Running manager. CTRL-C to quit.")
	err = shell.RunAndLog("docker", managerDockerArgs...)
	if err != nil {
		return err
	}
	log.Printf("    Done.")
	return nil
}
