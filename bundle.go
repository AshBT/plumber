package plumb

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"
)

type templateContext struct {
	Wrapper string
	Plumb   *PlumbContext
}

const bundleConfig = ".plumb.yml"

// A Dockerfile template for python bundles. We'll make sure this is
// a standalone file in the future, so we can add different language
// types.
const dockerfileTemplate = `
FROM python:2.7.10-slim

RUN mkdir -p /usr/src/bundle
WORKDIR /usr/src/bundle

COPY . /usr/src/bundle
{{ if .Plumb.Install }}
{{ range .Plumb.Install }}
RUN {{ . }}
{{ end }}
{{ else }}
RUN pip install -r requirements.txt
{{ end }}
RUN pip install bottle
EXPOSE 9800
{{ range .Plumb.Env }}
ENV {{ . }}
{{ end }}
CMD ["python", "{{ .Wrapper }}"]
`

const wrapperTemplate = `
import {{ .Plumb.Name }}
from bottle import post, route, run, request, HTTPResponse

__INFO = {'bundle': '{{ .Plumb.Name }}',
  'inputs': {
    {{ range .Plumb.Inputs }}'{{ .Name }}': '{{ .Description }}',{{ end }}
  },
  'outputs': {
    {{ range .Plumb.Outputs }}'{{ .Name }}': '{{ .Description }}',{{ end }}
  }
}

@route('/info')
def info():
    return __INFO

@post('/')
def index():
    data = request.json

    # validate input data
    if data is None:
        raise HTTPResponse(
            body = "No JSON data received; did you forget to set 'Content-Type: application/json'?",
            status = 400
        )

	# Instead of throwing errors when required fields are missing,
	# maybe we just return the payload instead and add a note?
    {{ range .Plumb.Inputs }}
    if not '{{ .Name }}' in data:
        raise HTTPResponse(
            body = "Missing required '{{ .Name }}' field in JSON data.",
            status = 400
        )
    {{ end }}
    # run the enhancer
    output = {{ .Plumb.Name }}.run(data)

    # validate output data
    {{ range .Plumb.Outputs }}
    if not '{{ .Name }}' in output:
        raise HTTPResponse(
            body = "Unexpected output; missing '{{ .Name }}' field.",
            status = 501
        )
    {{ end }}

	# Instead of assuming the output contains a superset of the input
	# fields, we can add the input fields to the output here? It's a
	# performance hit, but it means the output is *always* a superset
	# of the input.

    return output

run(host='0.0.0.0', port=9800)
`

func removeTempFile(f *os.File) {
	filename := f.Name()
	log.Printf(" |  Removing '%s'", filename)
	err := os.RemoveAll(filename)
	if err != nil {
		log.Printf("    %v", err)
	} else {
		log.Printf("    Removed.")
	}
}

func Bundle(path string) error {
	ctx := PlumbContext{}
	log.Printf("==> Creating bundle from '%s'", path)
	defer log.Printf("<== Bundling complete.")

	log.Printf(" |  Parsing bundle config.")
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, bundleConfig))
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, &ctx); err != nil {
		return err
	}
	log.Printf("    %v", ctx)

	log.Printf(" |  Making temp file for python wrapper")
	wrapper, err := ioutil.TempFile(path, "plumb")
	defer removeTempFile(wrapper)
	if err != nil {
		return err
	}
	log.Printf("    Created '%s'", wrapper.Name())

	templateCtx := templateContext{
		Wrapper: wrapper.Name(),
		Plumb:   &ctx,
	}

	log.Printf(" |  Writing wrapper.")
	tmpl, err := template.New("wrapper").Parse(wrapperTemplate)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(wrapper, templateCtx); err != nil {
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Making temp file for Dockerfile")
	dockerfile, err := ioutil.TempFile(path, "plumb")
	defer removeTempFile(dockerfile)
	if err != nil {
		return err
	}
	log.Printf("    Created '%s'", dockerfile.Name())

	log.Printf(" |  Writing Dockerfile.")
	tmpl, err = template.New("dockerfile").Parse(dockerfileTemplate)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(dockerfile, templateCtx); err != nil {
		return err
	}
	log.Printf("    Done.")

	log.Printf(" |  Building container.")
	imageName := fmt.Sprintf("plumb/%s", ctx.Name)
	cmd := exec.Command("docker", "build", "-t", imageName, "-f", dockerfile.Name(), path)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	multi := io.MultiReader(stdout, stderr)
	in := bufio.NewScanner(multi)

	for in.Scan() {
		log.Printf("    %s", in.Text())
	}
	if err := in.Err(); err != nil {
		log.Printf("    error: %s", err)
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	log.Printf("    Container 'plumb/%s' built.", ctx.Name)
	return nil
}
