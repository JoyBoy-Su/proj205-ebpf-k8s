package kube

import (
	"context"
	"os/exec"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: 避免使用命令行执行
func ConfigMapCreate(name string, namespace string, from string) {
	param := "--from-file=" + from
	command := exec.Command("kubectl", "create", "cm", name, param, "-n", namespace)
	err := command.Run()
	if err != nil {
		panic(err)
	}
}

func ConfigMapDelete(name string, namespace string) {
	clientset := ClientSet()
	cms := clientset.CoreV1().ConfigMaps(namespace)
	cms.Delete(context.TODO(), name, metav1.DeleteOptions{})
}
