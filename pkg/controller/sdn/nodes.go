package sdn

import (
	"context"
	"errors"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NodeList struct {
	MasterNodes        map[string]string
	WorkerNodes        map[string]string
	DefultApiServer    string
	MasterNodesStr     string
	NodesStr           string
}

type TungstenRoleList struct {
	ControllerNodes   map[string]bool
	VppNodes          map[string]bool
}

var nodeList *NodeList

func newNodeList() (n *NodeList) {
	nl := new(NodeList)
	nl.MasterNodes = make(map[string]string)
	nl.WorkerNodes = make(map[string]string)
	return nl
}

func SetNodeLabels(cl client.Client, nodeName string, labels []string, dType string) error {
	node := &corev1.Node{}

	err := cl.Get(context.TODO(), client.ObjectKey{Name: nodeName}, node)
	if err != nil {
		log.Info("Failed to get node with error " + err.Error())
		return err
	}

	newNode := node.DeepCopy()
	// delete previous labels
	delete(newNode.Labels, "node-role.tungsten.io/vpp")
	delete(newNode.Labels, "node-role.tungsten.io/agent")
	delete(newNode.Labels, NodeRoleAnalytics)
	delete(newNode.Labels, NodeRoleAnalyticsAlarm)
	delete(newNode.Labels, NodeRoleAnalyticsSnmp)
	delete(newNode.Labels, NodeRoleAnalyticsDb)
	delete(newNode.Labels, NodeRoleConfig)
	delete(newNode.Labels, NodeRoleConfigDb)
	delete(newNode.Labels, NodeRoleControl)
	delete(newNode.Labels, NodeRoleWebui)

	newNode.Labels[NodeRoleDatapath] = dType
	for _,label := range labels {
		newNode.Labels[label] = ""
	}

	err = cl.Patch(context.TODO(), newNode, client.MergeFrom(node))
	if err != nil {
		log.Info("Failed to patch node with error " + err.Error())
		return err
	}


	return nil
}

func FetchNodeList(client client.Client) (*NodeList, error) {
	nodeList = newNodeList()
	nodes := &corev1.NodeList{}
	err := client.List(context.TODO(), nodes)
	if err == nil {
		for ix := range nodes.Items {
			node := nodes.Items[ix]
			var ipAddress string
			addresses := node.Status.Addresses
			for _, address := range addresses {
				if address.Type == corev1.NodeInternalIP {
					ipAddress = address.Address
					break
				}
			}
			labels := node.GetLabels()
			var isMaster bool
			for key, _ := range labels {
				if key == "node-role.kubernetes.io/master" {
					isMaster = true
				}
			}
			if ipAddress != "" {
				if isMaster {
					nodeList.MasterNodes[ipAddress] = node.Name
					if nodeList.DefultApiServer == "" {
						nodeList.DefultApiServer = ipAddress
					}
					if nodeList.MasterNodesStr == "" {
						nodeList.MasterNodesStr = ipAddress
					} else {
						nodeList.MasterNodesStr = nodeList.MasterNodesStr + "," + ipAddress
					}
				} else {
					nodeList.WorkerNodes[ipAddress] = node.Name
				}
				if nodeList.NodesStr == "" {
					nodeList.NodesStr = ipAddress
				} else {
					nodeList.NodesStr = nodeList.NodesStr + "," + ipAddress
				}
			} else {
				log.Info("Error! discovered node: " + node.Name + ", without ip")
				return nil, errors.New("node discovered without ip")
			}
		}
	} else {
		log.Info("Failed reading node information from cluster")
		return nil, err
	}
	return nodeList, nil
}
