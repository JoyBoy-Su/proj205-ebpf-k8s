package test

import (
	"context"
	"fmt"

	"fudan.edu.cn/swz/bpf/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientSet() {
	clientset := kube.ClientSet()
	// 获取pod资源，测试clientset
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("clientset get pods error:\n")
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("There are %d pods in this cluster\n", len(pods.Items))

	namespace := "default"
	pod := "nginx-app-5c64488cdf-75wxb"
	_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n", pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		fmt.Printf("clientset get pods error:\n")
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("Found pod %s in namesapce %s\n", pod, namespace)
	}
}
