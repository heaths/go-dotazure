# Dotazure

[![releases](https://img.shields.io/github/v/release/heaths/go-dotazure.svg?logo=github)](https://github.com/heaths/go-dotazure/releases/latest)
[![reference](https://pkg.go.dev/badge/github.com/heaths/go-dotazure.svg)](https://pkg.go.dev/github.com/heaths/go-dotazure)
[![ci](https://github.com/heaths/go-dotazure/actions/workflows/ci.yml/badge.svg?event=push)](https://github.com/heaths/go-dotazure/actions/workflows/ci.yml)

Locate and load environment variables defined when provisioning an [Azure Developer CLI] project.

## Getting Started

If you do not already have an [Azure Developer CLI] (azd) project, you can create one:

```sh
azd init
```

After you define some resources e.g., an [Azure Key Vault](https://github.com/heaths/go-dotazure/blob/main/infra/resources.bicep),
you can provision those resources which will create a `.env` file with any `output` parameters:

```sh
azd up
```

## Example

After `azd up` provisions resources and creates a `.env` file, you can call `Load()` to load those environment variables
from the default environment e.g.,

```go
package main

import (
    "errors"
    "fmt"
    "os"

    "github.com/heaths/go-dotazure"
)

func main() {
    if loaded, err := dotazure.Load(); err != nil {
        panic(err)
    } else if loaded {
        fmt.Fprintln(os.Stderr, "loaded environment variables")
    }

    // Assumes bicep contains e.g.
    //
    // output AZURE_KEYVAULT_URL string = kv.properties.vaultUri
    vaultURL, _ := os.LookupEnv("AZURE_KEYVAULT_URL")
    fmt.Printf("AZURE_KEYVAULT_URL=%q\n", vaultURL)
}
```

If you want to customize behavior, you can call `Load()` with various `With*` option functions.

## License

Licensed under the [MIT](https://github.com/heaths/go-dotazure/blob/refactor/LICENSE.txt) license.

[Azure Developer CLI]: https://aka.ms/azd
