package kube

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// 后台删除，不夯前台，提升删除速度
	deletePolicy metav1.DeletionPropagation = "Background"
	// 立即删除
	gracetime     int64                = 0
	deleteOptions metav1.DeleteOptions = metav1.DeleteOptions{
		PropagationPolicy:  &deletePolicy,
		GracePeriodSeconds: &gracetime,
	}
)

func PodCreate(namespace string, podSpec *apiv1.Pod) {
	clientset := ClientSet()
	pods := clientset.CoreV1().Pods(namespace)
	fmt.Println("Creating Pod...")
	pod_result, err := pods.Create(context.Background(), podSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Pod %q .\n", pod_result.GetObjectMeta().GetName())
}

func PodDelete(namespace string, pod_name string) {
	clientset := ClientSet()
	pods := clientset.CoreV1().Pods(namespace)
	pods.Delete(context.TODO(), pod_name, deleteOptions)
	fmt.Printf("Delete Pod %q .\n", pod_name)
}
