# waluigi

[![Go Reference](https://pkg.go.dev/badge/github.com/giantswarm/waluigi.svg)](https://pkg.go.dev/github.com/giantswarm/waluigi)


CLI tool to pretty print the logs from CAPI controllers.

It was vibe coded using ChatGTP. Feel free to burn the code, change it or rewrite it.

## Installation
```bash
go install github.com/giantswarm/waluigi
```

## Usage

```bash
$ waluigi --help
Usage of waluigi:
  -controller string
    	Filter logs by the 'controller' field
  -name string
    	Filter logs by the 'name' field
  -namespace string
    	Filter logs by the 'namespace' field
```

### Example

```bash
kubectl logs -f -n capi-system deploy/capi-controller-manager | waluigi
```



### Some suggestions for your README

After you have created your new repository, you may want to add some of these badges to the top of your README.

- **CircleCI:** After enabling builds for this repo via [this link](https://circleci.com/setup-project/gh/giantswarm/REPOSITORY_NAME), you can find badge code on [this page](https://app.circleci.com/settings/project/github/giantswarm/REPOSITORY_NAME/status-badges).
