package models

// win input table
type Win struct {
	DiscoverID string
	IP         string
	Username   string
	Password   string
	IpType     string // WMI or SSh or SNMP
	IPRange    Range
	SubnetMask Subnet
}

//use timestamp

// customer - clientID
//              org
