#!/bin/bash
export KUBECONFIG=/tmp/kube.config
kubectl apply -f project.yaml
kubectl apply -f application.yaml