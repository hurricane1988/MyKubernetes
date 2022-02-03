package main

import "MyKubernetes/pkg/deploy"

var mydeployment deploy.MyDeployment

func main() {
	//pods.GetPodsFromNamespace("kube-system", true)
	//deploy.GetDeployFromNamespace()
	//deploy.DeleteDeployment("nginx-test")
	mydeployment.Namespace = "default"
	mydeployment.Name = "web"
	mydeployment.PortName = "tcp-80"
	mydeployment.ContainerName = "web"
	mydeployment.ContainerPort = 80
	mydeployment.Replicas = 2
	mydeployment.Image = "nginx:1.19.3"
	mydeployment.MatchLabel = map[string]string{
		"app": "nginx",
	}
	mydeployment.ContainerName = "web"
	deploy.CreateDeployment(mydeployment)

}
