package kube

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func JobCreate(namespace string, jobSpec *batchv1.Job) {
	clientset := ClientSet()
	jobs := clientset.BatchV1().Jobs(namespace)
	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func JobCompleted(namespace string, name string) bool {
	clientset := ClientSet()
	job, err := clientset.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return job.Status.CompletionTime != nil
}

func JobFailed(namespace string, name string) bool {
	clientset := ClientSet()
	job, err := clientset.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return job.Status.Failed == *job.Spec.BackoffLimit
}

func JobDelete(namespace string, name string) {
	clientset := ClientSet()
	propagationPolicy := metav1.DeletePropagationBackground
	err := clientset.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	})
	if err != nil {
		panic(err)
	}
}
