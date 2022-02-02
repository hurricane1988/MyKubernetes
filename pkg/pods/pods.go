package pods

import (
	"MyKubernetes/configs"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

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
	log.SetPrefix("[pods] ")
	log.SetFlags(log.Ltime | log.Ldate | log.Lmicroseconds)
}

// Pods 定义pod结构体信息
type Pods struct {
	ID        string      `json:"id,omitempty"`
	Name      string      `json:"name,omitempty"`
	IP        string      `json:"ip,omitempty"`
	HostIP    string      `json:"host_ip,omitempty"`
	StartTime interface{} `json:"start_time,omitempty"`
	Namespace string      `json:"namespace,omitempty"`
}

// GetPodsFromNamespace 获取指定namespace中的pod信息
func GetPodsFromNamespace(namespace string, all bool) {
	// 定义pod slice信息
	var PodSlice []interface{}
	// 初始化subconfig
	config, err := clientcmd.BuildConfigFromFlags("", configs.KubeconfigPath)
	_, err = clientcmd.BuildConfigFromFlags("", configs.KubeconfigPath)
	if err != nil {
		log.Fatalf("获取kubeconfig文件失败,错误信息: %e\n", err)
		return
	}
	client, _ := kubernetes.NewForConfig(config)
	if all == true {
		pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
			log.Fatalf("获取namespace %s中pod信息失败, 错误信息%e\n", namespace, err)
			return
		}
		for _, pod := range pods.Items {
			log.Printf("获取namespace %v指定POD信息成功, pod名称: %v IP: %v\n", pod.Namespace, pod.Name, pod.Status.PodIP)
			// fmt.Printf(" 命名空间是：%v\n pod名字：%v\n IP：%v\n\n", pod.Namespace, pod.Name, pod.Status.PodIP)
			results := Pods{
				ID:        string(pod.UID),
				Name:      pod.Name,
				IP:        pod.Status.PodIP,
				HostIP:    pod.Status.HostIP,
				StartTime: pod.Status.StartTime,
				Namespace: pod.Namespace,
			}
			PodSlice = append(PodSlice, results)
			fmt.Println(results)
		}
	} else if all == false {
		pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
			log.Fatalf("获取namespace %s中pod信息失败, 错误信息%e\n", namespace, err)
			return
		}
		for _, pod := range pods.Items {
			log.Printf("获取namespace %v指定POD信息成功, pod名称: %v IP: %v\n", pod.Namespace, pod.Name, pod.Status.PodIP)
			// fmt.Printf(" 命名空间是：%v\n pod名字：%v\n IP：%v\n\n", pod.Namespace, pod.Name, pod.Status.PodIP)
		}
	}
	fmt.Println(PodSlice[:])
}
