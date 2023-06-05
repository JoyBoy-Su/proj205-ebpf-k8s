package kube

import (
	"context"

	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func JobCreate(namespace string, jobSpec *batchv1.Job) {
	clientset := ClientSet()
	jobs := clientset.BatchV1().Jobs(namespace)
	fmt.Println("Creating Job...")
	result, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Job %q .\n", result.GetObjectMeta().GetName())
}

func JobCompleted(namespace string, name string) bool {
	clientset := ClientSet()
	job, err := clientset.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return job.Status.CompletionTime != nil
}
