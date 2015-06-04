# Plumber
Plumber is a tool to facilitate **distributed data exploration**. It comes with a `plumber` command line tool to deploy and manage data enhancers on a [Kubernetes](https://github.com/GoogleCloudPlatform/kubernetes) cluster.

Based on information provided in a `.plumb.yml` file, `plumber` can compose a set of enhancers such that their dependencies are satisfied and deploy the data processing pipeline to a Kubernetes cluster. At the moment, we only support a local deploy (with Docker) and a cloud deploy with Google Cloud.

## Rationale
Data processing tasks for ETL or data science typically involve data cleaning, data munging, or just adding new fields to a pre-set schema. Often, the trick is corralling raw data from a variety of sources into a form usable by algorithms that typically operate on floating point numbers or categorical data.

This process can be thought of as a series of operations or transformations on raw data: we term these *enhancers*. Each enhancer can be as simple (e.g., a regex match) or as complex (e.g., a database lookup) as necessary to provide additional data. The only requirement for enhancers is that they take a map in and provide a map out.

## Installation
Download somewhere.

You can also use

    go get github.com/qadium/plumber

However, this will not display the git SHA1 information with `plumber version`.

### Prerequisites
You'll need `git` and `docker` installed on the command line. For use with Google Cloud, you'll need the the Google Cloud SDK command line tools. You'll also need to make sure you installed kubernetes via `gcloud`.

## Enhancers and linkers

Developers create enhancers and linkers by adhering to a simple programmatic interface and providing a YAML config file in their repository. The requirement is simple: implement a (public) `run` function that takes a map (or dictionary) in and returns a new map.

The `plumber` tool will take care of creating the necessary wrappers to enable use in the `plumber` ecosystem.

## Testing
By decoupling the transformations on each piece of data, we can also programmatically test and document enhancers and linkers. This gives end-users a high level of assurance that their data processing pipeline is correct, preventing garbage in and garbage out.

# Examples
## Hello, world
First, you'll need to bootstrap `plumber` by creating a `manager` container.

    plumber bootstrap

Next, you'll need some data enhancers.

### Data enhancers
For the "hello, world" demo of `plumber`, you will need to clone two repositories:

    git clone github.com/qadium/plumber-hello
    git clone github.com/qadium/plumber-host

The `plumber-hello` repository contains a piece of code that reads in a field `name` from an input JSON and adds a field `hello` that contains the text `Hello, {name}, my name is {my_name}`.

The `plumber-host` respository contains a piece of code that reads in a field `hostname` from an input JSON and adds a field `name` that contains the resolution of that hostname to an IP address. If no IP can be resolved, it uses a default character, such as "?" for unknown.

First, we'll play with the `plumber bundle` command. After cloning the repositories, run

    cd plumber-hello && plumber bundle .

This will package the source code into a Docker container which you can run locally with `docker run -p 9800:9800 plumber/hello`.

You can then `curl localhost:9800 -d '{"name": "qadium"}' -H 'Content-Type: application/json'` and receive the output from the greeter, which should look something like

    {"name": "qadium", "hello", "Hello, qadium, my name is bazbux"}

You can do the same with `plumber-host`. First, navigate to its directory and run `plumber bundle .`. You can run a container locally with `docker run -p 9800:9800 plumber/host`. You can again use `curl` to send data to the server and see its response.

### Pipelines
We now create a pipeline

    plumber create foo

To this pipeline, we add our two bundles

    plumber add foo PATH/TO/plumber-host PATH/TO/plumber-hello

This will copy their `.plumb.yml` files into the `foo` pipeline directory. Note that you can add them before bundling or bundle them then add to a pipeline.

Finally, ensure that `plumber-host` and `plumber-hello` have had their pipelines built via `plumber bundle`, now run

    plumber start foo

Along with the `plumber/host` and `plumber/hello` containers, this will start a `manager` container which forwards requests to `plumber/host` and `plumber/hello` in an order that satisfies the dependency graph. You can now run

    curl localhost:9800 -d '{"hostname": "qadium.com"}' -H 'Content-Type: application/json'

This should produce a response

    {"hostname": "qadium.com", "name": "54.67.80.178", "hello": "Hello, 54.67.80.178, my name is bazbux"}

If you're running on OSX, replace the curl command with

    curl `boot2docker ip`:9800 -d '{"hostname": "qadium.com"}' -H 'Content-Type: application/json'

### Run on Google Cloud
Running on Google Cloud is very straightforward. First, ensure you have an account and have installed the Google Cloud SDK. Log in with

    gcloud auth login

Ensure you've installed the alpha components with

    gcloud components update alpha

Now, create a container cluster (if you haven't already)

    gcloud alpha container clusters create ...

Finally, start your pipeline on Google Cloud with

    plumber start --gce PROJECT-ID foo

The `PROJECT-ID` for your cloud can be found from the Google Cloud developer console.

## A note on data sources and sinks
Data sources and sinks do not fit nicely into our model of "JSON in, JSON out", since a data source is essentially "nothing in, JSON out," and a data sink is "JSON in, nothing out." While we plan to provide support for creating data sources and sinks, these can be emulated with simple HTTP get requests (data source) and proper handling of the response (data sink).

# Command line tool
Here's the help-text for `plumber`
```
NAME:
   plumber - a command line tool for managing distributed data pipelines

USAGE:
   plumber [global options] command [command options] [arguments...]

VERSION:
   0.0.1-dev

AUTHOR(S):

COMMANDS:
   add		add a plumber-enabled bundle to a pipeline
   create	create a pipeline managed by plumber
   bootstrap	bootstrap local setup for use with plumber
   start	start a pipeline managed by plumber
   bundle	bundle a node for use in a pipeline managed by plumber
   version	more detailed version information for plumber
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

# Alternatives
## Storm
Storm topologies... very similar, not as dynamic. Not container based. Performance?

## Docker compose aka Fig
For more generic services; explicit linking. Full control of docker containers.

## Others?
I don't know of any others.

# Roadmap
*v0.1.0*

- (✓) `plumb bundle` functionality for Python
- ~~`plumb submit` to a local `plumbd` instance~~
- (✓) `plumb start` on local docker instances
- ~~`plumbd` as nothing more than docker wrapper~~
- (✓) `plumb start` on GCE with kubernetes
- unit tests
- bintray deploy?

*v0.2.0*

- `plumb bundle` runs automated unit tests
- ~~`plumbd` handles dependency graphs~~
- `plumb install` for different languages
- automatic discovery of inputs and outputs based on test cases?

*v0.3.0*

- `plumb compile`?
