# Infra - Step 2 (Cloud Run Service)

Use KPT to setup Cloud Runa and deploy the service

## Setup
*NOTE:* The `setters.yaml` is a symbolic link to the `setters.yaml` file in the parent `infra` folder.

The `app-path` varible is used here to set the full path to build docker file in artifact registry

## Deploy

Once the variables are setup in the `setters.yaml` file, run:

`kpt fn render`\
`kpt live init`\
`kpt live apply`