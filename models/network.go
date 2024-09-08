package models

import "time"

// Network input table
type SNMPInput struct {
	DiscoveryID    string
	IP             string
	IPType         string
	IPRange        Range
	SnmpParameters SnmpParams
	DeviceProvider string
	SubnetMask     Subnet
}
type SnmpParams struct {
	Version          string        // snmp version
	Timeout          time.Duration // request timeout (default = 5sec)
	Retries          uint          //default = 0
	Community        string        // V1 or V2 specific
	Username         string        //V3 specific (security name)
	SecurityLevel    string        //V3 specific
	AuthPassword     string        // Authentication password(V3 specific)
	AuthProtocol     string        //Authentication protocol(V3 specific)
	PrivPassword     string        //Privacy password (V3 specific)
	PrivProtocol     string        // Privacy protocol(V3 specific)
	SecurityEngineID string        //(V3 specific)
	ContextEngineID  string        //(V3 specific)
	ContextName      string        //(V3 specific)
}
