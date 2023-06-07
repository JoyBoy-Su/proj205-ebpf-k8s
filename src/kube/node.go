package kube

import (
	"context"
	"math/rand"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// master结点的labels会有一个key为"node-role.kubernetes.io/control-plane"
var control_plane_key string = "node-role.kubernetes.io/control-plane"

func LoadNodes() []string {
	clientset := ClientSet()
	node_client := clientset.CoreV1().Nodes()
	node_list, err := node_client.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var nodes []string
	for _, node := range node_list.Items {
		// 若存在key则为控制平面，不包含在load结点中
		if _, ok := node.Labels[control_plane_key]; !ok {
			nodes = append(nodes, node.Name)
		}
	}
	return nodes
}

func LoadNodeRandom() string {
	nodes := LoadNodes()
	rand.Seed(time.Now().UnixNano())
	// rand.Intn(n): [0, n)
	index := rand.Intn(len(nodes))
	return nodes[index]
}

func NodeIsMaster() bool {
	// 获取host
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	clientset := ClientSet()
	node_client := clientset.CoreV1().Nodes()
	node, err := node_client.Get(context.TODO(), host, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	_, ok := node.Labels[control_plane_key]
	return ok
}
