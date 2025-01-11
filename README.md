# AutoGit
Another attempt at a Github automation script 

### TODO: 

- [x] Pass YAML file through for config
- [x] Outputs latest Deployment environment information
- [x] Converted to CLI (Moved old code to ./.old and new `autogit` cli is under ./cli)
- [x] Add output to YAML file.
- [ ] Test multiple repos
    - [ ] Add a create command
        - [ ] Repos
        - [ ] Deployment Environments
        - [ ] Secrets & Variables
        - [ ] Repo setup/policies/branches
    - [ ] Add a Update command
        - [ ] Repos
        - [ ] Deployment Environments
        - [ ] Secrets & Variables
        - [ ] Repo setup/policies/branches
- [ ] Add Background colour to the status line
- [ ] Trigger deployments
- [ ] Approval of deployments
- [ ] Add conditions to add ARGS for specific output.
### NOTES 
To build and run
`go build -o ./bin/autogit && ./bin/autogit fetch config.yaml`

To just run 
`go run main.go fetch config.yaml`

- use github api to automate simple/complext tasks 
- use mongodb for local storage `mongod --dbpath ~/mongodb-data`