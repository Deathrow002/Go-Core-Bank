#!/usr/bin/env bash
set -euo pipefail

###############################################
# Go-Core-Bank Docker build helper
#
# Builds (and optionally pushes) the core bank
# microservice images in a consistent manner.
#
# Supported services (from docker-compose.yml):
#   authentication-service  -> Authentication-Service/dockerfile
#   customer-service        -> Customer-Service/Dockerfile
#   account-service         -> Account-Service/Dockerfile
#   transaction-service     -> Transaction-Service/Dockerfile
#
# Usage:
#   ./dockerbuild.sh                # build with auto tag (git SHA or 'latest')
#   ./dockerbuild.sh -t v1.2.3       # build all with tag v1.2.3
#   ./dockerbuild.sh -t v1.2.3 -r myregistry.io/myproj  # prefix images with registry path
#   ./dockerbuild.sh -p -r ghcr.io/org/repo -t $(git rev-parse --short HEAD)
#   ./dockerbuild.sh -s account-service,customer-service  # build subset
#   (Works from any directory; paths resolved relative to script.)
#
# Options:
#   -t <tag>        Image tag (default: git short SHA or 'latest')
#   -r <registry>   Registry prefix (example: ghcr.io/org/repo)
#   -p              Push images after build
#   -n              Disable Docker build cache (adds --no-cache)
#   -s <list>       Comma-separated service list to build (default: all)
#   -j              Parallel build (background jobs)
#   -L <target>     Load built images into local cluster: 'kind', 'minikube', or 'docker-desktop'
#   -h              Show help
#
# Notes:
#   * Requires Docker daemon running.
#   * On Windows, run via Git Bash / WSL for best compatibility.
#   * Uses Docker BuildKit if available.
###############################################

if ! command -v docker >/dev/null 2>&1; then
	echo "Error: docker CLI not found in PATH" >&2
	exit 1
fi

export DOCKER_BUILDKIT=1

# Resolve script directory and repository root (assumes script lives in Deployment/ at repo root)
SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd -- "$SCRIPT_DIR/.." && pwd)"

# Derive default tag from git if possible
DEFAULT_TAG="latest"
if command -v git >/dev/null 2>&1 && git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
	if GIT_SHA=$(git rev-parse --short HEAD 2>/dev/null); then
		DEFAULT_TAG="$GIT_SHA"
	fi
fi

TAG="$DEFAULT_TAG"
REGISTRY=""
PUSH=false
NO_CACHE=false
PARALLEL=false
REQUESTED_SERVICES=""
LOAD_TARGET=""

print_help() {
	grep '^#' "$0" | sed 's/^# //'
}

while getopts ":t:r:ps:njhL:" opt; do
	case $opt in
		t) TAG="$OPTARG" ;;
		r) REGISTRY="$OPTARG" ;;
		p) PUSH=true ;;
		n) NO_CACHE=true ;;
		s) REQUESTED_SERVICES="$OPTARG" ;;
		j) PARALLEL=true ;;
		L) LOAD_TARGET="$OPTARG" ;;
		h) print_help; exit 0 ;;
		:) echo "Missing value for -$OPTARG" >&2; exit 1 ;;
		\?) echo "Invalid option: -$OPTARG" >&2; print_help; exit 1 ;;
	esac
done
shift $((OPTIND-1))

# Canonical service list (compose service names)
ALL_SERVICES=(
	authentication-service
	customer-service
	account-service
	transaction-service
)

# Map service -> directory & dockerfile
service_dir() {
	case "$1" in
		authentication-service) echo "Authentication-Service" ;;
		customer-service)       echo "Customer-Service" ;;
		account-service)        echo "Account-Service" ;;
		transaction-service)    echo "Transaction-Service" ;;
		*) return 1 ;;
	esac
}

service_dockerfile() {
	case "$1" in
		authentication-service) echo "dockerfile" ;;
		*) echo "Dockerfile" ;;
	esac
}

IFS=',' read -r -a BUILD_SET <<< "${REQUESTED_SERVICES}" || true

SERVICES_TO_BUILD=()
if [[ -n "$REQUESTED_SERVICES" ]]; then
	for s in "${BUILD_SET[@]}"; do
		s_trim=$(echo "$s" | xargs)
		if [[ -z "$s_trim" ]]; then continue; fi
		found=false
		for a in "${ALL_SERVICES[@]}"; do
			if [[ "$a" == "$s_trim" ]]; then found=true; break; fi
		done
		if ! $found; then
			echo "Warning: Unknown service '$s_trim' - skipping" >&2
			continue
		fi
		SERVICES_TO_BUILD+=("$s_trim")
	done
else
	SERVICES_TO_BUILD=("${ALL_SERVICES[@]}")
fi

if [[ ${#SERVICES_TO_BUILD[@]} -eq 0 ]]; then
	echo "No valid services selected. Nothing to do." >&2
	exit 1
fi

echo "==> Building services: ${SERVICES_TO_BUILD[*]}"
echo "==> Tag: $TAG"
if [[ -n "$REGISTRY" ]]; then echo "==> Registry prefix: $REGISTRY"; fi
echo "==> Push: $PUSH"; 
echo "==> No cache: $NO_CACHE"; 
echo "==> Parallel: $PARALLEL"
if [[ -n "$LOAD_TARGET" ]]; then echo "==> Load target: $LOAD_TARGET"; fi

docker_build_args=()
if $NO_CACHE; then docker_build_args+=("--no-cache"); fi

build_service() {
	local svc="$1"
	local dir dockerfile image ctx dockerfile_path
	dir=$(service_dir "$svc") || { echo "Unknown service '$svc'" >&2; return 1; }
	dockerfile=$(service_dockerfile "$svc")
	ctx="$ROOT_DIR/$dir"
	dockerfile_path="$ctx/$dockerfile"
	if [[ ! -f "$dockerfile_path" ]]; then
		echo "Error: Dockerfile '$dockerfile_path' not found for service '$svc'" >&2
		return 1
	fi
	image="${svc}:$TAG"
	if [[ -n "$REGISTRY" ]]; then
		image="${REGISTRY}/${svc}:$TAG"
	fi
	echo "--- Building $image (context: $ctx, dockerfile: $dockerfile)" 
	docker build "${docker_build_args[@]}" -f "$dockerfile_path" -t "$image" "$ctx"
	if $PUSH; then
		echo "--- Pushing $image"
		docker push "$image"
	fi
	# Optional: load into local k8s cluster (kind or minikube)
	if [[ -n "$LOAD_TARGET" ]]; then
		case "$LOAD_TARGET" in
			kind)
				if command -v kind >/dev/null 2>&1; then
					echo "--- Loading image into kind: $image"
					kind load docker-image "$image" || { echo "Failed to load into kind" >&2; return 1; }
				else
					echo "Error: 'kind' not found in PATH" >&2; return 1
				fi
				;;
			minikube)
				if command -v minikube >/dev/null 2>&1; then
					echo "--- Loading image into minikube: $image"
					minikube image load "$image" || { echo "Failed to load into minikube" >&2; return 1; }
				else
					echo "Error: 'minikube' not found in PATH" >&2; return 1
				fi
				;;
			docker-desktop)
				# Docker Desktop Kubernetes shares the Docker daemon, so locally built images
				# are available to the cluster as long as manifests use the same name:tag
				# and imagePullPolicy is IfNotPresent. No extra load step required.
				echo "--- Docker Desktop detected; using local images (${image}). Ensure imagePullPolicy=IfNotPresent."
				;;
			*)
				echo "Error: unsupported load target '$LOAD_TARGET' (use 'kind' or 'minikube')" >&2; return 1
				;;
		esac
	fi
}

FAILED=()

if $PARALLEL; then
	echo "Parallel build enabled"
	for svc in "${SERVICES_TO_BUILD[@]}"; do
		build_service "$svc" &
	done
	wait
else
	for svc in "${SERVICES_TO_BUILD[@]}"; do
		if ! build_service "$svc"; then
			FAILED+=("$svc")
		fi
	done
fi

if [[ ${#FAILED[@]} -gt 0 ]]; then
	echo "Build failures: ${FAILED[*]}" >&2
	exit 2
fi

echo "All requested services built successfully." 
if $PUSH; then
	echo "Images pushed with tag '$TAG'."
fi

