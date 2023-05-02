# INSERT_COLLECTOR_HERE

To leverage this template, you will need to replace some things that were copied:

* module name in `go.mod`
* project name in `goreleaser.yml`
* Github release name in `goreleaser.yml` templates are not supported in this block. It must be hardcoded
* goimports local prefixes in `golangci.yml`
* binary name in `Dockerfile`
* local imports in `main.go`
* update defaults in `CollectorDefaults` as appropriate for your collector in `configuration.go`

For the most part you should be able to do a find and replace on `locomotive-collector-template` with your new collector name.

Then depending on the type of API you are integrating:

* Pick the poller or streamer integration option
* Delete all the stuff related to the other options from `main.go`, `configuration.go`, and the actual struct handling that
* ???
* PROFIT!!!

## Caveat around Batching

With the v2 integration package that was introduced, we added support for batching (off by default).
However this should be used with EXTERME caution as batching ONLY support homogenous workloads.
This means that if the workloads are the same data type, destination, customer, etc (essentially the important things in metadata) you are safe to use batching.
If any of those differ however, leave batching off because you could potentially cross contaminate or send data to a location unnecessarily.

## Make Commands

### Dependency Related Commands

#### clean

Deletes the `vendor` and `dist` for a fresh start in case some state got corrupted in either.

#### upgrade

Runs the command to ensure all transitive dependencies are pulled in and that they are on the latest "passive" version available.

#### vendor

Pulls the libraries described in `go.mod` into the `vendor` folder of this repo.

#### tidy

Cleans up unused and/or unnecessary libraries within `go.mod`.

### CI Related Commands

#### fmt

Uses tools to auto-format source files using custom tools that are not built-in to IDEs.

#### lint

Run linting tool to uncover any linting issues present in source files.

#### test

Run all test files with test coverage within source files.

#### tools

This will download tools needed for linting/testing/building/etc, anytime we have a tool/CLI that we would want to run against our code, it should be added here. If the tool is common enough, it should be added to [the CDP Collector Template's Makefile](https://github.com/FishtechCSOC/locomotive-collector-template/blob/main/Makefile).

### Build Related Commands

### build

This builds a "snapshot" version of the binary and Docker image using goreleaser but does not publish anything.

### snapshot

This is used by our CI to build and publish "snapshot" versions of your code.

**Note**

Should only be used in case CI pipelines are not working

### release

This is used by our CI to build and publish "snapshot" versions of your code.

**Note**

Should only be used in case CI pipelines are not working

## Running locally

In order to test your collector, you are able to run locally. However there are some things to take note of.

* If you want to send things through our actual dev pipeline to avoid changing code, you must have service account credentials stored locally. To do this you will need to set the environment variable `GOOGLE_APPLICATION_CREDENTIALS` with the path to your credentials file.
* If you want to just logout the results of your collector, you will need to adjust code to use either the devnull or logger dispatcher to run it.


Building and running the binary is the easiest way to debug.  Running locally in a Docker container is more difficult, but sometimes necessary.

### To run the binary directly

1. Copy a local configuration yaml to the repo root directory, and name it `configuration.yaml`
1. Run `go run cmd/<repo_name>/main.go`

### To run with Docker

Ensure that `/infra/local` has a valid `configuration.yaml` for use with your collector.

```sh
docker run -p 8888:8888 -v <PATH_TO_REPO>/infra/local:/etc/integration <docker_image>
```

**Note** 

If you run with docker and are trying to send to the dev pipeline, you will need to mount the credentials as a volume into the container:

```sh
docker run -p 8888:8888 -v <PATH_TO_REPO>/infra/local:/etc/integration -v $GOOGLE_APPLICATION_CREDENTIALS:/cred.json --env GOOGLE_APPLICATION_CREDENTIALS=/cred.json <docker_image>
```

The crucial requirement is that the GOOGLE_APPLICATION_CREDENTIALS environment variable must match the bind mount location inside the container.

### To run with Binary

Ensure that `/infra/local` has a valid `configuration.yaml` for use with your collector and copy it to the root of your repo.

Find the binary appropriate for your machine's OS in the `/dist` folder and move it to the root of your repo, alternatively use your IDE to build/run.


