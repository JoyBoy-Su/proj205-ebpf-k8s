/*
Copyright © 2023 swz@fudan.edu.cn

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"

	"fudan.edu.cn/swz/bpf/bpf"
	"fudan.edu.cn/swz/bpf/kube"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all bpf running in cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// 获取client set
		clientset := kube.ClientSet()
		// 读取BPF_HOME目录，得到所有的bpf_name
		files := bpf.ListBPF()
		// 遍历files 处理每个bpf的信息
		fmt.Println("BPF\tNODE\tSTART\tSRC")
		podclient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)
		for _, bpf_name := range files {
			pod_name, src_name := bpf.ReadBPF(bpf_name)
			pod, err := podclient.Get(context.TODO(), pod_name, metav1.GetOptions{})
			if err != nil {
				panic(err)
			}
			node := pod.Spec.NodeName
			stamp := pod.ObjectMeta.CreationTimestamp
			start := stamp.Time.Format("2006-01-02 15:04:05")
			fmt.Printf("%s\t%s\t%s\t%s", bpf_name, node, start, src_name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
