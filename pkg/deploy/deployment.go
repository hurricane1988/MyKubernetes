package deploy

import (
	"MyKubernetes/configs"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// DeploymentStruct 定义deploy结构体
type DeploymentStruct struct {
	ID        string `json:"id,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`

	Replicas int `json:"replicas,omitempty"`
}

// DeploymentSlice deployment切片信息
var DeploymentSlice []interface{}

// GetDeployFromNamespace 获取deployment信息
func GetDeployFromNamespace() {
	config, err := clientcmd.BuildConfigFromFlags("", configs.KubeconfigPath)
	if err != nil {
		panic(err.Error())
	}
	client, _ := kubernetes.NewForConfig(config)
	deploymentList, err := client.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range deploymentList.Items {
		fmt.Printf("命名空间: %v\n deployment服务名称: %v\n 副本个数: %v\n\n", v.Namespace, v.Name, v.Status.Replicas)
		results := DeploymentStruct{
			ID:        string(v.UID),
			Namespace: v.Namespace,
			Name:      v.Name,
			Replicas:  int(v.Status.Replicas),
		}
		DeploymentSlice = append(DeploymentSlice, results)
	}
	fmt.Println(DeploymentSlice)
}
