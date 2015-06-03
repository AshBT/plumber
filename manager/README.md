# Plumb manager
The manager for plumb is included in the `plumb` distribution by packaging this folder using [go-bindata](https://github.com/jteeuwen/go-bindata). This folder exists so that the build chain tests it, but does the `plumb` tool does not link against it.

 Instead, the `plumb` tool contains a copy of these source files in binary form and uses the source files to build the manager inside a Docker container.
