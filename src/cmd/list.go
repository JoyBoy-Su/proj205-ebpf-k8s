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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all bpf running in cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// 获取client set
		// clientset := kube.ClientSet()
		// 读取BPF_HOME目录，得到所有的bpf_name
		fmt.Println("INST\tPACKAGE\tSRC_LIST\tCEATTED")
		insts := bpf.InstList()
		var inst_info bpf.InstInfo
		for _, inst_name := range insts {
			bpf.InstInfoClear(&inst_info)
			// 遍历files 处理每个bpf的信息
			bpf.InstRead(inst_name, &inst_info)
			// 输出
			pod := kube.PodStatus("default", inst_name)
			fmt.Printf("%s\t %s\t %q\t %q\n",
				inst_name,
				bpf.InstInfoGetPackageName(&inst_info),
				bpf.InstInfoGetSrcList(&inst_info),
				pod.GetCreationTimestamp(),
			)

		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
