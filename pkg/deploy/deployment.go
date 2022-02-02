package deploy

import (
	"MyKubernetes/configs"
	"context"
	"encoding/json"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
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

// 初始化日志信息
func init() {
	LogPath := configs.LogPath
	LogFile, err := os.OpenFile(LogPath,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		fmt.Printf("创建日志路径%s失败, 错误信息%e\n", LogPath, err)
		return
	}
	log.SetOutput(LogFile)
	log.SetPrefix("[deployment] ")
	log.SetFlags(log.Ltime | log.Ldate | log.Lmicroseconds)
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
		// fmt.Printf("命名空间: %v\n deployment服务名称: %v\n 副本个数: %v\n\n", v.Namespace, v.Name, v.Status.Replicas)
		results := DeploymentStruct{
			ID:        string(v.UID),
			Namespace: v.Namespace,
			Name:      v.Name,
			Replicas:  int(v.Status.Replicas),
		}
		// json格式化处理
		results01, _ := json.Marshal(results)
		results02, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(results01))
		fmt.Println("\n", string(results02))
		DeploymentSlice = append(DeploymentSlice, string(results02))
	}
	fmt.Println(DeploymentSlice)
}

// CreateDeployment 创建deployment方法
func CreateDeployment() {
	config, err := clientcmd.BuildConfigFromFlags("", configs.KubeconfigPath)
	if err != nil {
		panic(err.Error())
	}
	client, _ := kubernetes.NewForConfig(config)

	deploymentClient := client.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-test",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http-80",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
}
