package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
)

func main() {
	// 初始化kubeconfig
	var kubeconfig *string
	home := homedir.HomeDir()
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// 获取config 与 client set
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	// 获取deployments client，以此为入口访问deployment
	depolymentsClient := clientset.AppsV1().Depolyments(apiv1.NamespaceDefault)

	// 初始化一个deployment （按照yaml资源清单的方式）
	depolyment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-depolyment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Containers{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	// 调用接口创建deployment
	fmt.Println("creating deployment")
	result, err := depolymentsClient.Create(context.TODO(), depolyment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created depolyment %q.\n", result.GetObjectMeta().GetName())

	// 调用接口更新deployment
	prompt()
	fmt.Println("update deployment")
	// 更新deployment的两种方式
	// 1、更新deployment变量，修改其信息，并调用Update(deployment)
	// 2、调用get得到result，修改result（deployment的另一种形式），retry
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := depolymentsClient.Get(context.TODO(), "demo-depolyment", metav1.CetOptions{})
		if getErr != nil {
			panic(getErr)
		}
		result.Spec.Replicas = int32Ptr(1)                      // 修改replicas
		result.Spec.Template.Containers[0].Image = "nginx:1.13" // 修改版本
		_, updateErr := depolymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(retryErr)
	}
	fmt.Println("Updated deployment")

	// 调用接口，列出所有的deployment
	prompt()
	fmt.Printf("Listing depolyments in namespace %q. \n", apiv1.NamespaceDefault)
	list, err := depolymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// 遍历deployment
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas) \n", d.Name, *d.Spec.Replicas)
	}

	// 调用接口删除deployment
	prompt()
	fmt.Println("delete deployment")
	deletePolicy := metav1.DeletePropagationForeground
	if err := depolymentsClient.Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment")
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func int32Ptr(i int32) *int32 { return &i }
