/*
Copyright © 2023 jiadisu@fudan.edu.cn

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
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show package details",
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		package_name := args[0]
		// 检查是否存在
		exist, err := bpf.PackageExist(package_name)
		if err != nil {
			panic(err)
		}
		if !exist {
			err := fmt.Errorf("package %s not exist", package_name)
			if err != nil {
				panic(err)
			}
			return
		}
		// 存在后查找详细信息
		var package_info bpf.PackageInfo
		bpf.PackageRead(package_name, &package_info)
		// 逐项显示
		fmt.Printf("Name: %s\n", package_name)
		fmt.Printf("InstList: %q\n", bpf.PackageInfoGetInstList(&package_info))
		fmt.Printf("SrcList: %q\n", bpf.PackageInfoGetSrcList(&package_info))
		for _, src_name := range bpf.PackageInfoGetSrcList(&package_info) {
			fmt.Printf("========== %s ==========\n", src_name)
			src_str := bpf.SrcRead(package_name, src_name)
			fmt.Printf("%s\n", src_str)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
