# Link plumber
The `plumb` tool is designed to facilitate easy deployment of data enhancers and linkers. It is designed to be language-agnostic and allows developers the flexibility to build tools in whichever language they are familiar with.

Its primary goal is to facilitate knowledge discovery through simple deployment of "hypotheses" and "conclusions" in the form of "enhancers" and "linkers". Developers create these enhancers and linkers by adhering to a simple programmatic interface and providing a YAML config file in their repository.

This interface is simple: JSON in and JSON out. The `plumb` tool will take care of creating the necessary wrappers to enable use in the `plumb` ecosystem.

## Installation
Download binaries.

# Examples

## Server example
To start a `plumb` scheduler, run `plumbd` to start the server daemon. This does not need to be run locally (except for development purposes). When `plumb submit` is invoked, it sends graph and dependency information to `plumbd`. Every time a new `plumb` is submitted to the daemon, the dependency graph is recomputed.

If the named pipeline is currently running, then any upstream dependencies will be paused and buffered while the new `plumb` is installed. You can give this a whirl at `plumb-demo.qadium.com` by running commands with `plumb --server plumb-demo.qadium.com submit foobar .`. You can direct your browser to `plumb-demo.qadium.com` to see the results.

## Enhancers and linkers
An "enhancer" is defined as a compute node that takes some JSON in, adds information to it, and returns it as more JSON. A linker.... TBD.

### Data sources
A data source, such as a database or the Twitter firehose, is an enhancer with no upstream dependencies. You can see an example of a data source via

    git clone github.com/qadium/plumb-twitter
    cd plumb-twitter && plumb bundle .

You'll need to add your own Twitter credentials in the `.plumb.yml` file before running `plumb bundle`. You can start the source with `docker run plumb/twitter`. Note that this subscribes to the Twitter feed, but does not do anything with it.

### Data sinks
A data sink, such as a database writer or a GUI display, is an enhancer with no downstream dependencies. You can see an example of a data sink via

  git clone github.com/qadium/plumb-events
  cd plumb-events && plumb bundle .

You can start the sink with `docker run -p 8000:8000 -p 9100:9100 plumb/events`. If you point your browser to `localhost:8000`, you'll see a moving timeline with heights corresponding to requests per minute.

You can `curl localhost:9100 -d '{}' -H 'Content-Type: applicatoin/json'` to see an event recorded. Note that sinks receive *any* input, since the idea is to uniformly handle all messages that arrive at the destination. Plumb does not support multiple heterogenous sinks at the moment; you can, however, write the same data to multiple databases if you'd like.

### Data enrichers
To add an enhancer that adds twitter data to a pipeline that produces JSON with a `twitter` key, first

    git clone github.com/qadium/twitter-linker
    cd twitter-linker && plumb bundle .

This will package the source code in the twitter-linker into a Docker container which you can run locally with `docker run twitter-linker`.

You can then `curl localhost:8000 -d '{"twitter": "qadium"}' -H 'Content-Type: application/json'` and receive the output from the twitter linker.

This will additionally run any unit tests you have specified in the YAML file as well.

To schedule it into a production pipeline, run `plumb create foobar` and `plumb submit foobar .`. To run the pipeline, run `plumb start foobar`.

## Put it all together
Let's put the source, the sink, and the enhancers together into a data pipeline.

    plumb create hello-pipe
    plumb submit hello-pipe plumb-twitter plumb-events plumb-hello
    plumb start hello-pipe

Now navigate to `localhost:8000` and watch the count of tweets coming in!

### Some details
The command `plumb create` creates a bolt bucket to store the enhancers in the pipeline. When submitting, `plumb` will recompute the dependency graph and modify the key-val store for that bucket. (XXX: Does this mean we no longer need a server / daemon to manage all this stuff?)

The `plumb start` command will query the bucket for all its key vals and construct the proper order for starting the containers and hook them up to talk to each other in the right order

# Roadmap
*v0.1.0*

- `plumb bundle` functionality for Python
- `plumb submit` to a local `plumbd` instance
- `plumb start` on local docker instances
- `plumbd` as nothing more than docker wrapper
- unit tests
- bintray deploy?

*v0.2.0*

- `plumb bundle` runs unit tests
- `plumbd` handles dependency graphs
- `plumb install` for different languages
- automatic discovery of inputs and outputs based on test cases?

*v0.3.0*

- `plumb compile`?
