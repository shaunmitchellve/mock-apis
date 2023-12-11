#!/bin/bash

while getopts p:r:v:a: flag
do
    case "${flag}" in
        p) projectid=${OPTARG};;
        r) region=${OPTARG};;
        a) repo=${OPTARG};;
        v) version=${OPTARG};;
    esac
done

if [ ${#projectid} -eq 0 ] || [ ${#region} -eq 0 ] || [ ${#repo} -eq 0 ] || [ ${#version} -eq 0 ]; then
    echo "missing required field(s)"
    echo "Usage: ./deploy.sh -p <project-id> -r <region> -v <version> -a <repo>"
    exit 1
fi


gcloud builds submit . --project=$projectid --config=cloudbuild.yaml --substitutions=_PROJECT_ID=$projectid,_REGION=$region,_VERSION=$version,_REPO=$repo