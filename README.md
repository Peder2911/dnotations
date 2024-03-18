
# dnotations 

A simple tool for organizing your SystemD units with annotations via comment headers.

## Usage

First, add a comment-header to the SystemD unit file that you want to keep track of metadata for, like so:

```
# annotations:
#   part-of: dnotationstest 
#   component: absolute-unit 
#   version: 0.0.0
#   managed-by: peder
#

[Unit]
Description="My absolute unit"
...
```

Units with such annotations can be queried with the 'list' command:

```
>>> dens list | python3 -m json.tool 
{
    "units": [
        {
            "part-of": "dnotationstest",
            "component": "absolute-unit",
            "version": "0.0.0",
            "managed-by": "peder",
            "hostname": ""
        }
    ]
}
```

Use this to keep track of what your units are for, and what they are doing. Makes managing multiple interrelated units a bit easier!
