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
	"os"

	"fudan.edu.cn/swz/bpf/bpf"
	"github.com/spf13/cobra"
)

/**
 * run逻辑如下（随机选择一个node执行）：
 * 1、parse得到对应的src和name
 * 2、判断该name是否已经存在
 * 2、启动compiler-job编译得到package.json
 * 3、创建config-map挂载package.json
 * 4、启动runner-pod执行bpf程序
 * 5、在BPF_HOME下创建并进入bpfname对应的目录
 * 6、创建pod文件，文件内容是runner-pod的name
 * 7、创建src文件，文件内容是.c和.h
 * TODO:
 * 1、不要每次都重新编译，先查src，若没有再编译（需要将package.json添加到src下）即缓存机制
 * 2、去除hard code，且避免job的name重名
 */

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "kubectl run *.bpf.c *.bpf.h",
	Long:  "kubectl compile *.bpf.c *.bpf.h to pack.json, and run it on one pod",
	// 最少1个参数，
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
		// 获取src和name
		inst_name, err := cmd.Flags().GetString("bpfname")
		if err != nil {
			panic(err)
		}
		// 获取当前目录，用来挂载到compiler中
		exePath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		// 编译
		package_name := bpf.Compile(bpf.BPF_EMPTY_PACKAGE_NAME, exePath, args)
		// 创建configMap，通过cmd的方式
		bpf.MountPackageByConfigMap(package_name)
		// 运行
		bpf.Run(inst_name, package_name)
		fmt.Println("Bpf program run successfully")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("bpfname", "m", "default", "name for bpf instance")
	runCmd.MarkFlagRequired("bpfname")
}
