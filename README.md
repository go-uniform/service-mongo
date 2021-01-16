# service-mongo
A templated starting point for uniform microservices

### Getting Started
First step is to compile resources and metadata into the project's source-code:
```
go generate
```
This will create a file `service/meta.go` which is ignored by the `.gitignore` and contains the project's resources and metadata.

### Future Plans
Instead of `.description` file, use github curl request as primary method to get service description:
```
curl https://api.github.com/repos/go-uniform/base-api
``` 
keep `.description` file as a fail-over for project not on source-control.

### CLI Commands

cmd command example `cmd/command.example-one.go`:
```
package cmd

import (
	"github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
	"go-uniform/service-mongo/service"
)

var exampleOneCmd = &cobra.Command{
	Use:   "command:example-one",
	Short: "Request the running " + AppName + " to execute the example-one command",
	Long:  "Request the running " + AppName + " to execute the example-one command",
	Run: func(cmd *cobra.Command, args []string) {
		service.Command("example-one", natsUri)
	},
}

func init() {
	exampleOneCmd.Flags().StringVarP(&natsUri, "nats", "n", nats.DefaultURL, "The nats cluster URI")
	rootCmd.AddCommand(exampleOneCmd)
}
```

service command example `service/command.example-one.go`:
```
package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local(command("example-one")), exampleOne)
}

func exampleOne(r uniform.IRequest, p diary.IPage) {
	// todo: write logic here
}
```

### Background Worker Process

Use a CLI Command and add it to a scheduled cronjob to avoid the background process from being executed multiple times when scaling service instances.
In other words this will work like a sync.Mutex but across all running instances of the given service, allowing us to add as many service instances as we need.

### Routines

service action example `service/routine.example-two.go`:
```
package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("example-two"), exampleTwo)
}

func exampleTwo(r uniform.IRequest, p diary.IPage) {
	// todo: write logic here
}
```

### Events

service event example `service/event.example-three.go`:
```
package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("example-three"), exampleThree)
}

func exampleThree(r uniform.IRequest, p diary.IPage) {
	// todo: write logic here
}
```