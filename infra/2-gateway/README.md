# Infra - Step 2 (API Gateway)

Use KPT to setup and deploy the API Gateway. The API Gateway KCC CRDs are in Alpha right now. Make sure they have been installed into your KCC environment.

## Setup
*NOTE:* The `setters.yaml` is a symbolic link to the `setters.yaml` file in the parent `infra` folder.

Update the file `api-spec.yaml` `address` field in the segment:

```
x-google-backend:
  address:
```

to the URL of the Cloud Run service.

## Deploy
This KPT process uses a custom KPT Function that has been compiled already in the `/app/kpt-fn/base6file` folder. This function will base64 encode the Swagger file `api-spec.yaml` and update the `api-gateway.yaml` with that encoded string.

Once the variables are setup in the `setters.yaml` file, run:

`kpt fn render --allow-exex` *NOTE* the --allow-exec flag\
`kpt live init`\
`kpt live apply`