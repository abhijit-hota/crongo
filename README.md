### Crongo

Crongo is a simple cron job runner written in Go. It allows you to schedule jobs using a cron expression.

## Usage

### As a library

The library exposes a single function `RunCronJob` which takes in a cron expression and a function to be executed.

```go
package main

import (
    "fmt"
    "time"

    "abhijithota.me/crongo"
)

func main() {
    crongo.RunCronJob("*/5 * * * * *", func() {
        fmt.Println("Hello, World!")
    })
}
```

### As a CLI

The CLI exposes a single command `crongo` which takes in a cron expression and a command to be executed, both as strings.

```bash
$ crongo "*/5 * * * * *" 'echo "Hello, World!"'
```

## Future (?)

- [ ] Multiple jobs
- [ ] More cron expressions
- [ ] `crontab` file support 

## Inspiration

This is a weekend project, made because I just thought about implementing a Cron job runner PoC. Is nt tested properly. Use at your own risk. Or don't use it at all. Yeah, that's better. 