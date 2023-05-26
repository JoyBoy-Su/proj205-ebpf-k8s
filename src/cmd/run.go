/*
Copyright © 2023 wzh@fudan.edu.cn

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
	"flag"
	"path/filepath"

	"fmt"

	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// "context"
// 	"flag"
// 	"os/exec"
// 	"path/filepath"

// 	"fmt"

//	"github.com/spf13/cobra"
//	appsv1 "k8s.io/api/apps/v1"
//	apiv1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/tools/clientcmd"
//	"k8s.io/client-go/util/homedir"
// batchv1 "k8s.io/api/batch/v1"

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "kubectl run *.bpf.c *.bpf.h",
	Long:  "kubectl compile *.bpf.c *.bpf.h to pack.json, and run it on one pod",
	// 最少1个参数，
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		jobs := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
		var completions int32 = 1
		var hostpathdirectory apiv1.HostPathType = apiv1.HostPathDirectory
		jobSpec := &batchv1.Job{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Job",
				APIVersion: "batch/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "compile-bpf-job",
			},
			Spec: batchv1.JobSpec{
				Completions: &completions,
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name: "compile-bpf",
					},
					Spec: apiv1.PodSpec{
						RestartPolicy: apiv1.RestartPolicyOnFailure,
						Tolerations: []apiv1.Toleration{
							{
								Key:    "node-role.kubernetes.io/master",
								Effect: apiv1.TaintEffectNoSchedule,
							},
						},
						NodeName: "master ",
						Volumes: []apiv1.Volume{
							{
								Name: "bpf-src-volume",
								VolumeSource: apiv1.VolumeSource{
									HostPath: &apiv1.HostPathVolumeSource{
										Path: "/home/ubuntu/jiadisu/bpf",
										Type: &hostpathdirectory,
									},
								},
							},
						},
						Containers: []apiv1.Container{
							{
								Name:  "compile-bpf",
								Image: "jiadisu/ecc-min-ubuntu-x86:0.1",
								Args:  args,
								VolumeMounts: []apiv1.VolumeMount{
									{
										Name:      "bpf-src-volume",
										MountPath: "/code",
									},
								},
							},
						},
					},
				},
			},
		}

		result, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
		_ = result
		if err != nil {
			panic(err)
		}
		// // Create Deployment
		// fmt.Println("Creating deployment...")
		// result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
		// // 编译
		// compiler_cmd := exec.Command("ecc", args...)
		// output, err := compiler_cmd.Output()
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Print(string(output))
		// // 创建configMap
		// var namespaceName string = "default"
		// configMap := &apiv1.ConfigMap{}
		// configMapInfo, err := clientset.CoreV1().ConfigMaps(namespaceName).Create(context.TODO(), configMap, metav1.CreateOptions{})
		// if err != nil {
		// 	panic(err)
		// }

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
