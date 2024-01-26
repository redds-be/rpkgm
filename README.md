<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/redds-be/redd-go-template">
    <img src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_LightBlue.png" alt="Logo" width="128" height="128">
  </a>

<h3 align="center">redd's Go Template</h3>

  <p align="center">
    A template for my go projects.
    <br />
    <br />
    <a href="#">View Demo</a>
    ·
    <a href="https://github.com/redds-be/redd-go-template/issues">Report Bug</a>
    ·
    <a href="https://github.com/redds-be/redd-go-template/issues">Request Feature</a>
  </p>
</div>

<!-- PROJECT SHIELDS -->
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/redds-be/redd-go-template/golangci-lint.yml?label=Golangci-lint)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/redds-be/redd-go-template/gotest.yml?label=Go%20test)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/redds-be/redd-go-template/gobuild.yml?label=Go%20build)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/redds-be/redd-go-template/docker-build-test.yml?label=Docker%20build)
![GitHub pull requests](https://img.shields.io/github/issues-pr/redds-be/redd-go-template)
![GitHub issues](https://img.shields.io/github/issues/redds-be/redd-go-template)
![GitHub License](https://img.shields.io/github/license/redds-be/redd-go-template)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/redds-be/redd-go-template)

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#features">Features</a></li>
    <li>
        <a href="#usage">Usage</a>
        <ul>
            <li><a href="#source">Source</a></li>
            <li><a href="#docker">Docker</a></li>
        </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About redd-go-template

A template for my go projects. It contains everything I need to start a new Go project.

This template (helloworld) is an example of a [hexagonal architecture](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3) implementation.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Features

- A `Hello, World!` message.
- HTTP server to serve the `Hello, World!` at "/" endpoint
- sqlite database to keep track of the history (date+helloworld)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage

### Source

Just run it:

```console
user@host:~$ go run cmd/main.go
```

Or compile and run it:

```console
user@host:~$ go build -o helloworld cmd/main.go && ./helloworld
```

### Docker

See [docker instruction](docker/).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ROADMAP -->
## Roadmap

- [ ] Upgrade the template

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

I don't expect anyone other than me to contribute, but you should follow these steps :

**Fork -> Patch -> Push -> Pull Request**

The **Go** code is linted with [`golangci-lint`](https://golangci-lint.run) and
formatted with [`golines`](https://github.com/segmentio/golines) (width 120) and
[`gofumpt`](https://github.com/mvdan/gofumpt). See the Makefile targets.
If there are false positives, feel free to use the
[`//nolint:`](https://golangci-lint.run/usage/false-positives/#nolint-directive) directive
and justify it when committing to your branch or in your pull request.

For any contribution to the code, make sure to create tests/alter the already existing ones according to the new code.

Make sure to run `make prep` before committing any code.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## License

*Project under the [GPLv3 License](https://www.gnu.org/licenses/gpl-3.0.html).*

*Copyright (C) 2024 redd*

<p align="right">(<a href="#readme-top">back to top</a>)</p>
