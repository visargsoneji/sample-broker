package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetCatalog(c *gin.Context) {
	result := JSONReader("catalog.json")
	c.JSON(http.StatusOK, result)
}

func GetServiceInstance(c *gin.Context) {
	instanceId := c.Param("instance_id")

	if string(instanceId) != "instance1234" { // error
		c.JSON(http.StatusNotFound, gin.H{
			"error":             "404",
			"description":       "Instance Not Found",
			"instance_usable":   true,
			"update_repeatable": true,
		})
	} else {
		result := JSONReader("serviceInstance.json")
		c.JSON(http.StatusOK, result)
	}
}

func ProvisionServiceInstance(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	c.String(200, string(body))
}

func UpdateServiceInstance(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))

	c.JSON(http.StatusAccepted, gin.H{
		"dashboard_url": "dashboard_url",
		"operation":     "string",
	})
}

func DeprovisionServiceInstance(c *gin.Context) {
	serviceId := c.Query("service_id")
	planId := c.Query("plan_id")

	c.JSON(http.StatusAccepted, gin.H{
		"operation": "Deprovisioning",
		"serviceId": serviceId,
		"planId":    planId,
	})
}

func GetInstanceLastOperation(c *gin.Context) {

	if c.Param("instance_id") == "instance1234" {
		c.JSON(http.StatusOK, gin.H{
			"state":             "in progress",
			"description":       "lastOperation",
			"instance_usable":   true,
			"update_repeatable": true,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error":             "404",
			"description":       "instance not found",
			"instance_usable":   true,
			"update_repeatable": true,
		})
	}
}

func GetServiceBinding(c *gin.Context) {
	instanceId := c.Param("instance_id")
	bindingId := c.Param("binding_id")

	fmt.Println(instanceId, bindingId)

	result := JSONReader("serviceBinding.json")
	c.JSON(http.StatusOK, result)
}

func GenerateServiceBinding(c *gin.Context) {
	instanceId := c.Param("instance_id")
	bindingId := c.Param("binding_id")

	fmt.Println(instanceId, bindingId)

	result := JSONReader("serviceBinding.json")
	c.JSON(http.StatusCreated, result)
}

func DeleteServiceBinding(c *gin.Context) {
	serviceId := c.Query("service_id")
	planId := c.Query("plan_id")

	fmt.Println(serviceId, planId)
	c.JSON(http.StatusAccepted, gin.H{
		"operation": "DeleteServiceBinding successful",
	})
}

func GetBindingLastOperation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"state":             "in progress",
		"description":       "lastOperation",
		"instance_usable":   true,
		"update_repeatable": true,
	})
}

func JSONReader(fileName string) map[string]interface{} {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened %s\n", fileName)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

func main() {
	r := gin.Default()

	// Catalog
	r.GET("/v2/catalog", GetCatalog)

	// ServiceInstances
	r.GET("/v2/service_instances/:instance_id", GetServiceInstance)
	r.PUT("/v2/service_instances/:instance_id", ProvisionServiceInstance)
	r.PATCH("/v2/service_instances/:instance_id", UpdateServiceInstance)
	r.DELETE("/v2/service_instances/:instance_id", DeprovisionServiceInstance)
	r.GET("/v2/service_instances/:instance_id/last_operation", GetInstanceLastOperation)

	//ServiceBindings
	r.GET("/v2/service_instances/:instance_id/service_bindings/:binding_id", GetServiceBinding)
	r.PUT("/v2/service_instances/:instance_id/service_bindings/:binding_id", GenerateServiceBinding)
	r.DELETE("/v2/service_instances/:instance_id/service_bindings/:binding_id", DeleteServiceBinding)
	r.GET("/v2/service_instances/:instance_id/service_bindings/:binding_id/last_operation", GetBindingLastOperation)

	r.Run() // listen and serve on 0.0.0.0:8080
}
