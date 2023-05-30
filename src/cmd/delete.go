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
	"os"

	"fudan.edu.cn/swz/bpf/bpf"
	"fudan.edu.cn/swz/bpf/kube"
	"github.com/spf13/cobra"
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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete bpf by bpfname",
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		clientset := kube.ClientSet()
		bpf_name := args[0]
		var dir string = bpf.BPF_HOME + bpf_name
		pod_name, err := os.ReadFile(dir + "/" + bpf.POD_FILE)
		var podName string = string(pod_name)
		if err != nil {
			panic(err)
		}
		// 不处理异常反而不报错
		clientset.CoreV1().Pods(bpf.BPF_NAMESPACE).Delete(context.TODO(), podName, deleteOptions)
		// 删除文件夹
		os.RemoveAll(dir)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
