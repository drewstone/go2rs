# go2rs
[![PkgGoDev](https://pkg.go.dev/badge/drewstone/go2rs)](https://pkg.go.dev/github.com/drewstone/go2rs)

## What is go2rs?
- go2rs is a Rust struct generator from Go structs
- Automatically recognizes Go modules in the directory and generates equivalent Rust types
- Handles Go-to-Rust type conversions with appropriate derives and attributes

## Installation
```console
$ go get github.com/drewstone/go2rs
```

## Usage

```go
// ./example/main.go
package main

import (
    "time"
)

type Status string

const (
    StatusOK Status = "OK"
    StatusFailure Status = "Failure"
)

type Param struct {
    Status    Status
    Version   int
    Action    string
    CreatedAt time.Time
}
```

```console
$ go2rs ./example
```

Generates:

```rust
use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub enum Status {
    OK,
    Failure,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct Param {
    pub status: Status,
    pub version: i32,
    pub action: String,
    #[serde(with = "chrono::serde::ts_seconds")]
    pub created_at: DateTime<Utc>,
}
```

## Features
- Converts Go types to idiomatic Rust types
- Handles common Go patterns like string enums
- Adds appropriate serde derives and attributes
- Supports time.Time conversion to chrono::DateTime
- Maintains field visibility and naming conventions
- Generates documentation from Go comments

## TODO
- Handle custom MarshalJSON/UnmarshalJSON implementations

## Acknowledgements
This is entirely built using [go2ts](https://github.com/go-generalize/go2ts) by [go-generalize](https://github.com/go-generalize) as a reference and porting over the same concepts to Rust.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.