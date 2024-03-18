# Validator Duties

This is a CLI meant for managing the active and upcoming duties of an ethereum validator.
Currently, it allows for checking the gaps between upcoming attestations. 

## Installation

To install:
```shell
go install github.com/Alec1017/validator-duties@latest
```

Check if it is installed:
```shell
validator-duties --help
```

## Usage

To use, run: 

```shell
validator-duties [global options]
```

### Flags

- `--validator`: Index of the validator to get duties for.
- `--timezone`: Timezone for displaying timestamps. (default: "UTC")
- `--beacon-node-endpoint`: Endpoint URL for the beacon node API. (default: "http://localhost:5052")
