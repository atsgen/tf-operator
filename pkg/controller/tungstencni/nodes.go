package tungstencni

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NodeList struct {
	MasterNodes    map[string]bool
	Nodes          map[string]bool
}

type TungstenRoleList struct {
	ControllerNodes   map[string]bool
	VppNodes          map[string]bool
}

var nodeList *NodeList

func newNodeList() (n *NodeList) {
	nl := new(NodeList)
	nl.MasterNodes = make(map[string]bool)
	nl.Nodes = make(map[string]bool)
	return nl
}

func FetchNodeList(client client.Client) (n *NodeList) {
	nodeList = newNodeList()
	nodes := &corev1.NodeList{}
	err := client.List(context.TODO(), nodes)
	if err == nil {
		for ix := range nodes.Items {
			node := nodes.Items[ix]
			var ipAddress string
			log.Info("Prabhjot found node: " + node.Name)
			addresses := node.Status.Addresses
			for _, address := range addresses {
				if address.Type == corev1.NodeInternalIP {
					ipAddress = address.Address
					log.Info("Prabhjot IP: " + address.Address)
					break
				}
			}
			labels := node.GetLabels()
			var isMaster bool
			for key, element := range labels {
				log.Info("Prabhjot label: " + key + ", value: " + element)
				if key == "node-role.kubernetes.io/master" {
					isMaster = true
				}
			}
			if ipAddress != "" {
				if isMaster {
					nodeList.MasterNodes[ipAddress] = true
				}
				nodeList.Nodes[ipAddress] = true
			}
		}
	} else {
		log.Info("Prabhjot: failed reading nodes")
	}
	for key, _ := range nodeList.MasterNodes {
		log.Info("Master node found: " + key)
	}
	for key, _ := range nodeList.Nodes {
		log.Info("Node found: " + key)
	}
	return nodeList
}
