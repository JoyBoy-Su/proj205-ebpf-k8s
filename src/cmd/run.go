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
	"os"

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
// "os"
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
		exePath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		// 创建Job
		jobs := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
		var completions int32 = 1
		var hostpathdirectory apiv1.HostPathType = apiv1.HostPathDirectory
		var ttlSecondsAfterFinished int32 = 5
		jobSpec := &batchv1.Job{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Job",
				APIVersion: "batch/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "compile-bpf-job",
			},
			Spec: batchv1.JobSpec{
				Completions:             &completions,
				TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
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
						NodeName: "master",
						Volumes: []apiv1.Volume{
							{
								Name: "bpf-src-volume",
								VolumeSource: apiv1.VolumeSource{
									HostPath: &apiv1.HostPathVolumeSource{
										Path: exePath,
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
		// Create Job
		fmt.Println("Creating Job...")
		result, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created Job %q   and compile.\n", result.GetObjectMeta().GetName())

		// // 创建configMap
		// var namespaceName string = "default"
		// configMap := &apiv1.ConfigMap{}
		// configMapInfo, err := clientset.CoreV1().ConfigMaps(namespaceName).Create(context.TODO(), configMap, metav1.CreateOptions{})
		// if err != nil {
		// 	panic(err)
		// }
		// 创建Pod
		run_pod := clientset.CoreV1().Pods(apiv1.NamespaceDefault)
		var allowPrivilegeEscalation bool = true
		var privileged bool = true
		run_pod_Spec := &apiv1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "ecli-x86-ubuntu-pod",
			},
			Spec: apiv1.PodSpec{
				Containers: []apiv1.Container{
					{
						Name:    "ecli-x86-ubuntu",
						Image:   "ngccc/ecli_x86_ubuntu",
						Command: []string{"/bin/sh", "-c", "./ecli run /var/ebpfPackage/package.json"},
						VolumeMounts: []apiv1.VolumeMount{
							{
								Name:      "logs",
								MountPath: "/sys/kernel/debug",
							},
							{
								Name:      "config-vol",
								MountPath: "/var/ebpfPackage/",
							},
						},
						SecurityContext: &apiv1.SecurityContext{
							AllowPrivilegeEscalation: &allowPrivilegeEscalation,
							Privileged:               &privileged,
						},
					},
				},
				Volumes: []apiv1.Volume{
					{
						Name: "logs",
						VolumeSource: apiv1.VolumeSource{
							HostPath: &apiv1.HostPathVolumeSource{
								Path: "/sys/kernel/debug",
								Type: &hostpathdirectory,
							},
						},
					},
					{
						Name: "config-vol",
						VolumeSource: apiv1.VolumeSource{
							ConfigMap: &apiv1.ConfigMapVolumeSource{
								LocalObjectReference: apiv1.LocalObjectReference{
									Name: "ebpf-config",
								},
							},
						},
					},
				},
			},
		}
		// Create Pod
		fmt.Println("Creating Pod...")
		pod_result, err := run_pod.Create(context.Background(), run_pod_Spec, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created Pod %q   and run.\n", pod_result.GetObjectMeta().GetName())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
