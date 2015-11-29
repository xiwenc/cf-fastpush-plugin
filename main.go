package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/parnurzeal/gorequest"
	"github.com/xiwenc/cf-fastpush-controller/lib"
	"github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/cf/terminal"
	"github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/plugin"
	"github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/simonleung8/flags"
	"strings"
	"encoding/json"
	"io/ioutil"
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

	cliLogged, err := cliConnection.IsLoggedIn()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	if cliLogged == false {
		panic("cannot perform fast-push without being logged in to CF")
	}


	if args[0] == "fast-push" || args[0] == "fp" {
		if len(args) == 1 {
			c.showUsage(args)
			return
		}
		// set flag for dry run
		fc := flags.New()
		fc.NewBoolFlag("dry", "d", "bool dry run flag")

		err := fc.Parse(args[1:]...)
		if err != nil {
			c.ui.Failed(err.Error())
		}
		// check if the user asked for a dry run or not
		if fc.IsSet("dry") {
			dryRun = fc.Bool("dry")
		} else {
			c.ui.Warn("warning: dry run not set, commencing fast push")
		}

		c.ui.Say("Running the fast-push command")
		c.ui.Say("Target app: %s \n", args[1])
		c.FastPush(cliConnection, args[1], dryRun)
	} else if args[0] == "fast-push-status" || args[0] == "fps" {
		c.FastPushStatus(cliConnection, args[1])
	} else {
		return
	}

}

func (c *FastPushPlugin) FastPushStatus(cliConnection plugin.CliConnection, appName string) {
	apiEndpoint := c.GetApiEndpoint(cliConnection, appName)
	status := lib.Status{}
	request := gorequest.New()
	_, body, err := request.Get(apiEndpoint + "/status").End()
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(body), &status)
	c.ui.Say(status.Health)
}

func (c *FastPushPlugin) FastPush(cliConnection plugin.CliConnection, appName string, dryRun bool) {
	// Please check what GetApp returns here
	// https://github.com/cloudfoundry/cli/blob/master/plugin/models/get_app.go

	if dryRun {
		// NEED TO HANDLE DRY RUN
		c.ui.Warn("warning: No changes will be applied, this is a dry run !!")
	}

	apiEndpoint := c.GetApiEndpoint(cliConnection, appName)
	request := gorequest.New()
	_, body, err := request.Get(apiEndpoint + "/files").End()
	if err != nil {
		panic(err)
	}
	remoteFiles := map[string]*lib.FileEntry{}
	json.Unmarshal([]byte(body), &remoteFiles)

	localFiles := lib.ListFiles()

	filesToUpload := c.ComputeFilesToUpload(localFiles, remoteFiles)
	payload, _ := json.Marshal(filesToUpload)
	_, body, err = request.Put(apiEndpoint + "/files").Send(string(payload)).End()
	if err != nil {
		panic(err)
	}
	status := lib.Status{}
	json.Unmarshal([]byte(body), &status)
	c.ui.Say(status.Health)
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
					Usage: "cf fast-push APP_NAME\n   cf fp APP_NAME",
					Options: map[string]string{
						"dry": "--dry, dry run for fast-push",
					},
				},
			},
			plugin.Command{
				Name:     "fast-push-status",
				Alias:    "fps",
				HelpText: "fast-push-status shows the current state of your application",
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

func (c *FastPushPlugin) GetApiEndpoint(cliConnection plugin.CliConnection, appName string) string {
	results, err := cliConnection.CliCommandWithoutTerminalOutput("app", appName)
	if err != nil {
		c.ui.Failed(err.Error())
	}


	for _, line := range results {
		match, _ :=regexp.MatchString("^urls:.*", line)
		if match {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return "https://" + parts[1] + "/_fastpush"
			}
		}
	}
	panic("Could not find usable route for this app. Make sure at least one route is mapped to this app")
}

func (c *FastPushPlugin) ComputeFilesToUpload(local map[string]*lib.FileEntry, remote map[string]*lib.FileEntry) map[string]*lib.FileEntry {
	filesToUpload := map[string]*lib.FileEntry{}
	for path, f := range local {
		if remote[path] == nil {
			c.ui.Say("[NEW] " + path)
			f.Content, _ = ioutil.ReadFile(path)
			filesToUpload[path] = f
		} else if remote[path].Checksum != f.Checksum {
			c.ui.Say("[MOD] " + path)
			f.Content, _ = ioutil.ReadFile(path)
			filesToUpload[path] = f
		}
	}
	return filesToUpload
}
