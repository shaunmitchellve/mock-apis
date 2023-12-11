# Copyright 2023 Shaun Mitchell

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#  	http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash

while getopts p:r:a: flag
do
    case "${flag}" in
        p) project=${OPTARG};;
        r) region=${OPTARG};;
        a) repo=${OPTARG};;
    esac
done

if [ ${#project} -eq 0 ] || [ ${#region} -eq 0 ] || [ ${#repo} -eq 0 ]; then
    echo "missing required field(s)\n"
    echo "Usage: ./deploy.sh -p <project-id> -r <region> -a <repo>"
    exit 1
fi

APP_VERSION="v0.0.4"

echo "Installing API Gateway alpha CRDs into Config Controller"
git clone https://github.com/GoogleCloudPlatform/k8s-config-connector
cd k8s-config-connector
echo 'v'$(kubectl get ns cnrm-system -o jsonpath='{.metadata.annotations.cnrm\.cloud\.google\.com/version}') | xargs git checkout
kubectl apply -f crds/apigateway_v1alpha1_apigatewayapi.yaml
kubectl apply -f crds/apigateway_v1alpha1_apigatewayapiconfig.yaml
kubectl apply -f crds/apigateway_v1alpha1_apigatewaygateway.yaml
cd ..
rm -rf k8s-config-connector

echo "Enabling APIs and creating artifact registry in project $project"
sed -i "s|project-id:.*$|project-id: $project|g" infra/setters.yaml
sed -i "s|region:.*$|region: $region|g" infra/setters.yaml
sed -i "s|app-path:.*$|app-path: $region-docker.pkg.dev/$project/mock-apis/aq-mock-api:$APP_VERSION|g" infra/setters.yaml

kpt fn render infra/0-apis
kpt live init infra/0-apis
kpt live apply infra/0-apis

echo "Building AQ Mock API app"
cd app/aq-api
./build.sh -p $project -r $region -v $APP_VERSION -a $repo
cd ../../

echo "Creating Firestore Database"
gcloud alpha firestore databases create \
--database=aq-mock-data \
--location=nam5 \
--type=firestore-native \
--project="$project"

echo "Creating Cloud Run Service"
kpt fn render infra/1-run
kpt live init infra/1-run
kpt live apply infra/1-run

echo "Getting Run URL"
url=$(gcloud run services describe aq-mock-service --project=$project --region=$region --format='value(status.address.url)')
sed -i "s|address:.*$|address: $url|" infra/2-gateway/api-spec.yaml

echo "Creating API Gateway"
cd infra/2-gateway
kpt fn render . --allow-exec
kpt live init .
kpt live apply .
cd ../../


echo "Enabling API Managed Service for API Key"
ms=$(gcloud api-gateway apis describe airqualityapi --project=$project --format='value(managedService)')
gcloud services enable $ms --project=$project

echo "Creating and securing API Key"
gcloud beta services api-keys create --display-name="AQ API Key" \
--api-target=service=$ms \
--project=$project

echo "API Gateway URL: "
gcloud alpha api-gateway gateways describe airqualityapi-gateway --project=$project --location=$region --format='value(defaultHostname)'