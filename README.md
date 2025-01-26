# AutoGit 🚀

AutoGit is a GitHub automation script designed to simplify and automate various tasks related to GitHub repositories. It provides a CLI tool to manage deployment environments, secrets, variables, repository setup, and more.

## Features ✨

- 📄 Pass YAML file for configuration
- 📊 Output latest deployment environment information
- 🛠️ CLI tool for easy interaction
- 🔒 Create and update deployment environments, secrets, variables, and repository settings
- 🚀 Trigger and approve deployments
- ⚙️ Customizable output with ARGS

## Installation 🛠️

To build and run the AutoGit CLI tool, use the following commands:

```sh
go build -o ./bin/autogit && ./bin/autogit fetch config.yaml
```

To run the tool without building:
```sh
go run main.go fetch config.yaml
```

## Usage 📖
Fetch Deployment Information 📊
To fetch the latest deployment environment information:
```sh
./bin/autogit fetch config.yaml
```
## Create Resources 🛠️
To create resources such as repositories, deployment environments, secrets, and variables:
```sh
./bin/autogit create --type <resource-type> config.yaml
```
Replace <resource-type> with one of the following options:

- `deployment-env`
- `repository`
- `secrets`
- `variables`
- `secrets-variables`

## Configuration ⚙️
The configuration file (`config.yaml`) should be structured as follows:
```YAML
org: ORG-NAME
repos:
	Repo1:
		Dev:
			createDeploymentEnv: true # Ceates the deployment environments 
			fetchReleases: false # Fetches the latest release information on the environment -- Not included in ALL command
			createVariables: true # Create Variables within the repo and deployment environment
			createSecrets: true # Create Secrets within the repo and deployment environment
			variables:
					var1: "I AM A VARIABLE"
			secrets:
					secret1: "I AM A SECRET"
	Repo2:
		Dev:
			createDeploymentEnv: true # Ceates the deployment environments 
			fetchReleases: false # Fetches the latest release information on the environment -- Not included in ALL command
			createVariables: true # Create Variables within the repo and deployment environment
			createSecrets: true # Create Secrets within the repo and deployment environment
			variables:
				var1: "I AM A VARIABLE"
			secrets:
				secret1: "I AM A SECRET"
```


## Acknowledgments 🙏
Hat tip to anyone whose code was used
Inspiration:
    - Having loads of Secrets and Variables is annoying to update :smile