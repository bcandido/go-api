#!/usr/bin/env bash

if [[ "$1" == "clean" ]]; then
    echo "==== cleaning up previous environment"
    kubectl delete -f deploy/kubernetes/
    sleep 2
    echo -e "==== finish clean up\n"
fi

echo "==== deploying new environment"
kubectl apply -f deploy/kubernetes/
echo "==== finish deploy"
