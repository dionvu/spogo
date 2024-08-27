package cli

import (
	"fmt"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/icons"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// Prints the command to print help information
// corresponding to the command that the user messed up on.
func handleBadUsage(ctx *cli.Context, err error) {
	fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), err)
	fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+ctx.App.Name+" help "+ctx.Command.Name))
}

// Prints the error message corresponding to
// command requiring flags but none were provided.
func handleNoFlag(ctx *cli.Context) {
	fmt.Printf("%v %v\n", color.RedString(icons.Warning+"Error:"), "no flags provided")
	fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" help "+ctx.Command.Name))
}

func printNoDeviceHelp() {
	fmt.Printf("%v\n", color.YellowString(icons.Question+"Help: "+config.APPNAME+" device"))
}
