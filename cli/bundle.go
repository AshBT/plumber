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
	"github.com/qadium/plumber/shell"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"
)

type templateContext struct {
	Wrapper string
	Plumber *Bundle
}

// A Dockerfile template for python bundles. We'll make sure this is
// a standalone file in the future, so we can add different language
// types.
const dockerfileTemplate = `
FROM python:2.7.10-slim

RUN mkdir -p /usr/src/bundle
WORKDIR /usr/src/bundle

COPY . /usr/src/bundle
{{ if .Plumber.Install }}
{{ range .Plumber.Install }}
RUN {{ . }}
{{ end }}
{{ else }}
RUN pip install -r requirements.txt
{{ end }}
RUN pip install bottle
EXPOSE 9800
{{ range .Plumber.Env }}
ENV {{ . }}
{{ end }}
CMD ["python", "{{ .Wrapper }}"]
`

const wrapperTemplate = `
import {{ .Plumber.Name }}
import datetime
from bottle import post, route, run, request, HTTPResponse

__INFO = {'bundle': '{{ .Plumber.Name }}',
  'inputs': {
    {{ range .Plumber.Inputs }}'{{ .Name }}': '{{ .Description }}',{{ end }}
  },
  'outputs': {
    {{ range .Plumber.Outputs }}'{{ .Name }}': '{{ .Description }}',{{ end }}
  }
}

# set of expected output fields
__EXPECTED_OUTPUTS = set([{{ range .Plumber.Outputs }}"""{{ .Name }}""",{{end}}])

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

	# first, get the program data
	program_data = data.get("data", None)
	if program_data is None:
		raise HTTPResponse(
			body = "Unexpected format for an enhancer; did you forget to wrap your data in a 'data' field?",
			status = 400
		)

	# ensure its metadata, errors, and history fields are set
	if "metadata" not in data:
		data["metadata"] = {}
	if "errors" not in data["metadata"]:
		data["metadata"]["errors"] = []
	if "history" not in data["metadata"]:
		data["metadata"]["history"] = []

	# now, get only the fields needed for this enhancer
	input = {
		{{ range .Plumber.Inputs }}"""{{ .Name }}""": program_data.get("""{{ .Name }}""", None), {{ end }}
	}

	# run the enhancer
	try:
		output = {{ .Plumber.Name }}.run(input)
	except Exception as e:
		error = {
			"bundle-name": """{{ .Plumber.Name }}""",
			"timestamp": datetime.datetime.now().isoformat(),
			"message": str(e)
		}
		data["metadata"]["errors"].append(error)
		return data

	# discard any updates to the inputs
	{{ range .Plumber.Inputs }}
	output.pop("""{{ .Name }}""")
	{{ end }}

	# update the program data
	for key in output:
		if key not in __EXPECTED_OUTPUTS:
			error = {
				"bundle-name": """{{ .Plumber.Name }}""",
				"timestamp": datetime.datetime.now().isoformat(),
				"message": "Field '{}' was not in set of expected outputs.".format(key)
			}
			data["metadata"]["errors"].append(error)
			continue
		history = {
			"bundle-name": """{{ .Plumber.Name }}""",
			"field-name": key
		}
		if key in data["data"]:
			history["action"] = "update"
			history["prev"] = data["data"][key]
			history["timestamp"] = datetime.datetime.now().isoformat()
		else:
			history["action"] = "new"
			history["prev"] = None
			history["timestamp"] = datetime.datetime.now().isoformat()

		data["data"][key] = output[key]
		data["metadata"]["history"].append(history)
	return data

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

// Bundle stuff...
// BUG(echu): need to figure out how to handle conflicts in bundle names
func (ctx *Context) Bundle(bundlePath string) error {
	log.Printf("==> Creating bundle from '%s'", bundlePath)
	defer log.Printf("<== Bundling complete.")

	log.Printf(" |  Parsing bundle config.")
	bundleConfig, err := ParseBundleFromDir(bundlePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("    Could not find '%s' in path '%s'.", bundleConfig, bundlePath)
			log.Printf("    Did not bundle '%s'.", bundlePath)
			return nil
		}
		return err
	}
	log.Printf("    %v", bundleConfig)

	log.Printf(" |  Making temp file for python wrapper")
	wrapper, err := ioutil.TempFile(bundlePath, "plumber")
	defer removeTempFile(wrapper)
	if err != nil {
		return err
	}
	log.Printf("    Created '%s'", wrapper.Name())

	templateCtx := templateContext{
		Wrapper: path.Base(wrapper.Name()),
		Plumber: bundleConfig,
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
	dockerfile, err := ioutil.TempFile(bundlePath, "plumber")
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
	err = shell.RunAndLog(ctx.DockerCmd, "build", "--pull", "-t", ctx.GetImage(bundleConfig.Name), "-f", dockerfile.Name(), bundlePath)
	if err != nil {
		return err
	}
	log.Printf("    Container '%s' built.", ctx.GetImage(bundleConfig.Name))
	return nil
}
