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

// 显示所有的package
// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "list all packages",
	Run: func(cmd *cobra.Command, args []string) {
		packages := bpf.PackageList()
		fmt.Printf("NAME\tSRC_LIST\tINST_LIST\tSIZE\n")
		var package_info bpf.PackageInfo
		for _, package_name := range packages {
			bpf.PackageRead(package_name, &package_info)
			fmt.Printf("%s\t%q\t%q\t%d\n", package_name,
				bpf.PackageInfoGetSrcList(&package_info),
				bpf.PackageInfoGetInstList(&package_info),
				bpf.PackageInfoGetSize(&package_info))
		}
	},
}

func init() {
	rootCmd.AddCommand(packagesCmd)
}
