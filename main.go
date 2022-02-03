package main

import (
	"MyKubernetes/pkg/deploy"
	"MyKubernetes/pkg/pods"
)

func main() {
	pods.GetPodsFromNamespace("kube-system", true)
	deploy.GetDeployFromNamespace()
	deploy.CreateDeployment()
}
