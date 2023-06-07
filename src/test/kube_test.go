package test

// func TestNodes(t *testing.T) {
// 	clientset := kube.ClientSet()
// 	node_client := clientset.CoreV1().Nodes()
// 	node_list, err := node_client.List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, node := range node_list.Items {
// 		t.Logf("%s%q\n", node.Name, node.Labels)
// 	}
// }

// func TestGetAllNodes(t *testing.T) {
// 	nodes := kube.LoadNodes()
// 	t.Logf("%q\n", nodes)
// }

// func TestRandomNode(t *testing.T) {
// 	var node string
// 	node = kube.LoadNodeRandom()
// 	t.Logf("node = %s\n", node)
// 	node = kube.LoadNodeRandom()
// 	t.Logf("node = %s\n", node)
// 	node = kube.LoadNodeRandom()
// 	t.Logf("node = %s\n", node)
// 	node = kube.LoadNodeRandom()
// 	t.Logf("node = %s\n", node)
// 	node = kube.LoadNodeRandom()
// 	t.Logf("node = %s\n", node)
// 	node = kube.LoadNodeRandom()
// 	t.Logf("node = %s\n", node)
// 	node = kube.LoadNodeRandom()
// 	t.Logf("node = %s\n", node)
// 	node = kube.LoadNodeRandom()
// 	t.Logf("node = %s\n", node)
// 	node = kube.LoadNodeRandom()
// 	t.Logf("node = %s\n", node)
// 	node = kube.LoadNodeRandom()
// 	t.Logf("node = %s\n", node)
// }

// func TestCreateConfigMap(t *testing.T) {
// 	package_name := "e2308b14-ccaa-4241-b52d-1f912a6b02db"
// 	param := "--from-file=" + bpf.BPF_PACKAGE_HOME + package_name + "/" + bpf.DATA_DIR_NAME
// 	// namespace := "-n bpf"
// 	command := exec.Command("kubectl", "create", "cm", package_name, param, "-n", "bpf")
// 	err := command.Run()
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func TestPodStatus(t *testing.T) {
// 	pod := kube.Pod("default", "nginx-app-5c64488cdf-75wxb")
// 	// t.Logf("pod.Status: %v\n", pod.Status)
// 	t.Logf("pod.Status: %v\n", pod.Status.Phase)
// }
