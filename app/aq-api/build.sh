#!/bin/bash

while getopts p:r:v: flag
do
    case "${flag}" in
        p) projectid=${OPTARG};;
        r) region=${OPTARG};;
        v) version=${OPTARG};;
    esac
done

if [ ${#projectid} -eq 0 ] || [ ${#region} -eq 0 ]; then
    echo "missing required field(s)\n"
    echo "Usage: ./deploy.sh -p <project-id> -r <region> -v <version>"
    exit 1
fi


gcloud builds submit . --project=$projectid --config=cloudbuild.yaml --substitutions=_PROJECT_ID=$projectid,_REGION=$region,_VERSION=$version