# CD/CD go template
CD/CD go template and test useful or not 

[![GoDoc](https://godoc.org/github.com/things-labs/cicd-go-template?status.svg)](https://godoc.org/github.com/things-labs/cicd-go-template)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-labs/cicd-go-template?tab=doc)
[![codecov](https://codecov.io/gh/things-labs/cicd-go-template/branch/main/graph/badge.svg)](https://codecov.io/gh/things-labs/cicd-go-template)
![Action Status](https://github.com/things-labs/cicd-go-template/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-labs/cicd-go-template)](https://goreportcard.com/report/github.com/things-labs/cicd-go-template)
[![Licence](https://img.shields.io/github/license/things-labs/cicd-go-template)](https://raw.githubusercontent.com/things-labs/cicd-go-template/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-labs/cicd-go-template)](https://github.com/things-labs/cicd-go-template/tags)

This is template that help you to quick implement some library using Go.

This repository is contains following.

- CI/CD
    - golangci-lint
    - go test
    - CodeQL Analysis (Go)
    - dependabot for github-actions and Go

## How to use
1. action Use this template and then create a repository
2. replace "things-labs" to your self username using sed(or others)
3. run make init 
4: done
   
## Features


## Usage

### Installation

Use go get.
```bash
    go get github.com/things-go/cicd-go-template
```

Then import the modbus package into your own code.
```bash
    import modbus "github.com/things-go/cicd-go-template"
```

### Example

[embedmd]:# (_examples/main.go go)
```go

```

## References
- [go-lib-template](https://github.com/skanehira/go-lib-template)

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.