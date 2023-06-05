package kube

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
