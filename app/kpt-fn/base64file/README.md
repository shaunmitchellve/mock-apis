## Base64File Encoding ##

This KPT Function will base64encode a file and update a specific field in a KRM with that base64 encoded string. It takes in a function config map that has the following key/value settings:

| key | description |
| -- | -- |
| target-kind | The kind of the target resource that will be updated with the base64 encoded string |\
| target-field | The full path to the field that updated. Use dot notation for the path, e.g: spec.openapiDocuments.document[0].contents |
| source-file | The full path to the file to base64 encode |

# Local Exec
The GO program in this folder has been build for the linux OS an amd64 architecture. If you want to re-build the program for your a different OS / Architecture, just run `build.sh`.

# Development
You can run `kpt fn source --fn-config=testdata/_fnconfig.yaml | go run main.go` to run the function against your testing KRM ResourceList files in the testdata folder.