package main

import (
    //"./dns"
    //"bufio"
    //"github.com/public0821/dnserver/db"
    //"github.com/public0821/dnserver/dns"
    //"github.com/public0821/dnserver/util"
    //"github.com/public0821/dnserver/web"
    "fmt"
    "github.com/codegangsta/cli"
    "os"
    //"io"
    //"log"
)

func main() {
    app := cli.NewApp()
    app.Name = "dnscli"
    app.Usage = "command line interface for dnserver"
    app.Action = func(c *cli.Context) {
        fmt.Println("Greetings")
    }
    app.Commands = []cli.Command{
        {
            Name:      "add",
            ShortName: "a",
            Usage:     "add a task to the list",
            Action: func(c *cli.Context) {
                println("added task: ", c.Args().First())
            },
            Flags: []cli.Flag{cli.StringFlag{"lang, l", "english", "language for the greeting"}},
        },
        {
            Name:      "complete",
            ShortName: "c",
            Usage:     "complete a task on the list",
            Action: func(c *cli.Context) {
                println("completed task: ", c.Args().First())
            },
        },
    }
    app.Run(os.Args)
}
