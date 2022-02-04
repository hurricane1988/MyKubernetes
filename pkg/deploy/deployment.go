package deploy

import (
	"MyKubernetes/configs"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"log"
	"os"
	"reflect"
)

// DeploymentStruct 定义deploy结构体
type DeploymentStruct struct {
	ID        string `json:"id,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	Replicas  int32  `json:"replicas,omitempty"`
}

// DeploymentSlice deployment切片信息
var DeploymentSlice []interface{}

// MyDeployment 定义创建deployment的结构体
type MyDeployment struct {
	Name          string `json:"name,omitempty"`
	Namespace     string `json:"namespace,omitempty"`
	MatchLabel    string `json:"matchLabel,omitempty"`
	ContainerName string `json:"containerName,omitempty"`
	PortName      string `json:"portName,omitempty"`
	ContainerPort int32  `json:"containerPort,omitempty"`
	Image         string `json:"image,omitempty"`
	Replicas      int    `json:"replicas,omitempty"`
}

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
			Replicas:  int32(int(v.Status.Replicas)),
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
func CreateDeployment(deploy MyDeployment) {
	fmt.Println(deploy)
	config, err := clientcmd.BuildConfigFromFlags("", configs.KubeconfigPath)
	if err != nil {
		panic(err.Error())
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploymentClient := client.AppsV1().Deployments(apiv1.NamespaceDefault)
	deploymentYaml := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploy.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploy.MatchLabel,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploy.MatchLabel,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  deploy.ContainerName,
							Image: deploy.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          deploy.PortName,
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: deploy.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}

	// create deployment
	fmt.Println("creating deployment...")
	result, err := deploymentClient.Create(context.TODO(), deploymentYaml, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("creating deployment %q.\n", result.GetObjectMeta().GetName())

	// List deployments
	prompt()
	fmt.Printf("Listing deployment in namespace %q\n", apiv1.NamespaceDefault)
	list, err := deploymentClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n ", d.Name, *d.Spec.Replicas)
	}

}

// DeleteDeployment 删除deployment函数
func DeleteDeployment(name, namespace string) {
	config, err := clientcmd.BuildConfigFromFlags("", configs.KubeconfigPath)
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploymentClient := client.AppsV1().Deployments(namespace)
	// delete deployment
	prompt()
	fmt.Println("deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentClient.Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Printf("删除deployment %s失败,错误信息: %e\n", name, err)
		panic(err.Error())
	}
	fmt.Printf("删除deployment %s成功!\n", name)
	log.Printf("删除deployment %s成功!\n", name)
}

// UpdateDeployment 更新deployment
func UpdateDeployment(name, namespace, image string, replicas int32) {
	config, err := clientcmd.BuildConfigFromFlags("", configs.KubeconfigPath)
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deploymentClient := client.AppsV1().Deployments(namespace)
	prompt()
	fmt.Println("updating deployment...")

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentClient.Get(context.TODO(), name, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("failed to get latest version of deployment: %v", getErr))
		}
		result.Spec.Template.Spec.Containers[0].Image = image // change nginx version
		if reflect.ValueOf(replicas).IsNil() {
			log.Printf("副本数设置为空")
			fmt.Println(replicas)
		}
		result.Spec.Replicas = int32Ptr(1) // reduce replica count
		_, updateErr := deploymentClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	fmt.Printf("updating failed: %v", retryErr)
}

//
func prompt() {
	fmt.Printf("-> Press return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}
	fmt.Println()
}

//
func int32Ptr(i int32) *int32 {
	return &i
}
