# Plumb
The `plumb` tool is designed to facilitate knowledge discovery through the deployment of data enhancers and linkers for the express purpose of exploring data and (informally) testing hypotheses.

Data processing tasks for ETL or data science typically involve data cleaning, data munging, or just adding new fields to a pre-set schema. Often, the trick is corralling raw data from a variety of sources into a form usable by algorithms that typically operate on floating point numbers or categorical data.

This process can be thought of as a series of operations or transformations on raw data: we term these *enhancers*. Each enhancer can be as simple (e.g., a regex match) or as complex (e.g., a database lookup) as necessary to provide additional data. The only requirement for enhancers is that they take a map in and provide a map out.

Furthermore, we provide some higher-level capabilities by allowing developers to deploy nodes that link related data together: we call these nodes *linkers*.

Based on information provided in a `.plumb.yml` file, the `plumb` tool can compose a set of enhancers or linkers such that their dependencies are satisfied and deploy the data processing pipeline to a (properly configured?) CoreOS cluster.

## Installation
We do not currently have a binary for you to download. The preferred way to install `plumb` at the moment is

    git clone https://github.com/qadium/plumb
    cd plumb && make && make install

You can also use

    go get github.com/qadium/plumb

However, this will not display the git SHA1 information with `plumb version`.

## Enhancers and linkers

Developers create enhancers and linkers by adhering to a simple programmatic interface and providing a YAML config file in their repository. The requirement is simple: implement a `run` function that takes a map (or dictionary) in and returns a new map.

The `plumb` tool will take care of creating the necessary wrappers to enable use in the `plumb` ecosystem.

An "enhancer" is defined as a compute node that takes some JSON in, adds information to it, and returns it as more JSON. A linker.... TBD.

## Testing
By decoupling the transformations on each piece of data, we can also programmatically test and document enhancers and linkers. This gives end-users a high level of assurance that their data processing pipeline is correct, preventing garbage in and garbage out.

## Infrastructure (TODO)
While `plumb` can run on a single machine, it is best run on a cluster hosted on AWS or GCE. We are planning support for Kubernetes and CoreOS. For now, use the included Terraform file for bootstrapping.

# Examples

## Data sources (TODO)
A data source, such as a database or the Twitter firehose, is an enhancer with no upstream dependencies. You can see an example of a data source via

    git clone github.com/qadium/plumb-csv
    cd plumb-csv && plumb bundle .

You can start the source with `docker run plumb/csv`. Note that this rotates through the included CSV file, but does not do anything with it.

## Data sinks (TODO)
A data sink, such as a database writer or a GUI display, is an enhancer with no downstream dependencies. You can see an example of a data sink via

  git clone github.com/qadium/plumb-count
  cd plumb-count && plumb bundle .

You can start the sink with `docker run -p 8000:8000 -p 9800:9800 plumb/count`. If you point your browser to `localhost:8000`, you'll see a moving timeline with heights corresponding to requests per minute.

You can `curl localhost:9800 -d '{}' -H 'Content-Type: application/json'` to see an event recorded. Note that sinks receive *any* input, since the idea is to uniformly handle all messages that arrive at the destination. Plumb does not support multiple heterogenous sinks at the moment; you can, however, write the same data to multiple databases if you'd like.

## Data enrichers
To add an enhancer that adds a `hello` text to a JSON containing the `name` field, first

    git clone github.com/qadium/plumb-hello
    cd plumb-hello && plumb bundle .

This will package the source code into a Docker container which you can run locally with `docker run plumb/hello -p 9800`.

You can then `curl localhost:9800 -d '{"name": "qadium"}' -H 'Content-Type: application/json'` and receive the output from the greeter.

(TODO) This will additionally run any unit tests you have specified in the YAML file as well.

## Put it all together
Let's put the source, the sink, and the enhancers together into a data pipeline.

    plumb create hello-pipe
    plumb submit hello-pipe plumb-twitter plumb-events plumb-hello
    plumb start hello-pipe

Now navigate to `localhost:8000` and watch the count of tweets coming in!

### Some details
The command `plumb create` creates a bolt bucket to store the enhancers in the pipeline. When submitting, `plumb` will recompute the dependency graph and modify the key-val store for that bucket. (XXX: Does this mean we no longer need a server / daemon to manage all this stuff?)

The `plumb start` command will query the bucket for all its key vals and construct the proper order for starting the containers and hook them up to talk to each other in the right order.

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
- `plumb start` on GCE with kubernetes
- unit tests
- bintray deploy?

*v0.2.0*

- `plumb bundle` runs unit tests
- `plumbd` handles dependency graphs
- `plumb install` for different languages
- automatic discovery of inputs and outputs based on test cases?

*v0.3.0*

- `plumb compile`?

# Internals / notes
What I'd really like is for `plumb push` to just create a bunch of Fleet service files, push the images, and submit the service files to a remote server. The whole thing is standalone at that point and can run without any other dependencies. A `plumb clone` command can clone a pipeline while a `plumb pull` can update your pipeline. This also allows pipeline diffs and rollbacks.

A `plumb create` command will create a local `git` repo for each pipeline.

This means you can have **distributed data processing pipelines**. (Issue: do we store environment variables that might contain secret keys?).

Another thing to "share" are the nodes themselves; this might make sense to share via Docker registries. After `plumb bundle` is called, we might want to share the bundles with other people. This requires some sort of central registry.

Sharing within an organization can happen via private git servers and private docker registries; do we care to share publicly on Github? It's easy to do if all you do is expose the `.yml` files. It's easy to expose bundles--not easy to expose pipelines. Do we care to expose pipelines?

One possibility is to host bundles on Github and "search" over that. The images are then "built" and stored on private registries; so we don't have any public images--although it's not hard to have public images.

One final thought about "sharing"--let's punt for now. We know it can be done via github, but let's build the core functionality. Let's `bundle` enhancers, `create` pipelines, `submit` bundles to a pipeline, and then `start` something. All in a Vagrant CoreOS cluster.

The `start` command will compute the dependency graph, build a "coordinator" container, and handle requests. So you can curl the coordinator container (say, at "coordinator.foo.org"), and it will return "coordinator.foo.org/1" for the first piece of data you sent. Navigating to the URL will "spin" until the coordinator receives data. It will (at some point) use an LRU cache. But it will evict data that it cannot hold in memory. If you try to access evicted data, mmm... it will return 404.
