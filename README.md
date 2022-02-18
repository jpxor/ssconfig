[![Go](https://github.com/jpxor/ssconfig/actions/workflows/go.yml/badge.svg)](https://github.com/jpxor/ssconfig/actions/workflows/go.yml)
[![CodeQL](https://github.com/jpxor/ssconfig/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/jpxor/ssconfig/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jpxor/ssconfig)](https://goreportcard.com/report/github.com/jpxor/ssconfig)

# ssconfig
Super Simple Config:  read from json file, overwrite with ENV vars.

 - values are read directly into your config struct, 
 - field names of the struct match the ENV vars (with optional prefix)
 - single source file
 - full test coverage
 - no dependencies

## Example Usage
```sh
go get github.com/jpxor/ssconfig
```

```go
package main

import (
    "github.com/jpxor/ssconfig"
)

type my_config struct {
  PORT        string
  WWW_ROOT    string
  PRIVATE_KEY string
}

func main() {
  var env my_config
  
  // default: load from environment variables matching struct field names
  ssconfig.Load(&env)
 
  // or set optional parameters to:
  err := ssconfig.Set{
    FilePath: "config.json",    // load from config file before env vars,
    EnvPrefix: "MY_APP_",       // set env var name prefix,
    Logger: myLogger,           // use custom log.Logger
  }.Load(&env)

  // any error will contain a list of fields that failed to load
  if err != nil {
    log.Println(err)
  }
  log.Printf("%+v\n", env)
}
```
## Documentation
See the example above, that's 99% what you need to know.
The other 1%:
 - ssconfig supports loading all types via "encoding/json" unmarshal. You can use the json struct tags for loading from config file.
 - ssconfig logs to default logger if one is not set. It names the Fields loaded, but not their values (keeps any private api keys or passwords out of the logs).
 - need default values? Then set them before loading; the load will overwrite the existing values.

## About
I tried a few different config packages from a 'curated list' and they all had lots of dependencies and features I didn't need, or returned errors I wasn't expecting (config was missing from config file, so it gave up instead of looking in the env vars). So I quickly whipped up a config package that does one simple thing, has no dependencies, and doesn't immediately give up when it can't find the config value for one of the fields.
