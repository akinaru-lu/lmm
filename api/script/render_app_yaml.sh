#!/bin/bash

python /gae/render_app_yaml.py --input ./template.yaml --output ./app.yaml \
	--ASSET_BUCKET_NAME ${ASSET_BUCKET_NAME} \
	--DATASTORE_PROJECT_ID ${DATASTORE_PROJECT_ID} \
	--GIN_MODE debug \
	--GCP_PROJECT_ID ${GCP_PROJECT_ID} \
	--GOOGLE_APPLICATION_CREDENTIALS ${GOOGLE_APPLICATION_CREDENTIALS} \
	--GO_VERSION ${GO_VERSION} \
	--LMM_API_TOKEN_KEY ${LMM_API_TOKEN_KEY} \
	--LMM_DOMAIN ${LMM_DOMAIN} \
	--PUBSUB_PROJECT_ID ${PUBSUB_PROJECT_ID}
