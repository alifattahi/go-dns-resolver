#!/bin/bash
export KUBECONFIG=/tmp/kube.config

set -e

# Check for necessary tools
if ! command -v kubeseal &> /dev/null; then
  echo "Error: kubeseal is not installed."
  echo "Install kubeseal: https://github.com/bitnami-labs/sealed-secrets"
  exit 1
fi

if ! command -v kubectl &> /dev/null; then
  echo "Error: kubectl is not installed."
  echo "Install kubectl: https://kubernetes.io/docs/tasks/tools/"
  exit 1
fi

# Check for yq installation
if ! command -v yq &> /dev/null; then
  echo "Error: yq is not installed."
  echo "Install yq: https://mikefarah.gitbook.io/yq/"
  exit 1
fi

# Prompt user for inputs
read -p "Enter GitHub Username: " GIT_USERNAME
read -s -p "Enter GitHub Password: " GIT_PASSWORD
echo


NAMESPACE="argocd"
SECRET_NAME="github-web-app"
TEMP_SECRET_FILE="temp-secret.yaml"
SEALED_SECRET_FILE="github-sealedsecret.yaml"

# Create a temporary Kubernetes Secret manifest
echo "[INFO] Creating temporary Secret YAML..."
cat <<EOF > $TEMP_SECRET_FILE
apiVersion: v1
kind: Secret
metadata:
  name: $SECRET_NAME
  namespace: $NAMESPACE
type: Opaque
stringData:
  url: https://github.com/alifattahi/go-dns-resolver.git
  username: $GIT_USERNAME
  password: $GIT_PASSWORD
EOF

echo "[INFO] Secret manifest created:"
cat $TEMP_SECRET_FILE

# Encrypt the Secret using kubeseal
echo "[INFO] Encrypting the Secret using kubeseal..."
kubeseal --format=yaml --scope=namespace-wide < $TEMP_SECRET_FILE > $SEALED_SECRET_FILE

# Verify the generated SealedSecret
if [ $? -eq 0 ]; then
  echo "[SUCCESS] SealedSecret created at $SEALED_SECRET_FILE"
else
  echo "[ERROR] Failed to create SealedSecret."
  rm -f $TEMP_SECRET_FILE
  exit 1
fi

# Display the SealedSecret
cat $SEALED_SECRET_FILE

# Cleanup temporary file
echo "[INFO] Cleaning up temporary files..."
rm -f $TEMP_SECRET_FILE

echo "[SUCCESS] SealedSecret is ready. Apply it with:"
echo "kubectl apply -f $SEALED_SECRET_FILE"



kubectl apply -f $SEALED_SECRET_FILE
kubectl apply -f project.yaml
kubectl apply -f application.yaml