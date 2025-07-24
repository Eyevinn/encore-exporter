# go-default

This is a default setup for Go projects with name github.com/Eyevinn/name

Use `cmd/` as root for programs (each in its own subdirectory)
Use `internal` directories for code that should not be exported.
Use `testdata` sub-directories for test files.

## Get started

You must edit the files so that the name of this repo gets inserted in the right places.
This includes replacing the string `REPO_NAME` in the Makefile and in `cmd/example/main.go`
with the name of this repo. 

1. Run `go mod init github/Eyevinn/REPO_NAME` where `REPO_NAME` is the name of this repo.
2. Edit the installed files appropriately
   * Replace `REPO_NAME` with the name of the REPO
   * README.md - See [mp4ff][mp4ff] for an example. (The default is rather Node-oriented)
   * CHANGELOG.md
   * cmd/example.go (fill in proper name to )
   * Makefile (with proper commands)
3. Place code for executables in directories with proper names under cmd/example etc
   and update the Makefile
4. For code coverage (open source projects), you can run `make coverage`.
   * You can also set up `coveralls.io` in the sam way as [mp4ff][coveralls].
5. Later you can run the license check program `wwhrd` for third-party licenses
   using `make check-licenses` after installing `wwhrd`.
   

## Included

The defaults for all go projects include:

- A .gitignore file
- Github actions for running tests and golang-ci-lint
- A Makefile for running tests, coverage, and update dependencies.
- An example command line application cmd/example with a version fetched from internal/version.go.
- A README skeleton (update badges to Go, see e.g. mp4ff)
- A CHANGELOG.md file that should be changed manually
- A config file for pre-commit (see https://pre-commit.com)

[mp4ff]: https://github.com/Eyevinn/mp4ff
[coveralls]: https://coveralls.io/github/Eyevinn/mp4ff