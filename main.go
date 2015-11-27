package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/simonleung8/flags"
)

/*
*	This is the struct implementing the interface defined by the core CLI. It can
*	be found at  "github.com/cloudfoundry/cli/plugin/plugin.go"
*
 */
type FastPushPlugin struct {
	ui terminal.UI
}

/*
*	This function must be implemented by any plugin because it is part of the
*	plugin interface defined by the core CLI.
*
*	Run(....) is the entry point when the core CLI is invoking a command defined
*	by a plugin. The first parameter, plugin.CliConnection, is a struct that can
*	be used to invoke cli commands. The second paramter, args, is a slice of
*	strings. args[0] will be the name of the command, and will be followed by
*	any additional arguments a cli user typed in.
*
*	Any error handling should be handled with the plugin itself (this means printing
*	user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
*	1 should the plugin exits nonzero.
 */
func (c *FastPushPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that the user called the command fast-push
	// alias fp is auto mapped
	var dryRun bool
	c.ui = terminal.NewUI(os.Stdin, terminal.NewTeePrinter())

	if args[0] == "fast-push" || args[0] == "fp" {
		// set flag for dry run
		fc := flags.New()
		fc.NewBoolFlag("dry", "d", "bool dry run flag")

		err := fc.Parse(args[1:]...)
		if err != nil {
			c.ui.Failed(err.Error())
		}
		if fc.IsSet("dry") {
			dryRun = fc.Bool("dry")
		}

	} else {
		return
	}

	cliLogged, err := cliConnection.IsLoggedIn()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	if cliLogged == false {
		panic("cannot perform fast-push without being logged in to CF")
	}

	if len(args) > 2 {
		fmt.Println("Running the fast-push command")
		fmt.Printf("Target app: %s /n", args[1])
		// check if the user asked for a dry run or not
		if dryRun {
			c.fastPush(cliConnection, args[1], true)
		} else {
			c.fastPush(cliConnection, args[1], false)
		}
	} else {
		c.showUsage(args)
	}

}

func (c *FastPushPlugin) fastPush(cliConnection plugin.CliConnection, appName string, dryRun bool) {
	// Please check what GetApp returns here
	// https://github.com/cloudfoundry/cli/blob/master/plugin/models/get_app.go

	if dryRun {
		c.ui.Warn("warning: No changes will be applied, this is a dry run !!")
	}

	app, err := cliConnection.GetApp(appName)
	if err != nil {
		c.ui.Failed(err.Error())
	}
	routes := app.Routes

	if len(routes) > 1 {
		for _, route := range routes {
			c.ui.Warn("multiple corresponding url's has been found")
			c.ui.Say(route.Host)
			c.ui.Say(route.Domain.Name)
		}
	} else {
		c.ui.Say("corresponding app host: %s", routes[0].Host)
		c.ui.Say("corresponding app domain: %s", routes[0].Domain.Name)
	}
	panic("NOT IMPLEMENTED YET!")
	// dispatch request TODO
	url := routes[0].Host
	var query = []byte(`query-here`)
	req, err := http.NewRequest("POST", url+"/fast-push", bytes.NewBuffer(query))
	req.Header.Set("X-Custom-Header", "somevalue")
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

/*
*	This function must be implemented as part of the plugin interface
*	defined by the core CLI.
*
*	GetMetadata() returns a PluginMetadata struct. The first field, Name,
*	determines the name of the plugin which should generally be without spaces.
*	If there are spaces in the name a user will need to properly quote the name
*	during uninstall otherwise the name will be treated as seperate arguments.
*	The second value is a slice of Command structs. Our slice only contains one
*	Command Struct, but could contain any number of them. The first field Name
*	defines the command `cf basic-plugin-command` once installed into the CLI. The
*	second field, HelpText, is used by the core CLI to display help information
*	to the user in the core commands `cf help`, `cf`, or `cf -h`.
 */
func (c *FastPushPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "FastPushPlugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "fast-push",
				Alias:    "fp",
				HelpText: "fast-push removes the need to deploy your app again for a small change",
				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "fast-push appname\n   cf fp appname",
					Options: map[string]string{
						"dry": "--dry, dry run for fast-push",
					},
				},
			},
		},
	}
}

/*
* Unlike most Go programs, the `Main()` function will not be used to run all of the
* commands provided in your plugin. Main will be used to initialize the plugin
* process, as well as any dependencies you might require for your
* plugin.
 */
func main() {
	plugin.Start(new(FastPushPlugin))
}

func (c *FastPushPlugin) showUsage(args []string) {
	for _, cmd := range c.GetMetadata().Commands {
		if cmd.Name == args[0] {
			fmt.Println("Invalid Usage: \n", cmd.UsageDetails.Usage)
		}
	}
}
