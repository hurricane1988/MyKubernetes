package main

import (
	"MyKubernetes/pkg/pods"
)

func main() {
	pods.GetPodsFromNamespace("kube-system", true)
}
