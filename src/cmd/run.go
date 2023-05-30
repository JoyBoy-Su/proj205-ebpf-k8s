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
	"os/exec"

	"fmt"
	"os"

	"fudan.edu.cn/swz/bpf/bpf"
	"fudan.edu.cn/swz/bpf/kube"
	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func parse(cmd *cobra.Command, args []string) ([]string, string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var src []string
	for _, value := range args {
		src = append(src, cwd + "/" + value)
	}
	name, err := cmd.Flags().GetString("bpfname")
	if err != nil {
		panic(err)
	}
	return src, name
}

func exist(pathname string) (bool, error) {
	_, err := os.Stat(pathname)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

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
		bpf_src, bpf_name := parse(cmd, args)
		// 校验name是否已存在
		exist, err := exist(bpf.BPF_HOME + bpf_name)
		if err != nil {
			panic(err)
		}
		if exist {
			fmt.Println("Name already exists")
			return
		}
		// 获取client set，与k8s集群交互
		clientset := kube.ClientSet()
		// 获取当前目录，用来挂载到compiler中
		exePath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		// 设置Job的资源清单
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
				Name:      "compile-bpf-job",
				Namespace: bpf.BPF_NAMESPACE,
			},
			Spec: batchv1.JobSpec{
				Completions:             &completions,
				TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "compile-bpf",
						Namespace: bpf.BPF_NAMESPACE,
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
		// 创建Job
		fmt.Println("Creating Job...")
		result, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created Job %q   and compile.\n", result.GetObjectMeta().GetName())

		// 创建configMap，通过cmd的方式
		param := "--from-file=" + exePath
		command := exec.Command("kubectl", "create", "cm", "bpf-package", param)
		command.Run()
		// 设置runner-pod的资源清单
		run_pod := clientset.CoreV1().Pods(apiv1.NamespaceDefault)
		var allowPrivilegeEscalation bool = true
		var privileged bool = true
		run_pod_Spec := &apiv1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "bpf-runner-pod",
			},
			Spec: apiv1.PodSpec{
				Containers: []apiv1.Container{
					{
						Name:    "bpf-runner-container",
						Image:   "ngccc/ecli_x86_ubuntu",
						Command: []string{"/bin/sh", "-c", "./ecli run /var/ebpfPackage/package.json"},
						VolumeMounts: []apiv1.VolumeMount{
							{
								Name:      "logs",
								MountPath: "/sys/kernel/debug",
							},
							{
								Name:      "bpf-package",
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
						Name: "bpf-package",
						VolumeSource: apiv1.VolumeSource{
							ConfigMap: &apiv1.ConfigMapVolumeSource{
								LocalObjectReference: apiv1.LocalObjectReference{
									Name: "bpf-package",
								},
							},
						},
					},
				},
			},
		}
		// 创建Pod
		fmt.Println("Creating Pod...")
		pod_result, err := run_pod.Create(context.Background(), run_pod_Spec, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created Pod %q   and run.\n", pod_result.GetObjectMeta().GetName())
		// 获取Pod的name
		pod_name := pod_result.GetName()
		// 创建bpf_name目录，并设置pod文件和src文件
		bpf.AddBPF(bpf_name, pod_name, bpf_src)
		fmt.Println("Bpf program run successfully")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("bpfname", "m", "default", "name for bpf instance")
	runCmd.MarkFlagRequired("bpfname")
}
