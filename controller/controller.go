package controller

import (
	"context"
	"discovery/logger"
	"discovery/models"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosnmp/gosnmp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func InitMongo() (err error) {
	options := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(context.Background(), options)
	if err != nil {
		logger.Log.Println("err in mongo connect: ", err)
		return
	}
	db = client.Database("Test")
	return err
}

func FindDeviceType(c *gin.Context) {
	var ip models.IP
	if err := c.ShouldBindJSON(&ip); err != nil {
		logger.Log.Println("bad request: ", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	deviceType := detectOS(ip.IP)
	logger.Log.Printf("Detected Device Type for IP %s: %s\n", ip, deviceType)
	c.JSON(http.StatusAccepted, deviceType)
}
func checkSNMP(ip string) (string, error) {
	g := &gosnmp.GoSNMP{
		Target:    ip,
		Port:      161,
		Version:   gosnmp.Version2c,
		Community: "public",
		Timeout:   time.Duration(2) * time.Second,
		Retries:   1,
	}

	err := g.Connect()
	if err != nil {
		return "", err
	}
	defer g.Conn.Close()

	oid := "1.3.6.1.2.1.1.1.0" // sysDescr OID
	result, err := g.Get([]string{oid})
	if err != nil {
		return "", err
	}

	for _, variable := range result.Variables {
		switch variable.Type {
		case gosnmp.OctetString:
			return string(variable.Value.([]byte)), nil
		}
	}
	return "", nil
}

func checkPort(ip string, port string) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", ip+":"+port, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func detectOS(ip string) string {
	// First, check if SNMP is available and get the system description
	snmpResult, err := checkSNMP(ip)
	if err == nil && snmpResult != "" {
		if strings.Contains(strings.ToLower(snmpResult), "windows") {
			return "Windows (SNMP)"
		} else if strings.Contains(strings.ToLower(snmpResult), "linux") {
			return "Linux (SNMP)"
		}
		return "Unknown (SNMP Device)"
	}

	// If SNMP doesn't work, fallback to port scanning
	if checkPort(ip, "135") || checkPort(ip, "445") {
		return "Windows (Port Scan)"
	} else if checkPort(ip, "22") {
		return "Linux (Port Scan)"
	}
	return "Unknown Device"
}
