#!/usr/bin/env bash
set -euo pipefail

###############################################
# Core Bank Kubernetes deployment script
#
# Deploys all necessary manifests to set up
# the core bank microservices in a Kubernetes cluster.
# Usage:
#   ./deploy.sh                    # deploy to default namespace
#   ./deploy.sh -n core-bank       # specify namespace
#   ./deploy.sh -y                 # skip interactive confirmation
#   ./deploy.sh -w                 # wait for workloads to be ready
#
# Options:
#   -n <namespace>   Target namespace (default: core-bank)
#   -y               Auto-confirm (no prompt)
#   -w               Wait for workloads to be ready
#   -c               Check images exist in registry before apply
#   -b               Auto-build missing images using dockerbuild.sh
#   -L <target>      Load built images to local cluster: 'kind', 'minikube', 'docker-desktop'
#   -h               Show help
#
# Notes:
#   * Assumes kubectl context is already set.
#   * Creates namespace if it does not exist.
###############################################

print_help() {
  grep '^#' "$0" | sed 's/^# //'
}

NAMESPACE="core-bank"
AUTO_CONFIRM=false
WAIT_READY=false
CHECK_IMAGES=false
AUTO_BUILD=false
LOAD_TARGET=""

while getopts ":n:ywchbL:" opt; do
  case $opt in
    n) NAMESPACE="$OPTARG" ;;
    y) AUTO_CONFIRM=true ;;
    w) WAIT_READY=true ;;
    c) CHECK_IMAGES=true ;;
    b) AUTO_BUILD=true ;;
    L) LOAD_TARGET="$OPTARG" ;;
    h) print_help; exit 0 ;;
    :) echo "Missing value for -$OPTARG" >&2; exit 1 ;;
    \?) echo "Invalid option: -$OPTARG" >&2; print_help; exit 1 ;;
  esac
done

shift $((OPTIND-1))
echo "🚀 Deploying Core Bank services to Kubernetes namespace '$NAMESPACE'..."

# Check if namespace exists
if ! kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
  echo "Namespace '$NAMESPACE' does not exist. Creating..."
  kubectl create namespace "$NAMESPACE"
else
  echo "Namespace '$NAMESPACE' already exists."
fi

MANIFESTS=(
    "K8s/namespace.yml"
    "K8s/authentication-service.yml"
    "K8s/customer-service.yml"
    "K8s/account-service.yml"
    "K8s/transaction-service.yml"
    "K8s/postgres.yml"
    "K8s/kafka.yml"
    "K8s/ingress-loadbalancer.yml"
    "K8s/network-policy.yml"
)

# If using a local cluster, ensure common third-party images are present locally
case "$LOAD_TARGET" in
  docker-desktop|minikube|kind)
    echo "🔧 Ensuring required third-party images are present locally..."
    THIRDPARTY_REQUIRED=(
      "postgres:14"
      "confluentinc/cp-kafka:latest"
    )
    for img in "${THIRDPARTY_REQUIRED[@]}"; do
      if docker image inspect "$img" >/dev/null 2>&1; then
        echo "✅ Present: $img"
      else
        echo "📥 Pulling $img ..."
        if docker pull "$img"; then
          echo "✅ Pulled: $img"
        else
          echo "⚠️  Failed to pull: $img" >&2
        fi
      fi
    done
    ;;
esac

# Optional: check container images referenced in manifests are available
if $CHECK_IMAGES; then
  echo "🔎 Checking images referenced in manifests..."
  if ! command -v docker >/dev/null 2>&1; then
    echo "Warning: docker CLI not found; skipping image checks" >&2
  else
    IMAGES=()
    for manifest in "${MANIFESTS[@]}"; do
      if [[ -f "$manifest" ]]; then
        while IFS= read -r raw; do
          # Strip leading 'image:' and any inline comments after the value
          val=$(echo "$raw" | sed -E 's/^[[:space:]]*image:[[:space:]]*//')
          val=$(echo "$val" | sed -E 's/[[:space:]]+#.*$//')
          # Trim whitespace
          val=$(echo "$val" | xargs)
          # Basic validation to avoid malformed tokens (e.g., typos like latestest)
          if [[ -n "$val" ]] && [[ "$val" =~ ^[A-Za-z0-9._:/-]+(:[A-Za-z0-9._-]+)?$ ]]; then
            IMAGES+=("$val")
          fi
        done < <(grep -E "^[[:space:]]*image:[[:space:]]*" "$manifest")
      fi
    done
    # unique
    mapfile -t UNIQUE_IMAGES < <(printf "%s\n" "${IMAGES[@]}" | awk '!seen[$0]++')
    if [[ ${#UNIQUE_IMAGES[@]} -eq 0 ]]; then
      echo "No images found in manifests to check."
    else
      echo "Found images: ${UNIQUE_IMAGES[*]}"
      FAILED_IMAGES=()
      for image in "${UNIQUE_IMAGES[@]}"; do
        echo "Checking $image..."
        # If targeting local clusters, prefer local daemon check
        local_check=false
        case "$LOAD_TARGET" in
          docker-desktop|minikube|kind) local_check=true ;;
        esac
        if $local_check; then
          if docker image inspect "$image" >/dev/null 2>&1; then
            echo "✅ Local image present: $image"
          else
            echo "❌ Local image missing: $image"
            FAILED_IMAGES+=("$image")
          fi
        else
          if docker manifest inspect "$image" >/dev/null 2>&1; then
            echo "✅ Exists: $image"
          else
            echo "❌ Missing: $image"
            FAILED_IMAGES+=("$image")
          fi
        fi
      done
      if [[ ${#FAILED_IMAGES[@]} -gt 0 ]]; then
        echo "❌ Missing images: ${FAILED_IMAGES[*]}" >&2
        # If targeting local clusters, try to pull third-party images automatically
        local_cluster=false
        case "$LOAD_TARGET" in
          docker-desktop|minikube|kind) local_cluster=true ;;
        esac
        if $local_cluster; then
          echo "📥 Attempting to pull third-party images locally..."
          FIRST_PARTY_SET="authentication-service customer-service account-service transaction-service"
          for img in "${FAILED_IMAGES[@]}"; do
            name_tag="${img##*/}"; name_only="${name_tag%%:*}"
            if [[ ! " $FIRST_PARTY_SET " =~ " $name_only " ]]; then
              echo "Pulling $img ..."
              if docker pull "$img"; then
                echo "✅ Pulled: $img"
                # Remove from FAILED_IMAGES if pull succeeded
                for i in "${!FAILED_IMAGES[@]}"; do
                  if [[ "${FAILED_IMAGES[$i]}" == "$img" ]]; then unset 'FAILED_IMAGES[i]'; fi
                done
              else
                echo "⚠️  Failed to pull: $img" >&2
              fi
            fi
          done
          # Repack FAILED_IMAGES to remove gaps
          mapfile -t FAILED_IMAGES < <(printf "%s\n" "${FAILED_IMAGES[@]}" | awk 'NF')
        fi
        if $AUTO_BUILD; then
          echo "⚙️  Auto-build enabled (-b). Attempting to build missing images..."
          # Only attempt to build first-party services
          FIRST_PARTY=(authentication-service customer-service account-service transaction-service)
          SERVICES_TO_BUILD=()
          img_registry=""; img_tag=""; inferred=false
          for missing in "${FAILED_IMAGES[@]}"; do
            name_tag="${missing##*/}"
            name_only="${name_tag%%:*}"
            tag_part="${name_tag##*:}"
            for fp in "${FIRST_PARTY[@]}"; do
              if [[ "$name_only" == "$fp" ]]; then
                SERVICES_TO_BUILD+=("$name_only")
                if ! $inferred; then
                  # infer common tag from the first first-party missing image
                  img_tag="$tag_part"
                  inferred=true
                fi
              fi
            done
          done
          # Deduplicate services
          mapfile -t UNIQUE_SERVICES < <(printf "%s\n" "${SERVICES_TO_BUILD[@]}" | awk '!seen[$0]++')
          if [[ ${#UNIQUE_SERVICES[@]} -eq 0 ]]; then
            echo "No first-party images to build. Skipping auto-build." >&2
            echo "Tip: third-party images (postgres, kafka, zookeeper) should exist in registry or be pulled at deploy time." >&2
            exit 2
          fi
          build_cmd=("bash" "$(dirname "$0")/dockerbuild.sh")
          if [[ -n "$img_tag" && "$img_tag" != "$name_tag" ]]; then
            build_cmd+=("-t" "$img_tag")
          fi
          # If load target provided, prefer local load over push
          if [[ -n "$LOAD_TARGET" ]]; then
            build_cmd+=("-L" "$LOAD_TARGET")
          else
            # Heuristic: auto-detect common local clusters by current kubectl context
            current_ctx=$(kubectl config current-context 2>/dev/null || echo "")
            case "$current_ctx" in
              kind-*) build_cmd+=("-L" "kind") ;;
              minikube) build_cmd+=("-L" "minikube") ;;
              docker-desktop) build_cmd+=("-L" "docker-desktop") ;;
            esac
          fi
          # Push only if no local load target set
          if [[ "${build_cmd[*]}" != *" -L "* ]]; then
            build_cmd+=("-p")
          fi
          build_cmd+=("-s" "$(IFS=,; echo "${UNIQUE_SERVICES[*]}")")
          echo "Running: ${build_cmd[*]}"
          "${build_cmd[@]}"
          echo "🔁 Re-checking images after build..."
          STILL_MISSING=()
          for image in "${FAILED_IMAGES[@]}"; do
            if [[ -n "$LOAD_TARGET" ]] || [[ "$current_ctx" == kind-* || "$current_ctx" == minikube || "$current_ctx" == docker-desktop ]]; then
              if docker image inspect "$image" >/dev/null 2>&1; then
                echo "✅ Now present locally: $image"
              else
                echo "❌ Still missing locally: $image"
                STILL_MISSING+=("$image")
              fi
            else
              if docker manifest inspect "$image" >/dev/null 2>&1; then
                echo "✅ Now exists: $image"
              else
                echo "❌ Still missing: $image"
                STILL_MISSING+=("$image")
              fi
            fi
          done
          if [[ ${#STILL_MISSING[@]} -gt 0 ]]; then
            echo "ERROR: Some images are still missing: ${STILL_MISSING[*]}" >&2
            exit 2
          fi
        else
          echo "ERROR: Images missing and auto-build disabled. Use -b or build manually." >&2
          echo "Tip: build and push with ./Deployment/dockerbuild.sh -p -r <registry> -t <tag>" >&2
          exit 2
        fi
      fi
      # Fallback: ensure all first-party images exist locally when using a local cluster and auto-build is enabled.
      # This guards against edge cases where a manifest image line was skipped and avoids init-container deadlocks.
      if $AUTO_BUILD; then
        local_cluster=false
        case "$LOAD_TARGET" in
          docker-desktop|minikube|kind) local_cluster=true ;;
        esac
        if $local_cluster; then
          EXPECTED_FP=(authentication-service account-service customer-service transaction-service)
          MISSING_FP=()
          for svc in "${EXPECTED_FP[@]}"; do
            if ! docker image inspect "$svc:latest" >/dev/null 2>&1; then
              MISSING_FP+=("$svc")
            fi
          done
          if [[ ${#MISSING_FP[@]} -gt 0 ]]; then
            echo "⚙️  Fallback: building missing first-party images: ${MISSING_FP[*]}"
            fb_cmd=("bash" "$(dirname "$0")/dockerbuild.sh" "-t" "latest" "-L" "$LOAD_TARGET" "-s" "$(IFS=,; echo "${MISSING_FP[*]}")")
            echo "Running: ${fb_cmd[*]}"
            "${fb_cmd[@]}"
          fi
        fi
      fi
    fi
  fi
fi

for manifest in "${MANIFESTS[@]}"; do
  echo "Applying manifest: $manifest"
  kubectl apply -n "$NAMESPACE" -f "$manifest"
done

if ! $AUTO_CONFIRM; then
  read -r -p "Proceed with deployment? (y/N) " ans
  case "${ans:-}" in
    y|Y) echo "Proceeding..." ;;
    *) echo "Aborted."; exit 0 ;;
  esac
fi

if $WAIT_READY; then
  echo "Waiting for deployments to be ready..."
  for deploy in authentication-service customer-service account-service transaction-service; do
    echo "Waiting for deployment/$deploy to be ready..."
    kubectl wait --namespace "$NAMESPACE" --for=condition=available --timeout=300s deployment/"$deploy"
  done
  echo "All deployments are ready."
fi

echo "✅ Deployment completed."