/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"fudan.edu.cn/swz/bpf/bpf"
	"fudan.edu.cn/swz/bpf/kube"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type LogOptions struct {
	namespace  string
	name       string
	flowOption bool
}

var logOptions LogOptions

func (o *LogOptions) Validate(cmd *cobra.Command, args []string) error {

	return nil
}

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "log the status of a running bpf",
	PreRunE: func(c *cobra.Command, args []string) error {
		return logOptions.Validate(c, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("log called")
		fmt.Println(logOptions)
		inst_name := args[0]
		getPodLog(kube.ClientSet(), "default", inst_name)

	},
}

func init() {
	logCmd.Flags().StringVar(&logOptions.name, "name", "", " Ebpf name to log")
	logCmd.Flags().StringVar(&logOptions.name, "namespace", bpf.BPF_NAMESPACE, "Namespace")
	logCmd.Flags().BoolVar(&logOptions.flowOption, "flow", false, "Watch")
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getPodLog(clientset *kubernetes.Clientset, namespace, podname string) {
	// 获取最后100行日志
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

//	func tailfLog(clientset *kubernetes.Clientset,namespace,podname string) {
//		log.Println("tail pod logs ...")
//		// var (podLogs io.ReadClosererr     error)
//		// Follow true实时查看日志，同kubectl -f
//		// 也可以添加TailLines:1 从最后一行开始，默认会打印之前所有日志
//		// 如果pod中有多个container需要指定container的名字,同kubectl命令的-c参数
//		req := clientset.CoreV1().Pods(namespace).GetLogs(podname, &v1.PodLogOptions{Follow: true, Container: "hello"})
//		podLogs, err = req.Stream()
//		if err != nil {
//			return
//		}
//		defer podLogs.Close()
//		r := bufio.NewReader(podLogs)
//		// for bytes, err := r.ReadBytes(''){
//		// 	if err != nil {
//		// 		log.Println(err)
//		// 		return
//		// 	}
//		// // handler meg
//		// 	fmt.Printf("%s", string(bytes))
//		// }
//	}
//
// 根据添加的label标签来获取pod name
// func getPodName(clientset *kubernetes.Clientset, namespace string) string {
// 	// LabelSelector 使用kubectl describe 来查看label，在创建服务时设计好label
// 	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: "job-name=job-demo"})
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
// 	for _, pod := range pods.Items {
// 		return pod.GetName()
// 	}
// 	return ""
// }
