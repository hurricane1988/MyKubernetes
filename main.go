package main

import (
	"MyKubernetes/pkg/deploy"
)

func main() {
	// pods.GetPodsFromNamespace("kube-system", true)
	deploy.GetDeployFromNamespace()
}
