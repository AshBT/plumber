package cli

import (
	"fmt"
	"log"
//	"gopkg.in/yaml.v2"
//	"io/ioutil"
	"text/template"
	"path/filepath"
	"github.com/qadium/plumb/graph"
	"github.com/qadium/plumb/shell"
	"os/exec"
	"os/signal"
	"os"
)

type pipelineContext struct {
	Pipeline []string	// host:port
}

const manager = `
package main

import (
	"net/http"
	"log"
	"io"
	"io/ioutil"
	"os/signal"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var err error
	client := &http.Client{}
	body := r.Body
	defer body.Close()

	var (
		req *http.Request
		resp *http.Response
	)

	{{ range .Pipeline }}
	req, err = http.NewRequest("POST", "{{ . }}", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	body = resp.Body
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	{{ end }}

	final, err := ioutil.ReadAll(io.LimitReader(body, 1048576))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(final)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<- c
		log.Printf("Received termination; qutting.")
		os.Exit(0)
	}()


	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":9800", nil))
}
`

func contextsToGraph(ctxs []*Context) []*graph.Node {
	nodes := make([]*graph.Node, len(ctxs))
	m := make(map[string]int) 	// this map maps inputs to the index of the
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
	pipectx := pipelineContext{}
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
			cmd := exec.Command("docker", "kill", string(containerId)[0:4])
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
		pipectx.Pipeline = append(pipectx.Pipeline, fmt.Sprintf("http://192.168.59.103:%s", string(portNum[:len(portNum)-1])))
	}
	log.Printf("    %v", pipectx)
	log.Printf("    Done.")

	log.Printf(" |  Writing pipeline manager.")
	managerFile := fmt.Sprintf("%s/manager.go", path)
	file, err := os.Create(managerFile)
	if err != nil {
		return err
	}

	tmpl, err := template.New("manager").Parse(manager)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(file, pipectx); err != nil {
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Compiling pipeline manager.")
	managerBinary := fmt.Sprintf("%s/manager", path)
	err = shell.RunAndLog("go", "build", "-o", managerBinary, managerFile)
	if err != nil {
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Running manager. CTRL-C to quit.")
	cmd := exec.Command(managerBinary)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		// wait for signal from CTRL-C
		for range sig {
			log.Printf("    Received CTRL-C; terminating manager.")
			cmd.Process.Signal(os.Interrupt)
		}
	}()

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	// stop listening to CTRL-C now
	signal.Stop(sig)
	log.Printf("    Done.")
	return nil
}
