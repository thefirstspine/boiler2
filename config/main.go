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
package config

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
)

// Main Projects JSON type. Will contain all the projects.
type Projects struct {
	// The projects entry
	Projects []Project `json:"projects"`
}

// Represents a project in Boiler. A project is a Docker container
// launched with environment variables & mapped to a domain in the
// nginx server of the host machine.
type Project struct {
	// Name is the project
	Name string `json:"name"`

	// The domain of the project
	Domain string `json:"domain"`

	// The repository of the project. Can be either SSH or HTTP.
	Repository string `json:"repository"`

	// All the environements variables passed to the
	// container. If you pass secrets in this configuration, you'd
	// better crypt them & store them in a secure place.
	Env []string `json:"env"`
}

func GetConfig(appname string) (project Project, ok bool) {

	// Open the jsonFile
	jsonFile, err := os.Open("./boiler.json")
	if err != nil {
		fmt.Println(err)
		var defVal Project
		return defVal, false
	}

	// Defer the closing of the file
	defer jsonFile.Close()

	// Read the file
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Unmarshal the byteArray to the projects variable
	var projects Projects
	json.Unmarshal(byteValue, &projects)

	// Iterate on the projects to find the right project
	for i := 0; i < len(projects.Projects); i++ {
		if projects.Projects[i].Name == appname {
			return projects.Projects[i], true
		}
	}

	// None project was found matching the name - return default
	// value with a wrong status.
	var defVal Project
	return defVal, false
}
