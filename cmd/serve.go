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
	"fmt"
	"net/http"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/thefirstspine/boiler2/config"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	path = "/webhooks"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve for a webhook",
	Long:  `Serve for a webhook`,
	Run: func(cmd *cobra.Command, args []string) {
		hook, _ := github.New(github.Options.Secret("MyGitHubSuperSecretSecrect...?"))

		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
			if err != nil {
				if err == github.ErrEventNotFound {
					// ok event wasn;t one of the ones asked to be parsed
				}
			}
			switch payload.(type) {

			case github.ReleasePayload:
				// Getting release struct
				release := payload.(github.ReleasePayload)
				// Getting project struct
				project, projectStatus := config.GetConfig(release.Repository.Name)
				if !projectStatus {
					return
				}
				// Deploy
				exec.Command(
					"/bin/bash",
					"-c",
					fmt.Sprintf("./boiler2 deploy %s --tag_or_branch=%s", project.Name, release.Release.TagName)).Output()
			}
		})
		http.ListenAndServe(":3000", nil)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
