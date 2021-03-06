# GoFigure

GoFigure is a small utility library for reading configuration files. It's usefuly especially if you want to load many files recursively (think `/etc/apache2/mods-enabled/*.conf` on Ubuntu).

It can support multiple formats, as long as you take a file and unmarshal it into a struct containing your configurations. 

Right now the only implemented formats are YAML and JSON files, but feel free to add more :)

## Example usage:

```go 

import (
	"fmt"
	"github.com/EverythingMe/gofigure"
	"github.com/EverythingMe/gofigure/yaml"
)

// this is our configuration container
var conf = &struct {
	Redis struct {
		Server  string
		Monitor int
		Timeout int
	}
}{}

func main() {
	
	//if we set some default, the loader will override it
	conf.Redis.Server = "localhost:6377"

	// init a loader with a YAML decoder in strict mode
	loader := gofigure.NewLoader(yaml.Decoder{}, true)

	// run recursively on some directory
	err := loader.LoadRecursive(conf, "/etc/myservice/conf.d")
	if err != nil {
		panic(err)
	}

	fmt.Println(conf.Redis.Server)
	//Output: localhost:6379
}

```

## Automatic -conf and -confdir flags

GoFigure can automatically add the optional `-conf ` and `-confdir` flags to your program's command line flags, and then
automatically read config files specified by them.

To do that, simply add an import to `"github.com/EverythingMe/gofigure/autoflag"`.

You can then use `autoflag.Load` to load either a single file or all files in the directory specified by these flags.

Modifying the above example to do this:


```go

import (
    // ... other imports ... 
    
	"github.com/EverythingMe/gofigure"
    "github.com/EverythingMe/gofigure/autoflag"
)

func main() {
	
	// in this example we also use the default loader which loads yaml files,
    // but any loader can work here
    err := autoflag.Load(gofigure.DefaultLoader, &conf)
    if err != nil {
		panic(err)
	}
    
}
```

## Reloading configurations on the fly

GoFigure provides a primitive utility for waiting on config reloads. Right now the only implemented method
is to send a SIGHUP to the process, and have the process using GoFigure use a `ReloadMonitor`, with a `Reloader`
that gets called when that signal is sent.

We do not deal with the actual loading, as each program has its own sensitivites to what parts can be reconfigured in 
runtime and which can't.

### Example:

```go

import (
    // ... other imports ... 
    
	"github.com/EverythingMe/gofigure"
    "github.com/EverythingMe/gofigure/autoflag"
)


// loadConfigs is used both to initially load the configs, and to refresh them
func loadConfigs() {
	if err := autoflag.Load(gofigure.DefaultLoader, &conf); err != nil {
        log.Println(err)
    }
}


func main() {
	
	loadConfigs()
    
    // Create a SignalMonitor that calls a reloader when SIGHUP is sent to our process
    m := gofigure.NewSignalMonitor()
    
    // make the monitor wait for siguhup and call configReload when it needs to
	m.Monitor(gofigure.ReloadFunc(loadConfigs))
    
}

```