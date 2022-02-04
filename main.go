package main

import "MyKubernetes/pkg/deploy"

var mydeployment deploy.MyDeployment

func main() {
	//pods.GetPodsFromNamespace("kube-system", true)
	//deploy.GetDeployFromNamespace()
	deploy.DeleteDeployment("web", "default")
	//mydeployment.Namespace = "default"
	//mydeployment.Name = "web"
	//mydeployment.PortName = "tcp-80"
	//mydeployment.ContainerName = "web"
	//mydeployment.ContainerPort = 80
	//mydeployment.Replicas = 2
	//mydeployment.Image = "nginx:1.19.3"
	//mydeployment.MatchLabel = "nginx"
	//mydeployment.ContainerName = "web"
	//deploy.CreateDeployment(mydeployment)
	//deploy.UpdateDeployment("web", "default", "nginx:1.19.3", 3)

}
