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
	"context"
	"encoding/json"
	"github.com/peder2911/dnotations/internal/cli"
	"github.com/peder2911/dnotations/pkg/io"
)

func main() {
	ctx := context.Background()
	s, err := io.NewDnotatedSystemd()
	if err != nil {
		panic(err)
	}

	units, err := s.ListUnits(ctx)
	output_model := cli.DnotatedUnitsListing{Units: *units}
	output, err := json.Marshal(output_model)
	if err != nil {
		panic(err)
	}
	print(string(output))
}
