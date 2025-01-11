#!/bin/bash
export KUBECONFIG=/tmp/kube.config


set -e

# Input check
if [ $# -eq 0 ]; then
  echo "Usage: $0 <Database_URL>"
  exit 1
fi

DATABASE_URL=$1
NAMESPACE="web-app"
SECRET_NAME="web-app-secret"
TEMP_SECRET_FILE="temp-secret.yaml"
SEALED_SECRET_FILE="temp-sealed-secret.yaml"
VALUES_FILE="helm-chart/values.yaml"

# Check for yq installation
if ! command -v yq &> /dev/null; then
  echo "Error: yq is not installed."
  echo "Install yq: https://mikefarah.gitbook.io/yq/"
  exit 1
fi

# Check for kubeseal installation
if ! command -v kubeseal &> /dev/null; then
  echo "Error: kubeseal is not installed."
  echo "Install kubeseal: https://github.com/bitnami-labs/sealed-secrets"
  exit 1
fi

# Check for kubectl installation
if ! command -v kubectl &> /dev/null; then
  echo "Error: kubectl is not installed."
  echo "Install kubectl: https://kubernetes.io/docs/tasks/tools/"
  exit 1
fi

echo "[INFO] Encoding DATABASE_URL..."
ENCODED_DATABASE_URL=$(echo -n "$DATABASE_URL" | base64 -w 0)
echo "Encoded DATABASE_URL: $ENCODED_DATABASE_URL"

echo "[INFO] Creating temporary Secret YAML..."
cat <<EOF > $TEMP_SECRET_FILE
apiVersion: v1
kind: Secret
metadata:
  name: $SECRET_NAME
  namespace: $NAMESPACE
data:
  DATABASE_URL: $ENCODED_DATABASE_URL
EOF
cat $TEMP_SECRET_FILE

echo "[INFO] Encrypting the Secret using kubeseal..."
kubeseal --format=yaml --scope=namespace-wide < $TEMP_SECRET_FILE > $SEALED_SECRET_FILE
cat $SEALED_SECRET_FILE

if ! grep -q "DATABASE_URL:" $SEALED_SECRET_FILE; then
  echo "Error: No secrets found in the encrypted output. Check the input or cluster configuration."
  rm -f $TEMP_SECRET_FILE $SEALED_SECRET_FILE
  exit 1
fi

echo "[INFO] Extracting encrypted DATABASE_URL..."
ENCRYPTED_DATABASE_URL=$(grep "DATABASE_URL:" $SEALED_SECRET_FILE | awk '{print $2}')
echo "Encrypted DATABASE_URL: $ENCRYPTED_DATABASE_URL"

echo "[INFO] Updating values.yaml..."
if [ -f "$VALUES_FILE" ]; then
  cp $VALUES_FILE "${VALUES_FILE}.bak"
  yq w -i $VALUES_FILE "env.secrets.DATABASE_URL" "$ENCRYPTED_DATABASE_URL"

else
  cat <<EOF > $VALUES_FILE
env:
  SERVER_PORT: 8585
  secrets:
    DATABASE_URL: "$ENCRYPTED_DATABASE_URL"
EOF
fi

echo "[INFO] Cleaning up temporary files..."
rm -f $TEMP_SECRET_FILE $SEALED_SECRET_FILE

echo "[SUCCESS] Encrypted DATABASE_URL added to $VALUES_FILE."


helm install web-app ./helm-chart --namespace web-app --create-namespace
