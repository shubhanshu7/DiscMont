package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OIDMapping represents the structure of the OID document in MongoDB
type OIDMapping struct {
	DeviceType  string `bson:"deviceType"`
	Field       string `bson:"field"`
	Description string `bson:"description"`
	OID         string `bson:"oid"`
}

// Sample OID data
var oidData = []OIDMapping{
	{DeviceType: "cisco", Field: "DeviceName", Description: "System Name (hostname)", OID: "1.3.6.1.2.1.1.5.0"},
	{DeviceType: "cisco", Field: "DeviceStatus", Description: "System Uptime", OID: "1.3.6.1.2.1.1.3.0"},
	{DeviceType: "cisco", Field: "Device_DNS_Name", Description: "Primary DNS", OID: "1.3.6.1.2.1.4.20.1.2"},
	{DeviceType: "cisco", Field: "DeviceDesc", Description: "System Description", OID: "1.3.6.1.2.1.1.1.0"},
	{DeviceType: "cisco", Field: "DeviceType", Description: "System Object ID", OID: "1.3.6.1.2.1.1.2.0"},
	{DeviceType: "cisco", Field: "Device_MacAddress", Description: "MAC Address", OID: "1.3.6.1.2.1.2.2.1.6.0"},
	{DeviceType: "cisco", Field: "DevicePrimaryIP", Description: "IP Address", OID: "1.3.6.1.2.1.4.20.1.1.0"},
	{DeviceType: "cisco", Field: "Device_long_desc", Description: "Detailed System Description", OID: "1.3.6.1.2.1.1.1.0"},
	{DeviceType: "cisco", Field: "Device_Manf_Name", Description: "Manufacturer Name", OID: "1.3.6.1.4.1.9.1.0"},
	{DeviceType: "cisco", Field: "DeviceModelName", Description: "Model Name", OID: "1.3.6.1.2.1.1.1.0"},
	{DeviceType: "cisco", Field: "DeviceSerialNumber", Description: "Serial Number", OID: "1.3.6.1.2.1.47.1.1.1.1.11.1"},
	{DeviceType: "cisco", Field: "Device_soft_rev", Description: "Software Revision", OID: "1.3.6.1.2.1.47.1.1.1.1.10.1"},
	{DeviceType: "cisco", Field: "Device_firm_rev", Description: "Firmware Revision", OID: "1.3.6.1.2.1.47.1.1.1.1.9.1"},
	{DeviceType: "cisco", Field: "DeviceOS", Description: "Operating System", OID: "1.3.6.1.2.1.1.1.0"},
	{DeviceType: "cisco", Field: "DeviceOS_version", Description: "Operating System Version", OID: "1.3.6.1.2.1.47.1.1.1.1.10.1"},
	{DeviceType: "cisco", Field: "Device_Total_Port", Description: "Total Number of Ports", OID: "1.3.6.1.2.1.2.1.0"},
	{DeviceType: "cisco", Field: "Device_Free_port", Description: "Available Ports", OID: "1.3.6.1.2.1.2.2.1.8.0"},
	{DeviceType: "cisco", Field: "Device_no_of_cpus", Description: "Number of CPUs", OID: "1.3.6.1.4.1.9.2.1.57.0"},
	{DeviceType: "cisco", Field: "Numberofhosts", Description: "Number of Connected Hosts", OID: "1.3.6.1.2.1.4.9.0"},
	{DeviceType: "cisco", Field: "Memory", Description: "Total Memory", OID: "1.3.6.1.4.1.9.9.48.1.1.1.5.1"},
	{DeviceType: "cisco", Field: "DiscStorage", Description: "Total Disk Storage", OID: "1.3.6.1.4.1.9.9.48.1.1.1.6.1"},
	{DeviceType: "cisco", Field: "PollingMethod", Description: "Polling Method", OID: "1.3.6.1.2.1.4.20.1.1.0"},
	{DeviceType: "cisco", Field: "DevOid", Description: "Device OID", OID: "1.3.6.1.2.1.1.2.0"},
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("localdb")
	collection := db.Collection("oids")

	// Insert OID data into the collection
	for _, oid := range oidData {
		_, err := collection.InsertOne(context.Background(), oid)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted OID: %s for DeviceType: %s\n", oid.OID, oid.DeviceType)
	}

	fmt.Println("OID data loaded successfully!")
}
