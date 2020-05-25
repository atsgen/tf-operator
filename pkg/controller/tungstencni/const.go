package tungstencni

const (
	// Operator object state
	TF_OPERATOR_OBJECT_IGNORED         = "Not-Ready"
	TF_OPERATOR_OBJECT_DEPLOYED        = "Deployed"
	TF_OPERATOR_OBJECT_NOT_SUPPORTED   = "Not-Supported"

	// node roles
	NODE_ROLE_VPP              = "node-role.tungsten.io/vpp"
	NODE_ROLE_VROUTER          = "node-role.tungsten.io/agent"
	NODE_ROLE_ANALYTICS        = "node-role.tungsten.io/analytics"
	NODE_ROLE_ANALYTICS_ALARM  = "node-role.tungsten.io/analytics_alarm"
	NODE_ROLE_ANALYTICS_SNMP   = "node-role.tungsten.io/analytics_snmp"
	NODE_ROLE_ANALYTICS_DB     = "node-role.tungsten.io/analyticsdb"
	NODE_ROLE_CONFIG           = "node-role.tungsten.io/config"
	NODE_ROLE_CONFIG_DB        = "node-role.tungsten.io/configdb"
	NODE_ROLE_CONTROL          = "node-role.tungsten.io/control"
	NODE_ROLE_WEBUI            = "node-role.tungsten.io/webui"
)
