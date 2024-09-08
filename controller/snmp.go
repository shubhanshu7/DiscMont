package controller

import (
	"context"
	"discovery/logger"
	"discovery/models"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosnmp/gosnmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SnmpInput(c *gin.Context) {
	var snmpsetting models.SNMPInput
	logger.Log.Println("started snmp input")

	if err := c.ShouldBindJSON(&snmpsetting); err != nil {
		logger.Log.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err := db.Collection("SnmpSettings").InsertOne(context.Background(), snmpsetting)
	if err != nil {
		logger.Log.Println("Error: ", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.Log.Println("snmp input completed")
	c.JSON(http.StatusCreated, snmpsetting)
}

func StartDiscovery(c *gin.Context) {
	var snmpsetting models.SNMPInput
	logger.Log.Println("Started SNMP Discovery")

	if err := c.ShouldBindJSON(&snmpsetting); err != nil {
		logger.Log.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	oids, err := fetchOIDsFromDB(snmpsetting.DeviceProvider)
	if err != nil {
		logger.Log.Println("Error fetching OIDs: ", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	// logger.Log.Println("fetched OIDs: ", oids)
	logger.Log.Println("OIDs count: ", len(oids))

	var wg sync.WaitGroup
	resultCh := make(chan struct {
		Field string
		Value string
		Error error
	}, len(oids))

	// Perform SNMP GET requests in parallel using WaitGroup
	for _, oidMapping := range oids {
		wg.Add(1) // Increment WaitGroup counter

		go func(oidMapping models.OIDMapping) {
			defer wg.Done()

			value, err := performSNMPGet(snmpsetting.IP, snmpsetting.SnmpParameters, oidMapping.OID)
			resultCh <- struct {
				Field string
				Value string
				Error error
			}{Field: oidMapping.Field, Value: value, Error: err}
		}(oidMapping)
	}

	// Close the result channel once all goroutines are done
	go func() {
		wg.Wait()       // Wait for all goroutines to finish
		close(resultCh) // Close the channel to signal completion
	}()
	logger.Log.Println("processing completed: ")

	oidResults := make(map[string]string)
	for res := range resultCh {
		if res.Error != nil {
			logger.Log.Println("Error performing SNMP GET for field:", res.Field, res.Error.Error())
		} else {
			oidResults[res.Field] = res.Value
		}
	}

	deviceDetails := populateDeviceDetails(snmpsetting.DiscoveryID, oidResults)
	logger.Log.Println("Final output: ", deviceDetails)
	logger.Log.Println("Completed snmp discovery: ")

	c.JSON(http.StatusOK, deviceDetails)
}
func performSNMPGet(ip string, params models.SnmpParams, oid string) (string, error) {
	logger.Log.Println("snmp walk in progress ")
	var snmpargs gosnmp.GoSNMP
	var authProtocol gosnmp.SnmpV3AuthProtocol
	var privProtocol gosnmp.SnmpV3PrivProtocol
	var msgFlag gosnmp.SnmpV3MsgFlags
	var contextName string

	if params.Version == "1" || params.Version == "2c" || params.Version == "2C" || params.Version == "2" {
		snmpargs.Transport = "udp"
		snmpargs.Version = gosnmp.Version2c
		snmpargs.Port = 161
		snmpargs.Target = ip
		snmpargs.Timeout = time.Duration(10) * time.Second
		snmpargs.Retries = 3
		snmpargs.NonRepeaters = 0
		snmpargs.MaxRepetitions = 0
		snmpargs.Community = params.Community
	} else if params.Version == "3" {

		if len(params.ContextName) > 1 {
			contextName = params.ContextName
		}
		if params.AuthProtocol == "MD5" {
			authProtocol = gosnmp.MD5
		} else if params.AuthProtocol == "SHA" {
			authProtocol = gosnmp.SHA
		} else {
			logger.Log.Println("auth protocol is not supported ")
		}

		if params.PrivProtocol == "DES" {
			privProtocol = gosnmp.DES
		} else if params.PrivPassword == "AES" || params.PrivProtocol == "AES128" {
			privProtocol = gosnmp.AES
		} else if params.PrivProtocol == "AES192" {
			privProtocol = gosnmp.AES192
		} else if params.PrivProtocol == "AES256" {
			privProtocol = gosnmp.AES256
		} else if params.PrivProtocol == "AES192C" {
			privProtocol = gosnmp.AES192C
		} else if params.PrivProtocol == "AES256C" {
			privProtocol = gosnmp.AES256C
		} else {
			logger.Log.Println("priv protocol is not supported ")
		}

		if params.SecurityLevel == "NoAuthNoPriv" {
			msgFlag = gosnmp.NoAuthNoPriv
		} else if params.SecurityLevel == "AuthNoPriv" {
			msgFlag = gosnmp.AuthNoPriv
		} else if params.SecurityLevel == "AuthPriv" {
			msgFlag = gosnmp.AuthPriv
		} else {
			logger.Log.Println("Security level is not supported ")
		}

		snmpargs.Transport = "udp"
		snmpargs.Version = gosnmp.Version2c
		snmpargs.Port = 161
		snmpargs.Target = ip
		snmpargs.Timeout = time.Duration(10) * time.Second
		snmpargs.Retries = 3
		snmpargs.NonRepeaters = 0
		snmpargs.MaxRepetitions = 0
		snmpargs.SecurityModel = gosnmp.UserSecurityModel
		snmpargs.SecurityParameters = &gosnmp.UsmSecurityParameters{
			AuthoritativeEngineID:    params.SecurityEngineID,
			UserName:                 params.Username,
			AuthenticationProtocol:   authProtocol,
			AuthenticationPassphrase: params.AuthPassword,
			PrivacyProtocol:          privProtocol,
			PrivacyPassphrase:        params.PrivPassword,
		}
		snmpargs.MsgFlags = msgFlag
		snmpargs.ContextName = contextName
	}
	// snmp := &gosnmp.GoSNMP{
	// 	Target:    ip,
	// 	Port:      161,
	// 	Version:   gosnmp.Version2c,
	// 	Community: params.Community,
	// 	Timeout:   params.Timeout,
	// 	Retries:   int(params.Retries),
	// }
	snmp := &snmpargs
	err := snmp.Connect()
	if err != nil {
		return "", err
	}
	defer snmp.Conn.Close()

	pdu, err := snmp.Get([]string{oid})
	if err != nil {
		return "", err
	}

	// Extract the value from the SNMP response
	if len(pdu.Variables) > 0 {
		return pdu.Variables[0].Value.(string), nil
	}
	return "", fmt.Errorf("no value found for OID %s", oid)
}
func fetchOIDsFromDB(deviceType string) ([]models.OIDMapping, error) {
	var oids []models.OIDMapping
	collection := db.Collection("oids")

	filter := primitive.M{"deviceType": deviceType}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var oidMapping models.OIDMapping
		if err := cursor.Decode(&oidMapping); err != nil {
			return nil, err
		}
		oids = append(oids, oidMapping)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return oids, nil
}
func populateDeviceDetails(discoveryID string, results map[string]string) models.Devicedetails {
	return models.Devicedetails{
		DiscoverID:         discoveryID,
		DeviceName:         results["DeviceName"],
		DeviceStatus:       results["DeviceStatus"],
		Device_DNS_Name:    results["Device_DNS_Name"],
		DeviceDesc:         results["DeviceDesc"],
		DeviceType:         results["DeviceType"],
		Device_MacAddress:  results["Device_MacAddress"],
		DevicePrimaryIP:    results["DevicePrimaryIP"],
		Device_long_desc:   results["Device_long_desc"],
		Device_Manf_Name:   results["Device_Manf_Name"],
		DeviceModelName:    results["DeviceModelName"],
		DeviceSerialNumber: results["DeviceSerialNumber"],
		Device_soft_rev:    results["Device_soft_rev"],
		Device_firm_rev:    results["Device_firm_rev"],
		DeviceOS:           results["DeviceOS"],
		DeviceOS_version:   results["DeviceOS_version"],
		Device_Total_Port:  results["Device_Total_Port"],
		Device_Free_port:   results["Device_Free_port"],
		Device_no_of_cpus:  results["Device_no_of_cpus"],
		Numberofhosts:      results["Numberofhosts"],
		Memory:             results["Memory"],
		DiscStorage:        results["DiscStorage"],
		PollingMethod:      results["PollingMethod"],
		DevOid:             results["DevOid"],
		Createdat:          time.Now(),
		Updatedat:          time.Now(),
	}
}
