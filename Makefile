# ---------------------------------------------------------------------
# -- Which container tool to use
# ---------------------------------------------------------------------
CONTAINER_TOOL ?= docker
# ---------------------------------------------------------------------
# -- Image URL to use all building/pushing image targets
# ---------------------------------------------------------------------
IMAGE_TAG ?= webhook-router:dev
# ---------------------------------------------------------------------
# -- It's required when you want to use k3s and nerdctl
# -- $ export CONTAINER_TOOL_NAMESPACE_ARG="--namespace k8s.io"
# ---------------------------------------------------------------------
CONTAINER_TOOL_NAMESPACE_ARG ?=
# -- Set additional argumets to container tool
# -- To use, set an environment variable:
# -- $ export CONTAINER_TOOL_ARGS="--build-arg GOARCH=arm64 --platform=linux/arm64"
# ---------------------------------------------------------------------
CONTAINER_TOOL_ARGS ?=

k3d-create:
	k3d cluster create myk3s
	kubectl get pod

k3d-delete:
	k3d cluster delete myk3s

webhook-router-deploy:
	kubectl apply -f k8s/secret-config.yaml
	kubectl apply -f k8s/secret-env.yaml
	kubectl apply -f k8s/service.yaml
	kubectl delete -f k8s/deployment.yaml
	kubectl apply -f k8s/deployment.yaml

build: ## Build a container
	$(CONTAINER_TOOL) build ${CONTAINER_TOOL_ARGS} -t ${IMAGE_TAG} . ${CONTAINER_TOOL_NAMESPACE_ARG}
	$(CONTAINER_TOOL) save ${CONTAINER_TOOL_NAMESPACE_ARG} ${IMAGE_TAG} -o webhook-router-image.tar

k3d_image: build ## rebuild the docker images and upload into your k3d cluster
	@k3d image import webhook-router-image.tar -c myk3s