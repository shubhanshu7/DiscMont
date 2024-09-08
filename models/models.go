package models

import "time"

type Devicedetails struct {
	DiscoverID         string
	DeviceName         string
	DeviceStatus       string
	Device_DNS_Name    string
	DeviceDesc         string
	DeviceType         string
	Device_MacAddress  string
	DevicePrimaryIP    string
	Device_long_desc   string
	Device_Manf_Name   string
	DeviceModelName    string
	DeviceSerialNumber string // only for snmp(maybe)
	Device_soft_rev    string
	Device_firm_rev    string
	DeviceOS           string
	DeviceOS_version   string
	Device_Total_Port  string
	Device_Free_port   string
	Device_no_of_cpus  string //how many are in use or free
	Numberofhosts      string
	Memory             string //total cpu // how much is occupied
	DiscStorage        string
	PollingMethod      string
	DevOid             string // only for snmp
	Createdat          time.Time
	Updatedat          time.Time
}

// Range used for input tables
type Range struct {
	StartIP string
	EndIP   string
}

type Subnet struct {
	StartIP string
	Limit   string
}

// OIDMapping represents the structure of the OID document in MongoDB
type OIDMapping struct {
	DeviceType  string `bson:"deviceType"`
	Field       string `bson:"field"`
	Description string `bson:"description"`
	OID         string `bson:"oid"`
}
type IP struct {
	IP string
}
