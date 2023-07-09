package main

import (
	"fmt"
	"os"
	"os/exec"

	crongo "github.com/abhijit-hota/crongo"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: cron \"<cron expression>\" \"<command>\"")
		os.Exit(1)
	}

	cronExpr := os.Args[1]
	toRun := os.Args[2]
	task := exec.Command("sh", "-c", toRun)

	err := crongo.RunCronJob(cronExpr, func() {
		err := task.Run()
		if err != nil {
			panic(err)
		}
	})

	if err != nil {
		panic(err)
	}
}
