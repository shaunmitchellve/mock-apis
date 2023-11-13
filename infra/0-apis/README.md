# Infra - Step 0 (Foundation)

Use KPT to enable the Google Cloud APIs and create an Docker repository in Artififact Registry

## Setup
*NOTE:* The `setters.yaml` is a symbolic link to the `setters.yaml` file in the parent `infra` folder. Once the variables are setup in the `setters.yaml` file, run:


## Deploy
`kpt fn render`
`kpt live init`
`kpt live apply`