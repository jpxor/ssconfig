[![Go](https://github.com/jpxor/ssconfig/actions/workflows/go.yml/badge.svg)](https://github.com/jpxor/ssconfig/actions/workflows/go.yml)
[![CodeQL](https://github.com/jpxor/ssconfig/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/jpxor/ssconfig/actions/workflows/codeql-analysis.yml)

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
  
  // load from environment variables only, no prefix
  ssconfig.Load(&env)
 
  // or load from file, and then from env vars with prefix
  err := ssconfig.Set{
    FilePath: "config.json",
    EnvPrefix: "MY_APP_",
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
 - ssconfig supports loading all types from file via "encoding/json" unmarshal. You can use the json struct tags.
 - ssconfig supports loading these types from environment variables:
    - Int8/Int16/Int32/Int64/Int
    - Float32/Float64,
    - String
 - ssconfig logs to default logger. It names the Fields loaded, but not their values (keeps any private api keys or passwords out of the logs).
 - need default values? Then set them before loading, the load will overwrite them.

## About
I tried a few different config packages from a 'curated list' and they all had lots of dependencies and features I didn't need, or returned errors I wasn't expecting (config was missing from config file, so it gave up instead of looking in the env vars). So I quickly whipped up a config package that does one simple thing, has no dependencies, and doesn't immediately give up when it can't find the config value for one of the fields.
