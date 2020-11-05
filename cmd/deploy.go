/*
Copyright © 2020 The First Spine

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thefirstspine/boiler2/commands"
	"github.com/thefirstspine/boiler2/config"
	"github.com/thefirstspine/boiler2/nginx"
	"github.com/thefirstspine/boiler2/ports"
)

func init() {
	c := color.New(color.FgMagenta)
	c.Println("██████╗      ██████╗     ██╗    ██╗         ███████╗    ██████╗ ")
	c.Println("██╔══██╗    ██╔═══██╗    ██║    ██║         ██╔════╝    ██╔══██╗")
	c.Println("██████╔╝    ██║   ██║    ██║    ██║         █████╗      ██████╔╝")
	c.Println("██╔══██╗    ██║   ██║    ██║    ██║         ██╔══╝      ██╔══██╗")
	c.Println("██████╔╝    ╚██████╔╝    ██║    ███████╗    ███████╗    ██║  ██║")
	c.Println("╚═════╝      ╚═════╝     ╚═╝    ╚══════╝    ╚══════╝    ╚═╝  ╚═╝")
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy {app}",
	Short: "Deploy an app.",
	Long:  `Deploy an app.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an `app` argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get project
		project, status := config.GetConfig(args[0])
		if !status {
			color.Red("Cannot find app `%s`", args[0])
			panic("Something goes wrong in the deployment. Old container still up & running.")
		}

		// Get common config
		common, status := config.GetCommon()

		// Check requirements
		color.Cyan("\nChecking requirements...")
		if !CheckRequirements() {
			color.Red("Requirements checkings failed.")
			color.Red("You can install required dependencies by calling `./setup.sh`")
			panic("Something goes wrong in the deployment. Old container still up & running.")
		}
		color.Green("Done")

		// Get the app from the git repository
		directory := fmt.Sprintf("boilerapp_%s", project.Name)
		color.Cyan(fmt.Sprintf("\nGetting project from repo %s...", project.Repository))
		if !commands.GitClone(project.Repository, directory) {
			color.Red("Failed to clone the repository.")
			color.Red("Please check you have to rights to perform this action.")
			panic("Something goes wrong in the deployment. Old container still up & running.")
		}
		color.Green("Done")

		// Build image
		imageName := fmt.Sprintf("boilerimage_%s", project.Name)
		color.Cyan(fmt.Sprintf("\nBuild image %s from directory %s...", imageName, directory))
		if !commands.DockerBuild(imageName, directory) {
			color.Red("Failed to build the image.")
			color.Red("Please ensure that the Docker daemon is running on this machine.")
			panic("Something goes wrong in the deployment. Old container still up & running.")
		}
		color.Green("Done")

		// Stop old container
		containerName := fmt.Sprintf("boilercontainer_%s", project.Name)
		commands.DockerStop(containerName)

		// Run the new container
		port := ports.GetFirstFreePort(1024, 49151, []int{1024, 1433, 1521, 3306, 5432})
		color.Cyan(fmt.Sprintf("\nRun image %s to container %s:%d...", imageName, containerName, port))
		fmt.Print(common.Env)
		if !commands.DockerRun(
			imageName,
			containerName,
			append(project.Env, common.Env...),
			fmt.Sprintf("%d:%d", 8080, port),
		) {
			color.Red("Failed to build the image.")
			color.Red("Please ensure that the Docker daemon is running on this machine.")
			panic("Something goes wrong in the deployment. Old container might be not up & running.")
		}
		color.Green("Done")

		// Write nginx config
		domain := project.Domain
		forward := fmt.Sprintf("127.0.0.1:%d", port)
		color.Cyan(fmt.Sprintf("\nGenerate nginx config forwarding %s to %s...", domain, forward))
		if !nginx.WriteConfig(domain, nginx.GenerateConfig(domain, forward)) {
			color.Red("Failed to write nginx config.")
			color.Red("Ensure that Boiler has access write to `/etc/nginx/sites-*`")
			panic("Something goes wrong in the deployment. Old container might be not up & running.")
		}
		color.Green("Done")

		// Call certbot
		color.Cyan(fmt.Sprintf("\nGenerate certificate for domain %s...", domain))
		if !commands.Certbot(domain) {
			color.Red("Failed to generate a certificate.")
			color.Red("This error is usually a problem with the Letsencrypt challenge. More infos here: https://certbot.eff.org/docs/challenges.html")
			panic("Something goes wrong in the deployment. Old container might be not up & running.")
		}
		color.Green("Done")

		// Clean
		if !commands.TestCommand(fmt.Sprintf("rm -rf %s", directory)) {
			color.Yellow("Failed to remove the project directory.")
			color.Yellow("You should rm it by yourself.")
		} else {
			color.Green("Done")
		}

		// All done!
		color.Green("\nDeployment done!")
	},
}

// Will check the requirements on the machin Boiler is launched.
// Required packages are: git, docker, nginx & certbot.
// To install these dependencies, you should launch the `setup.sh`
// script that will install everything for you.
func CheckRequirements() bool {
	// Git is required
	if !commands.TestCommand("git --version") {
		color.Red("Git not installed")
		return false
	}

	// Docker is required
	if !commands.TestCommand("docker --version") {
		color.Red("Docker not installed")
		return false
	}

	// nginx is required
	if !commands.TestCommand("nginx -v") {
		color.Red("nginx not installed")
		return false
	}

	// certbot is required
	if !commands.TestCommand("certbot --version") {
		color.Red("certbot not installed")
		return false
	}

	return true
}
