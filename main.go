package main

import (
	"discovery/controller"
	"discovery/logger"

	"github.com/gin-gonic/gin"
)

// 192.168.1.4
// Enable-PSRemoting -Force

func main() {

	logger.Log.Println("server started")

	// get ttl value (if ttlvalue<=0 host not reachable)
	controller.InitMongo()
	router := gin.Default()

	// IP type
	router.POST("/checkIPType", controller.FindDeviceType)

	// SNMP
	router.POST("/snmp", controller.SnmpInput)
	router.POST("/snmpget", controller.StartDiscovery)

	router.Run(":8080")

}
