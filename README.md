# memStress

memStress is a Linux command-line tool that creates memory pressure for testing.
It allocates a chosen amount of memory over a set time and holds it until you stop
the command. Chaos Mesh uses memStress for memory stress experiments.

## Usage

Download and extract the Linux archive for your architecture from [Releases](https://github.com/chaos-mesh/memStress/releases),
or build the binary with `make build`. On Linux, run the following command from
the directory containing `memStress`.

```sh
PATH="$PWD:$PATH" ./memStress -size 64MiB -time 5s -workers 2
```

This allocates 64 MiB total over five seconds. Each of the two workers allocates
32 MiB and holds it until you press Ctrl+C. The `PATH` setting lets memStress
find the executable when it starts its workers.

- `-size` sets the total across all workers. Specify a size such as `64MiB` or a
  percentage of the total memory reported by Linux, such as `1%`.
- `-time` sets how long it takes to reach the target size. The default, `0s`,
  allocates the memory immediately.
- `-workers` sets how many processes share the total. The default is `1`.

## Build

memStress runs on Linux. By default, the build, vet, and release targets use
`golang:<GO_VERSION>-alpine` in Docker. You can build the Linux binaries from
macOS with the same commands.

```sh
make build                    # -> ./memStress for the Docker image's architecture
make release VERSION=v0.3.1   # -> ./dist/memStress_v0.3.1-{x86_64,aarch64}-linux-gnu.tar.gz
```

Run `make` (or `make help`) to list all targets. Release tarballs contain a
single `memStress` binary at their root.
