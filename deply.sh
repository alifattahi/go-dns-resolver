#!/bin/bash
export KUBECONFIG=/tmp/kube.config

helm install web-app ./helm-chart --namespace web-app --create-namespace
