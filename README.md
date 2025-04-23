# waluigi

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/giantswarm/waluigi/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/giantswarm/waluigi/tree/main)
[![Go Reference](https://pkg.go.dev/badge/github.com/giantswarm/waluigi.svg)](https://pkg.go.dev/github.com/giantswarm/waluigi)


CLI tool to pretty print the logs from CAPI controllers. Companion to [https://github.com/giantswarm/luigi](https://github.com/giantswarm/luigi)

It was vibe coded using ChatGTP. Feel free to burn the code, change it or rewrite it.

## Installation
```bash
go install github.com/giantswarm/waluigi
```

## Usage

```bash
$ waluigi --help
Usage of waluigi:
  --controller string
    	Filter logs by the 'controller' field
  --name string
    	Filter logs by the 'name' field
  --namespace string
    	Filter logs by the 'namespace' field
```

### Example

```bash
kubectl logs -f -n capi-system deploy/capi-controller-manager | waluigi
```

![image](https://github.com/user-attachments/assets/7694aa10-b6ca-49d5-8b0b-bcdd60a27495)
