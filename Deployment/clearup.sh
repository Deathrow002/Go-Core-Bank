#!/usr/bin/env bash
set -euo pipefail

###############################################
# Core Bank Kubernetes cleanup helper
#
# Removes deployed manifests, workloads, and (optionally)
# the namespace & persistent storage for this project.
#
# Usage:
#   ./clearup.sh                    # delete manifests & workloads only
#   ./clearup.sh -n core-bank       # specify namespace
#   ./clearup.sh -F                 # force delete namespace at end
#   ./clearup.sh -P                 # also delete PVCs & PVs
#   ./clearup.sh -d                 # dry-run (show what would happen)
#   ./clearup.sh -y                 # skip interactive confirmation
#   ./clearup.sh -w                 # wait for deletions to finish
#   ./clearup.sh -i                 # remove local service images
#   ./clearup.sh -I                 # remove local service + third-party images
#
# Options:
#   -n <namespace>   Target namespace (default: core-bank)
#   -F               Delete namespace after resources
#   -P               Delete persistent volume claims (and matching PVs)
#   -d               Dry run (no changes, just output)
#   -y               Auto-confirm (no prompt)
#   -w               Wait for workloads to terminate
#   -i               Remove local Docker images for first-party services
#   -I               Remove local Docker images including third-party deps
#   -h               Show help
#
# Notes:
#   * Assumes kubectl context is already set.
#   * Ignores not-found errors gracefully.
#   * PV deletion is cluster-impacting; use cautiously.
###############################################

print_help() {
  grep '^#' "$0" | sed 's/^# //' 
}

# Resolve script directory so relative paths work regardless of cwd
SCRIPT_DIR=$(cd -- "$(dirname "$0")" && pwd)
ROOT_DIR=$(cd -- "$SCRIPT_DIR/.." && pwd)

NAMESPACE="core-bank"
DELETE_NAMESPACE=false
DELETE_PERSISTENCE=false
DRY_RUN=false
AUTO_CONFIRM=false
WAIT_TERMINATION=false
REMOVE_IMAGES=false
REMOVE_THIRDPARTY=false

while getopts ":n:FPdywhiI" opt; do
  case $opt in
    n)
      # Guard against missing namespace value (e.g., '-n -w')
      if [[ -z "$OPTARG" || "$OPTARG" == -* ]]; then
        echo "Error: -n requires a namespace value" >&2
        exit 1
      fi
      NAMESPACE="$OPTARG"
      ;;
    F) DELETE_NAMESPACE=true ;;
    P) DELETE_PERSISTENCE=true ;;
    d) DRY_RUN=true ;;
    y) AUTO_CONFIRM=true ;;
    w) WAIT_TERMINATION=true ;;
    i) REMOVE_IMAGES=true ;;
    I) REMOVE_IMAGES=true; REMOVE_THIRDPARTY=true ;;
    h) print_help; exit 0 ;;
    :) echo "Missing value for -$OPTARG" >&2; exit 1 ;;
    \?) echo "Invalid option: -$OPTARG" >&2; print_help; exit 1 ;;
  esac
done
shift $((OPTIND-1))

# Allow positional namespace if -n was omitted
if [[ $# -ge 1 && -z "$NAMESPACE_SET" ]]; then
  if [[ "$1" != -* ]]; then
    NAMESPACE="$1"
    shift
  fi
fi

echo "🧹 Cleaning up Core Bank services from Kubernetes..."
echo "================================================="
echo "Namespace: $NAMESPACE"
echo "Delete namespace: $DELETE_NAMESPACE"
echo "Delete persistence (PVC/PV): $DELETE_PERSISTENCE"
echo "Dry run: $DRY_RUN"
echo "Wait termination: $WAIT_TERMINATION"
echo "Remove images: $REMOVE_IMAGES"
echo "Include third-party images: $REMOVE_THIRDPARTY"

if ! command -v kubectl >/dev/null 2>&1; then
  echo "Error: kubectl not found in PATH" >&2
  exit 1
fi

if ! $AUTO_CONFIRM; then
  read -r -p "Proceed with cleanup? (y/N) " ans
  case "${ans:-}" in
    y|Y) echo "Proceeding..." ;;
    *) echo "Aborted."; exit 0 ;;
  esac
fi

run_cmd() {
  echo "[CMD] $*" >&2
  if ! $DRY_RUN; then
    "$@"
  fi
}

# 1. Delete manifests (if K8s directory exists)
K8S_DIR="$ROOT_DIR/K8s"
if [[ -d "$K8S_DIR" ]]; then
  echo "Deleting manifests from $K8S_DIR (ignore not found)..."
  if $DRY_RUN; then
    echo "(dry-run) kubectl delete -f $K8S_DIR -n $NAMESPACE --ignore-not-found"
  else
    kubectl delete -f "$K8S_DIR" -n "$NAMESPACE" --ignore-not-found || true
  fi
else
  echo "Note: K8s directory not found at $K8S_DIR; skipping manifest deletion."
fi

# 2. Explicit workload kinds to ensure thorough cleanup
KINDS=(deployment statefulset daemonset job cronjob service ingress configmap secret)
for kind in "${KINDS[@]}"; do
  echo "Scanning $kind objects in $NAMESPACE..."
  if kubectl get "$kind" -n "$NAMESPACE" 2>/dev/null | awk 'NR>1 {print $1}' | grep -q '.'; then
    objs=$(kubectl get "$kind" -n "$NAMESPACE" -o name 2>/dev/null || true)
    for o in $objs; do
      echo "Deleting $o"
      if $DRY_RUN; then
        echo "(dry-run) kubectl delete $o -n $NAMESPACE"
      else
        kubectl delete "$o" -n "$NAMESPACE" --ignore-not-found || true
      fi
    done
  else
    echo "(none)"
  fi
done

# 3. Optionally remove persistent volume claims & related PVs
if $DELETE_PERSISTENCE; then
  echo "Removing PersistentVolumeClaims in $NAMESPACE..."
  pvcs=$(kubectl get pvc -n "$NAMESPACE" -o name 2>/dev/null || true)
  for pvc in $pvcs; do
    echo "Deleting $pvc"
    if $DRY_RUN; then
      echo "(dry-run) kubectl delete $pvc -n $NAMESPACE"
    else
      kubectl delete "$pvc" -n "$NAMESPACE" --ignore-not-found || true
    fi
  done
  # Try delete orphan PVs labeled with namespace (if labeling policy exists)
  pvs=$(kubectl get pv -o jsonpath='{range .items[*]}{.metadata.name}{" "}{.spec.claimRef.namespace}{"\n"}{end}' 2>/dev/null | awk -v ns="$NAMESPACE" '$2==ns {print $1}' || true)
  for pv in $pvs; do
    echo "Deleting PV $pv (belongs to $NAMESPACE)"
    if $DRY_RUN; then
      echo "(dry-run) kubectl delete pv $pv"
    else
      kubectl delete pv "$pv" --ignore-not-found || true
    fi
  done
fi

# 4. Optionally delete namespace
if $DELETE_NAMESPACE; then
  echo "Deleting namespace $NAMESPACE..."
  if $DRY_RUN; then
    echo "(dry-run) kubectl delete namespace $NAMESPACE"
  else
    kubectl delete namespace "$NAMESPACE" --ignore-not-found || true
    if $WAIT_TERMINATION; then
      echo "Waiting for namespace termination..."
      while kubectl get namespace "$NAMESPACE" 2>/dev/null; do
        sleep 2
      done
      echo "Namespace $NAMESPACE terminated."
    fi
  fi
fi

# 5. Wait for residual pods if requested
if $WAIT_TERMINATION && ! $DELETE_NAMESPACE; then
  echo "Waiting for remaining pods to terminate in $NAMESPACE..."
  for i in {1..30}; do
    if ! kubectl get pods -n "$NAMESPACE" 2>/dev/null | awk 'NR>1 {exit 0} END {exit 1}'; then
      echo "All pods gone."; break
    fi
    sleep 2
  done
fi

echo "✅ Cleanup complete." 
if $DRY_RUN; then
  echo "(Dry run mode: no changes were applied)"
fi

# 6. Optionally remove local Docker images
if $REMOVE_IMAGES; then
  if ! command -v docker >/dev/null 2>&1; then
    echo "Note: docker CLI not found; skipping image removal" >&2
    exit 0
  fi
  echo "Removing local Docker images..."
  SERVICE_IMAGES=(
    "account-service:latest"
    "authentication-service:latest"
    "customer-service:latest"
    "transaction-service:latest"
  )
  THIRDPARTY_IMAGES=(
    "postgres:14"
    "apache/kafka:latest"
    "zookeeper:latest"
  )
  TARGET_IMAGES=("${SERVICE_IMAGES[@]}")
  if $REMOVE_THIRDPARTY; then
    TARGET_IMAGES+=("${THIRDPARTY_IMAGES[@]}")
  fi
  for img in "${TARGET_IMAGES[@]}"; do
    if docker image inspect "$img" >/dev/null 2>&1; then
      echo "Deleting image: $img"
      if $DRY_RUN; then
        echo "(dry-run) docker rmi -f $img"
      else
        docker rmi -f "$img" || true
      fi
    else
      echo "(skip) image not present: $img"
    fi
  done
  echo "🗑️  Image removal complete."
fi


