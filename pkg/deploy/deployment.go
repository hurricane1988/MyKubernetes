package deploy

import (
	"MyKubernetes/configs"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// 定义deploy结构体
type DeploymentStruct struct {
	ID        string `json:"id,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Replicas  int    `json:"replicas,omitempty"`
}

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
	}
}
