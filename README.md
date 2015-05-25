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

# Who handles what?
The `plumb` command *should* handle bundling and packaging, since this is the safest way to keep secrets from leaking in plain-text. Also, this lets people build and play around locally.

Unfortunately, it means that things like `docker`, `fleet`, `git`, etc. must be installed on the host machines. It's not a terribly tall order, but it can be annoying to install all the dependencies.

We can work around this requirement by letting `plumb` commands talk to a remote server that handles *all the work*. Unfortunately, this doesn't work as nicely, since someone can't play around without an internet connection. Also, if the connection is broken, etc. we have to recover, etc. etc.

Thus, it's safest to let all `plumb` commands work locally. Furthermore, we can work around this by hosting `plumb` on a server and letting folks `git push` to it. The only interface to `plumb` in this case is a `git push`. I'm happy to let that hide all the gory details.

Now, here's the big issue: how does one share pipelines and dependency graphs? The `plumb` tool can certainly compute it and store local information in a BoltDB, but if multiple users want to collaborate on a pipeline, they'd have to share the BoltDB--and the bundled images! This is where having a `plumbd` acting as a server might be a good idea. Plus, `git push` no longer works nicely since the user might want to be able to look at the graph from the command line, etc.

I can think of two solutions:

1. The `plumb` tool is stand alone; we can host `plumb` and its monitoring UIs on our cluster and show IT staffs how to install it. The interface is just a UI--so `git push` gets translated into a nice UI. The developers only need to write `.plumb.yml` files and `git push` to our server. (Kind of like Travis for data.)
2. We write a `plumbd` server that provides this information.

The first is `plumb` as a service, while the second exposes all the internals of `plumb` to users (sort of).

I don't know the right solution, but I do think simplicity is key. Users need to know roughly what happens when their code is run through the `plumb` tool. It's also nice to be able to run it themselves on a local machine.

What I'd really like is for `plumb submit` to just create a bunch of Fleet service files, push the images, and submit the service files to a remote server. The whole thing is standalone at that point and can run without any other dependencies.

The graph is stored locally, so you can see all that stuff locally, of course, but not remotely. One possible solution is to store all the graph / dependency metadata into `etcd` or just backup the `db` onto the server in its entirety (but that option copies over all graphs and not just the one you submitted). When another user connects to the server, the `plumb get` (or something) command queries `etcd` for the relevant info and makes a copy in Bolt. This means a user can actually have "graphs" that live on two different servers--kind of cool (and git-like). You'd have to do a `plumb sync` (or `plumb update`) to ensure your local graph is consistent with the server's graph, but that's kind of nice--not sure if Bolt has "commit" semantics, so I'm not sure if it's possible to rollback "diffs".

This means you can have **distributed data processing pipelines**. An alternative to using Bolt is to use `git` to track changes to the processing graph. We can keep that information in a *text file* under version control. Then the "graph" can be stored on the server via `git` and cloned by other users. One possibility is for this to just be a copy of all the `yml` files for the enhancers in the pipeline. Version control can add or remove files and do diffs pretty easily. The local `plumb` tool can recreate the graph for its local BoltDB. (This means IDs, etc. must be hashed instead of uniquely generated.)

Any modifications are done with the *client* talking to individual units in the server. This does open up the possibility of mutiple clients attempting to modify the graph *at the same time*. But how often does that happen with `git`? Furthermore, we could disallow multiple users from modifying the graph by requiring that all users have their graphs in sync--so a `plumb sync` will destroy (?) your local graph and sync you with the server.

## Ok, now what?

So after much internal debate: it seems like a good design for sharing data processing pipelines is the `git` model. Individual developers can have their own data processing pipelines to explore ideas and hypotheses. These pipelines are backed by Docker and must be run on CoreOS. When folks are happy with it, something like `plumb submit` can deploy the thing for real.

When something like `plumb submit` is called, this also sends metadata (TBD) to the server--assuming the server is ready for the "fast-forward". This metadata ought to be sufficient for anyone else to reconstruct the processing pipeline with `plumb pull` or `plumb update`. The best design for this is to let `git` handle changes--so it's best to store the metadata in a `git` repository (Issue: do we store environment variables that might contain secret keys?). I'm no longer sure if this requires a BoltDB, unless a lot of information must be computed about the graph each time.

This means `plumb create` needs to create a local `git` repo for each pipeline. And `plumb clone` will also create a local `git` repo for the pipeline it's trying to pull. The `plumb` tool can manage all the commits (and messages) locally. This puts the pipeline under version control.
