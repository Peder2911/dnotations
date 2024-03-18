/*
	dnotations

This is a (wip) tool for people who would like to organize their Systemd units with metadata and annotations.
This might be helpful to keep track of things like what services provide what functionality on a box, behind
which hostnames, and so on.

You do this by annotating your SystemD unit files with a comment header with some YAML data, like so:

```{example.service}
# annotations:
#   part-of: my-project
#   component: proxy
#   managed-by: someuser
#   version: 1.0.0
#   hostname: myservice.example.com
[Unit]
Description=A proxy server.
...
```
*/
package main

import (
	"log"
	"github.com/peder2911/dnotations/pkg/io"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	s,err := io.NewDnotatedSystemd()
	if err != nil {
		panic(err)
	}

	units,err := s.ListUnits(ctx)
	log.Printf("Loaded %v units\n",len(*units))
	for _,u := range *units {
		fmt.Println(u)
	}
	if err != nil {
		panic(err)
	}
}
