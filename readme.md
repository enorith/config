# Config component for [Enorith](https://github.com/enorith/framework)

## Basic usage

```go
package main

import (
    "github.com/enorith/config"
)
type FooConfig struct {
    // load yaml config, fallback environment variable
    Foo string `env:"ENV_FOO"`
}



func main() {
    var c FooConfig
	config.Unmarshal("config.yml", &c)
    
}
```