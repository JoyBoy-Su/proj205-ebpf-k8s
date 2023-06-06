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
	"fmt"

	"fudan.edu.cn/swz/bpf/bpf"
	"fudan.edu.cn/swz/bpf/kube"
	"github.com/spf13/cobra"
)

// 删除一个inst：
// 先判断inst是否存在；
// 若不存在则直接return
// 若存在则：
//
//	先把对应package的symbolic link删除；
//	再把inst的文件信息删除；
//	最后调用pod delete删除inst对应的pod
//
// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete bpf by bpfname",
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		for _, inst_name := range args {
			// 校验是否存在
			exist, err := bpf.InstExist(inst_name)
			if err != nil {
				panic(err)
			}
			if !exist {
				fmt.Printf("bpf instance '%s' not exist\n", inst_name)
				continue
			}
			// 若inst存在则执行删除
			pod_name := inst_name
			// 不处理异常反而不报错
			kube.PodDelete(bpf.BPF_NAMESPACE, pod_name)
			// 删除bpf的管理信息
			bpf.InstDelete(inst_name)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
