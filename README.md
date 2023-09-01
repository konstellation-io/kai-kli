# Kli
|  Build status  | Relase version  |  License  |
| :---------: | :-----:   |  :--------------------:  |
| [![Tests][tests-badge]][tests-link] | [![GitHub Release][release-badge]][release-link] | [![License][license-badge]][license-link] |

|  Component  | Coverage  |  Bugs  |  Maintainability Rating  | Go Report |
| :---------: | :-----:   |  :---: |  :--------------------:  |  :---: |
|  KLI  | [![coverage][kli-coverage]][kli-coverage-link] | [![bugs][kli-bugs]][kli-bugs-link] | [![mr][kli-mr]][kli-mr-link] | [![Go Report Card][report-badge]][report-link] |

[tests-badge]: https://img.shields.io/github/actions/workflow/status/konstellation-io/kli/quality.yml
[tests-link]: https://img.shields.io/github/actions/workflow/status/konstellation-io/kli/quality.yml

[release-badge]: https://img.shields.io/github/release/konstellation-io/kli.svg?logo=github&labelColor=262b30
[release-link]: https://github.com/konstellation-io/kli/releases

[report-badge]: https://goreportcard.com/badge/github.com/konstellation-io/kli
[report-link]: https://goreportcard.com/report/github.com/konstellation-io/kli

[license-badge]: https://img.shields.io/github/license/konstellation-io/kli
[license-link]: https://github.com/konstellation-io/kli/blob/master/LICENSE

[kli-coverage]: https://sonarcloud.io/api/project_badges/measure?project=kli&metric=coverage
[kli-coverage-link]: https://sonarcloud.io/api/project_badges/measure?project=kli&metric=coverage

[kli-bugs]: https://sonarcloud.io/api/project_badges/measure?project=kli&metric=bugs
[kli-bugs-link]: https://sonarcloud.io/component_measures?id=kli&metric=Reliability

[kli-mr]: https://sonarcloud.io/api/project_badges/measure?project=kli&metric=sqale_rating
[kli-mr-link]: https://sonarcloud.io/component_measures?id=kli&metric=Maintainability

---

This repo contains a CLI to access, query and manage KAI servers.

## Installation

### From releases page

Go to [release page](https://github.com/konstellation-io/kai-kli/releases) and download the binary you prefer.

### Homebrew

```
brew install konstellation-io/tap/kli
```

### Scoop (Windows)

Scoop installation is made via a Konstellation owned bucket.

```
scoop bucket add konstellation-io https://github.com/konstellation-io/scoop-bucket.git
scoop install konstellation-io/kli
```

### Installation script

Fetch the script and execute it locally.

```
$ curl -fsSL -o get_kli.sh https://raw.githubusercontent.com/konstellation-io/kai-kli/main/scripts/get-kli.sh
$ chmod 700 get_kli.sh
$ ./get_kli.sh
```

Use it with `--help` flag to get a list of options.

```
./get_kli.sh --help
```

## Frameworks and libraries

- [gomock](https://github.com/golang/mock) a mock library.
- [spf13/cobra](https://github.com/spf13/cobra) used as CLI framework.
- [joho/godotenv](https://github.com/joho/godotenv) used to parse env files.
- [golangci-lint](https://golangci-lint.run/) as linters runner.


## Development

You can compile the binary with this command: 

```bash
./scripts/generate_release.sh <version> <executable_name>
```

Then run any command: 
```bash
./kli help

# Output: 
Use Konstellation API from the command line.

Usage:
  kli [command]

Available Commands:
  kai         Manage KAI
  server      Manage servers for kli

Flags:
  -h, --help   help for kli

Use "kli [command] --help" for more information about a command.

```

Example: 

```bash
./kli server ls

# Output
SERVER URL                                  
local* http://api.kai.local                 
int    https://api.your-domain.com 
```


### Setting Version variables

1. You can set a Version variable as a key/value pair directly: 

```bash
./kli kai version config your-version --set SOME_VAR="any value"
# Output:
# [✔] Config updated for version 'your-version'.
```

2. Add a value from an environment variable:

```bash
export SOME_VAR="any value"
./kli kai version config your-version --set-from-env SOME_VAR
# Output:
# [✔] Config updated for version 'your-version'.
```

3. Add multiple variables from a file:

```text
# variables.env file
SOME_VAR=12345
ANOTHER_VAR="some long string... "
```

```bash
./kli kai version config your-version --set-from-file variables.env
# Output:
# [✔] Config updated for version 'your-version'.
```

NOTE: `godotenv` library currently doesn't support multiline variables, as stated in
[PR #118 @godotenv](https://github.com/joho/godotenv/pull/118). Use next example as a workaround. 


4. Add a file as value:

```bash
export SOME_VAR=$(cat any_file.txt) 
./kli kai version config your-version --set-from-env SOME_VAR
# Output:
# [✔] Config updated for version 'your-version'.
```



## Testing

To create new tests install [GoMock](https://github.com/golang/mock). Mocks used on tests are generated with 
**mockgen**, when you need a new mock, add the following:

```go
//go:generate mockgen -source=${GOFILE} -destination=$PWD/mocks/${GOFILE} -package=mocks
```

To generate the mocks execute:
```sh
$ go generate ./...
```

### Run tests

```sh
go test ./...
```


## Linters

`golangci-lint` is a fast Go linters runner. It runs linters in parallel, uses caching, supports yaml config, has
integrations with all major IDE and has dozens of linters included.

As you can see in the `.golangci.yml` config file of this repo, we enable more linters than the default and
have more strict settings.

To run `golangci-lint` execute:
```
golangci-lint run
```

## Versioning lifecycle

In the development lifecycle of KLI there are three main stages depend if we are going to add a new feature, release a new version with some features or apply a fix to a current release.

### Alphas

In order to add new features just create a feature branch from master, and after merger the Pull Request a workflow will run the tests and if everything pass a new alpha tag will be created (like *v0.0-alpha.0*) and a new release will be generaged with this tag.

### Releases

After some alpha versions we can create what we call a release, and to do that we have to run manual the Release Action. This workflow will create a new release branch and a new tag like *v0.0.0*. With this tag a new release will be generated.

### Fixes

If we find out a bug in a release, we can apply a bugfix just creating a fix branch from the specific release branch, and createing a Pull Request to the same release branc. When the Pull Request is merged, after pass the tests, a new fix tag will be created just increasing the patch number of the version, and a new release will be build and released.

### Release locally for debugging

A local release can be created for testing without creating anything official on the release page.

- Make sure [GoReleaser](https://goreleaser.com/install/) is installed
- Run: 
    `goreleaser --skip-validate --skip-publish --rm-dist`
- Find the built binaries under `dist/` folder.
