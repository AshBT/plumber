# Plumber manager
The manager for plumber is included in the `plumber` distribution by packaging this folder using [go-bindata](https://github.com/jteeuwen/go-bindata). This folder exists so that the build chain tests it, but the `plumber` tool does not link against it.

Instead, the `plumber` tool contains a copy of these source files in binary form and uses the source files to build the manager inside a Docker container.
