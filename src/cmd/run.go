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
	"errors"
	"fmt"
	"os"
	"strings"

	"fudan.edu.cn/swz/bpf/bpf"
	"fudan.edu.cn/swz/bpf/kube"
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

func compile(args []string) string {
	// 获取当前目录，用来挂载到compiler中
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// 编译
	return bpf.Compile(bpf.BPF_EMPTY_PACKAGE_NAME, exePath, args)
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "kubectl run *.bpf.c *.bpf.h",
	Long:  "kubectl compile *.bpf.c *.bpf.h to pack.json, and run it on one pod",
	// 最少1个参数，
	Args: func(cmd *cobra.Command, args []string) error {
		// 若指定了package直接通过校验
		package_name, err := cmd.Flags().GetString("package")
		if err == nil && strings.Compare(package_name, bpf.BPF_EMPTY_PACKAGE_NAME) != 0 {
			return nil
		}
		// 若未指定package则校验args的个数
		if len(args) < 1 {
			return errors.New("requires at least one arg to specifies the source code to be compiled")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 获取name
		inst_name, err := cmd.Flags().GetString("inst")
		if err != nil {
			panic(err)
		}
		// 若未指定package name则重新编译
		package_name, err := cmd.Flags().GetString("package")
		if err != nil {
			panic(err)
		}
		if strings.Compare(package_name, bpf.BPF_EMPTY_PACKAGE_NAME) == 0 {
			package_name = compile(args)
			// 创建configMap
			bpf.MountPackageByConfigMap(package_name)
		}
		// 运行
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			panic(err)
		}
		node, err := cmd.Flags().GetString("node")
		if err != nil {
			panic(err)
		}
		if all {
			// 依次启动
			for _, node := range kube.LoadNodes() {
				bpf.Run(inst_name, package_name, node, true)
			}
		} else if strings.Compare(node, bpf.BPF_EMPTY_NODE_NAME) != 0 {
			fmt.Println("node")
			bpf.Run(inst_name, package_name, node, false)
		} else {
			bpf.Run(inst_name, package_name, kube.LoadNodeRandom(), false)
		}
		fmt.Println("Bpf program run successfully")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	// 指定inst的base name
	runCmd.Flags().StringP("inst", "i", bpf.BPF_EMPTY_INSTANCE_NAME, "name for bpf instance")
	// 指定使用的package
	runCmd.Flags().StringP("package", "p", bpf.BPF_EMPTY_PACKAGE_NAME, "name for bpf package")
	// 指定是否分发到所有load node
	runCmd.Flags().BoolP("all", "a", false, "distribute bpf to all load nodes")
	// 指定是否分发到特定的node
	runCmd.Flags().StringP("node", "d", bpf.BPF_EMPTY_NODE_NAME, "specifies the node to run")
	// 分发到所有和特定node不能同时执行
	runCmd.MarkFlagsMutuallyExclusive("all", "node")
}
