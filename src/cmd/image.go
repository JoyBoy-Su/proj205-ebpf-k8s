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
	"context"
	"fmt"

	"fudan.edu.cn/swz/bpf/kube"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "A command to view the deployment management image",
	Run: func(cmd *cobra.Command, args []string) {
		clientset := kube.ClientSet()
		namespace, _ := rootCmd.PersistentFlags().GetString("namespace")
		deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			fmt.Printf("clientset.Deployments error:\n")
			fmt.Printf("err: %v\n", err)
		}
		for _, deployment := range deployments.Items {
			for _, container := range deployment.Spec.Template.Spec.Containers {
				fmt.Printf("NameSpace: %s\t", namespace)
				fmt.Printf("DeploymentName: %s\t", deployment.GetName())
				fmt.Printf("ContainerName: %s\t", container.Name)
				fmt.Printf("Image: %s\n", container.Image)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
}
