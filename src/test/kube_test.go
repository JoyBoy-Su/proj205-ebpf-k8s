package test

import (
	"context"
	"testing"

	"fudan.edu.cn/swz/bpf/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNodes(t *testing.T) {
	clientset := kube.ClientSet()
	node_client := clientset.CoreV1().Nodes()
	node_list, err := node_client.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, node := range node_list.Items {
		t.Logf("%s%q\n", node.Name, node.Labels)
	}
}

func TestGetAllNodes(t *testing.T) {
	nodes := kube.LoadNodes()
	t.Logf("%q\n", nodes)
}

func TestRandomNode(t *testing.T) {
	var node string
	node = kube.LoadNodeRandom()
	t.Logf("node = %s\n", node)
	node = kube.LoadNodeRandom()
	t.Logf("node = %s\n", node)
	node = kube.LoadNodeRandom()
	t.Logf("node = %s\n", node)
	node = kube.LoadNodeRandom()
	t.Logf("node = %s\n", node)
	node = kube.LoadNodeRandom()
	t.Logf("node = %s\n", node)
	node = kube.LoadNodeRandom()
	t.Logf("node = %s\n", node)
	node = kube.LoadNodeRandom()
	t.Logf("node = %s\n", node)
	node = kube.LoadNodeRandom()
	t.Logf("node = %s\n", node)
	node = kube.LoadNodeRandom()
	t.Logf("node = %s\n", node)
	node = kube.LoadNodeRandom()
	t.Logf("node = %s\n", node)
}
