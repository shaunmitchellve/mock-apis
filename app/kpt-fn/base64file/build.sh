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

echo 'Build local exec or build container?'
read -p 'local or container: ' buildoption

if [ buildoption = 'container' ]; then
    echo 'Update Artifact Repo variables\n'
    read -p 'Project ID: ' projectid
    read -p 'Repo Name: ' reponame
    read -p 'Region: ' region

    gcloud builds submit . --project=$projectid --config=cloudbuild.yaml --substitutions=_PROJECT_ID=$projectid,_REPO_NAME=$reponame,_REGION=$region
else
    echo "What OS and Architecture:"
    read -p 'OS: ' os
    read -p 'Arch: ' arch

   env GOOS=$os GOARCH=$arch go build
fi

