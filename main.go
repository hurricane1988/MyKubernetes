package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/admin/.kube/config")
	if err != nil {
		panic(err)
	}
	client, _ := kubernetes.NewForConfig(config)
	pods, err := client.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range pods.Items {
		fmt.Printf(" 命名空间是：%v\n pod名字：%v\n IP：%v\n\n", v.Namespace, v.Name, v.Status.PodIP)
	}
}
