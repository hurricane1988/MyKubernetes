// 参考链接(https://stackoverflow.com/questions/40975307/how-to-watch-events-on-a-kubernetes-service-using-its-go-client)

package events

import (
	"MyKubernetes/configs"
	"flag"
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"
)

// KubeConfig 定义全局变量信息
var (
// KubeConfig 以交互式输入方式获kubeconfig配置文件
//KubeConfig = flag.String("kubeconfig", "./config", "KubeConfig的绝对路径")
)

// WatchEvents 获取kubernetes所有事件信息
func WatchEvents(namespace string) {
	flag.Parse()
	// 获取config配置文件
	config, err := clientcmd.BuildConfigFromFlags("", configs.KubeconfigPath)
	if err != nil {
		log.Printf("获取kubeconfig配置文件失败,错误信息: %s\n", err)
		panic(err.Error())
	}
	// 初始化clientset客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("初始化clientset失败，错误信息 %s\n", err)
		panic(err.Error())
	}

	//
	watchlist := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		string(v1.ResourceServices),
		namespace,
		fields.Everything())

	_, controller := cache.NewInformer(
		watchlist,
		&v1.Service{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				fmt.Printf("service added: %s\n", obj)
			},
			DeleteFunc: func(obj interface{}) {
				fmt.Printf("deleted service: %s\n", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				fmt.Printf("update service: %s\n", oldObj)
			},
		})
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}
