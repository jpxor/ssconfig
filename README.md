[![Go](https://github.com/jpxor/ssconfig/actions/workflows/go.yml/badge.svg)](https://github.com/jpxor/ssconfig/actions/workflows/go.yml)
[![CodeQL](https://github.com/jpxor/ssconfig/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/jpxor/ssconfig/actions/workflows/codeql-analysis.yml)

# ssconfig
Super Simple Config:  read from json file, overwrite with ENV vars. no extra dependencies.

 - values are read directly into your config struct, 
 - field names of the struct match the ENV vars (with optional prefix)

```go
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
ssconfig supports loading all types from file via "encoding/json" unmarshal. You can use the json struct tags.

ssconfig supports loading these types from environment variables:
 - Int8/Int16/Int32/Int64/Int
 - Float32/Float64,
 - String

ssconfig logs to default logger. It names the Fields loaded, but not their values (keeps any private api keys or passwords out of the logs).
