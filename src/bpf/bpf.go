package bpf

import (
	"os/exec"

	"fmt"

	"fudan.edu.cn/swz/bpf/kube"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 执行bpf编译的过程：
// 创建package_name；
// 复制src到package_name；
// 创建Job编译；
// 复制package.json到package_name；
func Compile(package_name string, base_path string, files_list []string) string {
	// 创建package_name
	package_name = packageName(package_name)
	PackageCreate(package_name)
	// 复制src到package
	PackageAddSrcList(package_name, base_path, files_list)
	// 设置Job的资源清单
	var completions int32 = 1
	var hostpathdirectory apiv1.HostPathType = apiv1.HostPathDirectory
	var ttlSecondsAfterFinished int32 = 5
	var name = package_name
	jobSpec := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: BPF_NAMESPACE,
		},
		Spec: batchv1.JobSpec{
			Completions:             &completions,
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: BPF_NAMESPACE,
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
							Name: "bpf-data",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: BPF_PACKAGE_HOME + package_name + "/" + DATA_DIR_NAME,
									Type: &hostpathdirectory,
								},
							},
						},
					},
					Containers: []apiv1.Container{
						{
							Name:  "bpf-compiler",
							Image: CompileImage,
							Args:  files_list,
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "bpf-data",
									MountPath: CompileMountPath,
								},
							},
						},
					},
				},
			},
		},
	}
	// 创建Job编译
	kube.JobCreate(BPF_NAMESPACE, jobSpec)
	for {
		if kube.JobCompleted(BPF_NAMESPACE, package_name) {
			break
		}
	}
	fmt.Println("compile completed")
	return package_name
}

// 将package的data部分挂载到config map以供后续运行
// cm的名称与package的名称一致
// TODO: 避免使用命令行执行
func MountPackageByConfigMap(package_name string) {
	fmt.Println("mount config-map")
	param := "--from-file=" + BPF_PACKAGE_HOME + package_name + "/" + DATA_DIR_NAME
	command := exec.Command("kubectl", "create", "cm", package_name, param)
	err := command.Run()
	if err != nil {
		panic(err)
	}
}

// 从package中创建bpf instance到指定的node
func Run(inst_name string, package_name string, node string, serial bool) {
	// 创建bpf instance
	inst_name = instanceName(inst_name, serial)
	InstAdd(inst_name, package_name)
	var hostpathdirectory apiv1.HostPathType = apiv1.HostPathDirectory
	// 设置runner-pod的资源清单
	var allowPrivilegeEscalation bool = true
	var privileged bool = true
	podSpec := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: inst_name,
		},
		Spec: apiv1.PodSpec{
			NodeName: node,
			Containers: []apiv1.Container{
				{
					Name:    "bpf-runner",
					Image:   RunImage,
					Command: RunCommand,
					VolumeMounts: []apiv1.VolumeMount{
						{
							Name:      "logs",
							MountPath: "/sys/kernel/debug",
						},
						{
							Name:      "bpf-package",
							MountPath: RunMountPath,
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
								Name: package_name,
							},
						},
					},
				},
			},
		},
	}
	// 启动pod
	kube.PodCreate(BPF_NAMESPACE, podSpec)
}
