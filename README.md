[![CI](https://github.com/mgoeppe/grpc-gateway-csv/actions/workflows/ci.yml/badge.svg)](https://github.com/mgoeppe/grpc-gateway-csv/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mgoeppe/grpc-gateway-csv.svg)](https://pkg.go.dev/github.com/mgoeppe/grpc-gateway-csv)
[![Go Report Card](https://goreportcard.com/badge/github.com/mgoeppe/grpc-gateway-csv)](https://goreportcard.com/report/github.com/mgoeppe/grpc-gateway-csv)
[![License](https://img.shields.io/github/license/mgoeppe/grpc-gateway-csv)](LICENSE)

# grpc-gateway-csv

A CSV marshaler for [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).
Implements the [`runtime.Marshaler`](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/v2/runtime#Marshaler)
interface to render gRPC responses as CSV.

## Installation

```bash
go get github.com/mgoeppe/grpc-gateway-csv
```

## Usage

Register the marshaler with the grpc-gateway ServeMux:

```go
import csv "github.com/mgoeppe/grpc-gateway-csv"

mux := runtime.NewServeMux(
    runtime.WithMarshalerOption("text/csv", &csv.Marshaler{}),
)
```

### Options

- `RowDelim` — row delimiter (default: `\n`)
- `FieldDelim` — field delimiter (default: `;`)
- `InnerDelim` — delimiter for merged values within a field (default: `|`)
- `NoHeader` — suppress header row
- `Printf` — custom format function (e.g., for locale-aware number formatting)

Documentation on grpc-gateway custom marshalers: [Customizing Your Gateway](https://github.com/grpc-ecosystem/grpc-gateway/blob/main/docs/docs/mapping/customizing_your_gateway.md).
