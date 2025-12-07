.PHONY: all build run test proto clean harness execute

# Build settings
BINARY_NAME=server
CMD_PATH=./cmd/server

# Build the application
build:
	go build -o $(BINARY_NAME) $(CMD_PATH)

# Run the application
run: build
	./$(BINARY_NAME)

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-v:
	go test -v ./...

# Generate protobuf code
proto:
	buf generate

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f *.exe

# Download dependencies
deps:
	go mod tidy

# Run go vet
vet:
	go vet ./...

# Run harness script
harness:
	py scripts/execute_harness.py

# Show status only
harness-status:
	py scripts/execute_harness.py --status

# Dry run
harness-dry:
	py scripts/execute_harness.py --dry-run

# Build and test
check: vet test build

# Full rebuild
rebuild: clean deps build

# Docker settings
IMAGE_NAME=asia-northeast1-docker.pkg.dev/cloudsql-sv/postgres-prod/postgres-prod
ENVOY_IMAGE_NAME=asia-northeast1-docker.pkg.dev/cloudsql-sv/postgres-prod/postgres-prod-envoy
IMAGE_TAG=latest

# Cloud Run settings
PROJECT_ID=cloudsql-sv
REGION=asia-northeast1
CLOUDSQL_INSTANCE=postgres-prod
DB_NAME=myapp
DB_USER=747065218280-compute@developer
SERVICE_ACCOUNT=747065218280-compute@developer.gserviceaccount.com
FRONTEND_URL=https://mtama-front.mtamaramu.com/auth/callback

# Build Docker image locally
docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

# Build Envoy sidecar image
docker-build-envoy:
	docker build -f Dockerfile.envoy -t $(ENVOY_IMAGE_NAME):$(IMAGE_TAG) .

# gcloud path (Windows)
GCLOUD="/c/Users/mtama/AppData/Local/Google/Cloud SDK/google-cloud-sdk/bin/gcloud.cmd"

# Authenticate Docker to Artifact Registry
docker-auth:
	@TOKEN=$$($(GCLOUD) auth print-access-token) && docker login -u oauth2accesstoken -p $$TOKEN https://asia-northeast1-docker.pkg.dev

# Push Docker image to Artifact Registry
docker-push:
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

# Push Envoy image to Artifact Registry
docker-push-envoy:
	docker push $(ENVOY_IMAGE_NAME):$(IMAGE_TAG)

# Generate resolved service.yaml and deploy to Cloud Run
cloud-run-deploy:
	sed -e "s|\$${PROJECT_ID}|$(PROJECT_ID)|g" \
		-e "s|\$${REGION}|$(REGION)|g" \
		-e "s|\$${REPO_NAME}|postgres-prod|g" \
		-e "s|\$${CLOUDSQL_INSTANCE}|$(CLOUDSQL_INSTANCE)|g" \
		-e "s|\$${DB_NAME}|$(DB_NAME)|g" \
		-e "s|\$${DB_USER}|$(DB_USER)|g" \
		-e "s|\$${SERVICE_ACCOUNT}|$(SERVICE_ACCOUNT)|g" \
		-e "s|\$${TAG}|$(IMAGE_TAG)|g" \
		-e "s|\$${FRONTEND_URL}|$(FRONTEND_URL)|g" \
		service.yaml > service-resolved.yaml
	$(GCLOUD) run services replace service-resolved.yaml --region=$(REGION) --project=$(PROJECT_ID)

# Local build and deploy (no Cloud Build charges)
deploy-local: docker-build docker-build-envoy docker-auth docker-push docker-push-envoy cloud-run-deploy

# Force deploy with timestamp tag (always creates new revision)
deploy-force:
	$(eval TIMESTAMP := $(shell date +%Y%m%d%H%M%S))
	docker build -t $(IMAGE_NAME):$(TIMESTAMP) .
	docker build -f Dockerfile.envoy -t $(ENVOY_IMAGE_NAME):$(TIMESTAMP) .
	@TOKEN=$$($(GCLOUD) auth print-access-token) && docker login -u oauth2accesstoken -p $$TOKEN https://asia-northeast1-docker.pkg.dev
	docker push $(IMAGE_NAME):$(TIMESTAMP)
	docker push $(ENVOY_IMAGE_NAME):$(TIMESTAMP)
	sed -e "s|\$${PROJECT_ID}|$(PROJECT_ID)|g" \
		-e "s|\$${REGION}|$(REGION)|g" \
		-e "s|\$${REPO_NAME}|postgres-prod|g" \
		-e "s|\$${CLOUDSQL_INSTANCE}|$(CLOUDSQL_INSTANCE)|g" \
		-e "s|\$${DB_NAME}|$(DB_NAME)|g" \
		-e "s|\$${DB_USER}|$(DB_USER)|g" \
		-e "s|\$${SERVICE_ACCOUNT}|$(SERVICE_ACCOUNT)|g" \
		-e "s|\$${TAG}|$(TIMESTAMP)|g" \
		-e "s|\$${FRONTEND_URL}|$(FRONTEND_URL)|g" \
		service.yaml > service-resolved.yaml
	$(GCLOUD) run services replace service-resolved.yaml --region=$(REGION) --project=$(PROJECT_ID)
