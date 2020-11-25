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
package commands

import (
	"fmt"
	"os/exec"
	"strings"
)

// Test a command - on error, will return `false`.
func TestCommand(cmd string, log bool) bool {
	if log {
		fmt.Println("=> test command", cmd)
	}
	output, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if log {
		fmt.Println("  => error", err)
	}
	if log {
		fmt.Println("  => output", output)
	}
	return err == nil
}

// Launch a TestCommand to clone the repository inside a `boilerapp_{project}`
// directory.
func GitClone(repository string, destination string, tagOfBranch string) bool {
	return TestCommand(fmt.Sprintf("git clone -b %s %s %s", tagOfBranch, repository, destination), true)
}

func DockerBuild(imageName string, directory string) bool {
	return TestCommand(fmt.Sprintf("docker build -t %s %s", imageName, directory), true)
}

func DockerStop(containerName string) bool {
	return TestCommand(fmt.Sprintf("docker stop %s", containerName), true)
}

func DockerRm(containerName string) bool {
	return TestCommand(fmt.Sprintf("docker rm %s", containerName), true)
}

func DockerRun(imageName string, containerName string, envVars []string, portForwarding string) bool {
	envVarsStr := strings.Join(envVars[:], "\" -e \"")
	return TestCommand(
		fmt.Sprintf("docker run --name %s -e \"%s\" -p %s -d %s", containerName, envVarsStr, portForwarding, imageName),
		true)
}

func Certbot(domain string) bool {
	return TestCommand(fmt.Sprintf("certbot --nginx --email teddy@coretizone.com -d %s --agree-tos --non-interactive", domain), true)
}
