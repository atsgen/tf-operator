package sdn

const (
	// Operator object state

	// TFOperatorObjectIgnored - configuration object ignored
	// or not processed
	TFOperatorObjectIgnored         = "Not-Ready"
	// TFOperatorObjectPending - object waiting for resources
	TFOperatorObjectPending         = "Pending"
	// TFOperatorObjectUpdating - object being deployed
	TFOperatorObjectUpdating        = "Updating"
	// TFOperatorObjectDeployed - object deployment done
	TFOperatorObjectDeployed        = "Deployed"
	// TFOperatorNotSupported - object deployed but not supported
	TFOperatorNotSupported          = "Not-Supported"

	// datapath types

	// DatapathVpp - use vpp as datapath
	DatapathVpp               = "vpp"
	// DatapathVrouter - use vrouter as datapath
	DatapathVrouter           = "vrouter"

	// node roles
	// TODO(prabhjot) consolidate analytics roles in one label
	// consolidate config roles in one label

	// NodeRoleDatapath - role datapath
	NodeRoleDatapath         = "node-role.tungsten.io/datapath"
	// NodeRoleAnalytics - role analytics
	NodeRoleAnalytics        = "node-role.tungsten.io/analytics"
	// NodeRoleAnalyticsAlarm - role analytics alarm
	NodeRoleAnalyticsAlarm  = "node-role.tungsten.io/analytics_alarm"
	// NodeRoleAnalyticsSnmp - role analytics snmp
	NodeRoleAnalyticsSnmp   = "node-role.tungsten.io/analytics_snmp"
	// NodeRoleAnalyticsDb - role analytics db
	NodeRoleAnalyticsDb     = "node-role.tungsten.io/analyticsdb"
	// NodeRoleConfig - role config
	NodeRoleConfig           = "node-role.tungsten.io/config"
	// NodeRoleConfigDb - role config db
	NodeRoleConfigDb        = "node-role.tungsten.io/configdb"
	// NodeRoleControl - role control node
	NodeRoleControl          = "node-role.tungsten.io/control"
	// NodeRoleWebui - role webui
	NodeRoleWebui            = "node-role.tungsten.io/webui"

	// IP Forwarding values

	// IPForwardingEnabled - ip forwarding enable
	IPForwardingEnabled      = "enable"
	// IPForwardingSnat - ip forwarding enabled with SNAT
	IPForwardingSnat         = "snat"
)
