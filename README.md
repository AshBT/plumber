# Link
The `link` tool is designed to facilitate easy deployment of data enhancers and linkers. It is designed to be language-agnostic and allows developers the flexibility to build tools in whichever language they are familiar with.

Its primary goal is to facilitate knowledge discovery through simple deployment of "hypotheses" and "conclusions" in the form of "enhancers" and "linkers". Developers create these enhancers and linkers by adhering to a simple programmatic interface and providing a YAML config file in their repository.

This interface is simple: JSON in and JSON out. The `link` tool will take care of creating the necessary wrappers to enable use in the `link` ecosystem.

## Enhancers and linkers
An "enhancer" is defined as a compute node that takes some JSON in, adds information to it, and returns it as more JSON. A linker.... TBD.

## Example
To add an enhancer that adds twitter data to a pipeline that produces JSON with a `twitter` key, first

    git clone github.com/qadium/twitter-linker
    cd twitter-linker && link .

This will package the source code in the twitter-linker into a Docker container which you can run locally with `docker run twitter-linker`.

You can then `curl localhost:8000 -d "{'twitter': 'qadium'}"` and receive the output from the twitter linker.

This will additionally run any unit tests you have specified in the YAML file as well.

To schedule it into a production pipeline, run `link create foobar` and `link submit foobar .`. To run the pipeline, run `link start foobar`.

## Server example
To start a `link` scheduler, run `linkd` to start the server daemon. This does not need to be run locally (except for development purposes). When `link submit` is invoked, it sends graph and dependency information
to `linkd`. Every time a new `link` is submitted to the daemon, the dependency graph is recomputed.

If the named pipeline is currently running, then any upstream dependencies will be paused and buffered while the new `link` is installed. You can give this a whirl at `link-demo.qadium.com` by running commands with `link --server link-demo.qadium.com submit foobar .`.
