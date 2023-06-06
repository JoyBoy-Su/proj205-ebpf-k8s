package kube

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
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
	_, err := pods.Create(context.Background(), podSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func PodDelete(namespace string, pod_name string) {
	clientset := ClientSet()
	pods := clientset.CoreV1().Pods(namespace)
	pods.Delete(context.TODO(), pod_name, deleteOptions)
}

func GetPodLog(namespace, podname string) {
	// 获取最后100行日志
	clientset := ClientSet()
	var lines int64 = 100 // 如果不指定TailLines的话，会获取pod从运行到当前的所有日志
	req := clientset.CoreV1().Pods(namespace).GetLogs(podname, &v1.PodLogOptions{TailLines: &lines})
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	defer podLogs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		fmt.Println(err)
	}
	str := buf.String()
	// 处理日志
	fmt.Println(str)
}

func FllowLog(namespace, podname string) {
	// var (podLogs io.ReadClosererr     error)
	// Follow true实时查看日志，同kubectl -f
	// 也可以添加TailLines:1 从最后一行开始，默认会打印之前所有日志
	// 如果pod中有多个container需要指定container的名字,同kubectl命令的-c参数
	clientset := ClientSet()
	req := clientset.CoreV1().Pods(namespace).GetLogs(podname, &v1.PodLogOptions{Follow: true})
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return
	}
	defer podLogs.Close()
	r := bufio.NewReader(podLogs)
	outputBytes := make([]byte, 200)
	for {
		bytes, err := r.Read(outputBytes)
		if err != nil {
			log.Println(err)
			return
		} // handler meg
		fmt.Printf("%s", string(outputBytes[0:bytes]))
	}
}

func PodCreateTime(namespace string, name string) metav1.Time {
	clientset := ClientSet()
	pod, _ := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return pod.GetCreationTimestamp()
}

func PodStatus(namespace, podname string) string {
	clientset := ClientSet()
	pod, _ := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
	return string(pod.Status.Phase)
}

func PodFailed(namespace string, name string) bool {
	clientset := ClientSet()
	pod, _ := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return strings.Compare(string(pod.Status.Phase), "Failed") == 0
}

func PodRunning(namespace string, name string) bool {
	clientset := ClientSet()
	pod, _ := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return strings.Compare(string(pod.Status.Phase), "Running") == 0
}
