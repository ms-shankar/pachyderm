#!/bin/sh

set -Eex
# Parse flags
VERSION=v1.8.0
while getopts ":v:" opt; do
  case "${opt}" in
    v)
      VERSION="v${OPTARG}"
      ;;
    \?)
      echo "Invalid argument: ${opt}"
      exit 1
      ;;
  esac
done

#sudo CHANGE_MINIKUBE_NONE_USER=true minikube start --vm-driver=none --kubernetes-version="${VERSION}" --extra-config=kubelet.CertDir=/var/lib/kubelet/pki --extra-config=apiserver.ServiceClusterIpRange=10.0.0.1/12
#sudo CHANGE_MINIKUBE_NONE_USER=true minikube start --vm-driver=none --kubernetes-version="${VERSION}" --extra-config=kubelet.ClusterDNS=10.0.0.1 --extra-config=kubelet.Address=10.0.0.1
sudo CHANGE_MINIKUBE_NONE_USER=true minikube start --vm-driver=none --kubernetes-version="${VERSION}" --extra-config=kubelet.ResolvConf=/etc/resolv.conf
until kubectl version 2>/dev/null >/dev/null; do sleep 5; done
kubectl version