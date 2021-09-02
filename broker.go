package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Instance struct {
	instance_id string
	service_id  string
	plan_id     string
	bindings    []string
}

var InstanceList = []Instance{
	{instance_id: "instance99", service_id: "service99", plan_id: "plan99", bindings: []string{"binding99", "binding999"}},
	{instance_id: "instance88", service_id: "service88", plan_id: "plan88", bindings: []string{"binding88", "binding888"}},
}

func GetCatalog(c *gin.Context) {
	result := JSONReader("catalog.json")
	c.JSON(http.StatusOK, result)
}

func GetServiceInstance(c *gin.Context) {
	instanceId := c.Param("instance_id")
	dashboardUrl := "http://example-dashboard.example.com/"
	statusCode := 404

	for _, instance := range InstanceList {
		if instance.instance_id == instanceId {
			fmt.Printf("Matched Instance: %+v\n", instance)
			dashboardUrl += instanceId
			statusCode = 200
			c.JSON(statusCode, gin.H{
				"dashboard_url": dashboardUrl,
				"service_id":    instance.service_id,
				"plan_id":       instance.plan_id,
				"binding_ids":   fmt.Sprint(instance.bindings),
			})
		}
	}

	fmt.Println("-Get-")
	for i, instance := range InstanceList {
		fmt.Printf("%+v at %d\n", instance, i)
	}
	fmt.Println("--")

	if statusCode == 200 {
		return
	} else {
		c.JSON(statusCode, gin.H{
			"error":       "404",
			"description": "Instance Not Found",
		})
	}
}

func ProvisionServiceInstance(c *gin.Context) {
	byteValue, _ := ioutil.ReadAll(c.Request.Body)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	instanceId := c.Param("instance_id")

	serviceId, ok := result["service_id"]
	if !ok {
		fmt.Println("service_id not found in request body")
		return
	} else {
		fmt.Printf("ServiceId: %v\n", serviceId)
	}
	planId, ok := result["plan_id"]
	if !ok {
		fmt.Println("plan_id not found in request body")
		return
	} else {
		fmt.Printf("PlanId: %v\n", planId)
	}

	statusCode := 201
	for _, instance := range InstanceList {
		//fmt.Printf("%+v at %d\n", instance, i)
		if instance.instance_id == instanceId {
			statusCode = 200
			break
		}
	}
	service_id := fmt.Sprint(serviceId)
	plan_id := fmt.Sprint(planId)

	if statusCode == 201 {
		InstanceList = append(InstanceList, Instance{
			instance_id: instanceId,
			service_id:  service_id,
			plan_id:     plan_id,
			bindings:    []string{},
		})
	}

	fmt.Println("--")
	for i, instance := range InstanceList {
		fmt.Printf("%+v at %d\n", instance, i)
	}
	fmt.Println("--")

	dashboard_url := "http://example-dashboard.example.com/" + instanceId + "/" + service_id + "/" + plan_id
	c.JSON(statusCode, gin.H{
		"dashboard_url": dashboard_url,
		"operation":     "task_10",
		"metadata":      "",
	})
}

func UpdateServiceInstance(c *gin.Context) {
	byteValue, _ := ioutil.ReadAll(c.Request.Body)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	instanceId := c.Param("instance_id")

	serviceId, ok := result["service_id"]
	if !ok {
		fmt.Println("service_id not found in request body")
		return
	} else {
		fmt.Printf("ServiceId: %v\n", serviceId)
	}
	planId, ok := result["plan_id"]
	if !ok {
		fmt.Println("plan_id not found in request body")
		return
	} else {
		fmt.Printf("PlanId: %v\n", planId)
	}

	service_id := fmt.Sprint(serviceId)
	plan_id := fmt.Sprint(planId)

	statusCode := 400

	for i, instance := range InstanceList {
		if instance.instance_id == instanceId && instance.service_id == service_id {
			fmt.Println("Instance and Service ID matched..")
			if instance.plan_id != plan_id {
				fmt.Println("Updating with new PlanId: ", plan_id)
				InstanceList[i].plan_id = plan_id
				fmt.Println("New PlanId:", InstanceList[i].plan_id)
			}
			statusCode = 200

			break
		}
		//fmt.Printf("%+v\n", instance)
	}
	if statusCode == 200 {
		dashboard_url := "http://example-dashboard.example.com/" + instanceId + "/" + service_id + "/" + plan_id
		c.JSON(statusCode, gin.H{
			"dashboard_url": dashboard_url,
			"operation":     "update-plan",
		})
	} else {
		c.JSON(statusCode, gin.H{
			"error":       400,
			"description": "Bad Request",
		})
	}

	fmt.Println("-PATCH-")
	for i, instance := range InstanceList {
		fmt.Printf("%+v at %d\n", instance, i)
	}
	fmt.Println("--")
}

func removeInst(instList []Instance, i int) []Instance {
	instList[i] = instList[len(instList)-1]
	return instList[:len(instList)-1]
}

func DeprovisionServiceInstance(c *gin.Context) {
	serviceId := c.Query("service_id")
	planId := c.Query("plan_id")
	instanceId := c.Param("instance_id")
	statusCode := 410 // Gone

	for i, instance := range InstanceList {
		if instance.instance_id == instanceId && instance.service_id == serviceId && instance.plan_id == planId {
			statusCode = 200
			InstanceList = removeInst(InstanceList, i)
			break
		} else if instance.instance_id == instanceId && (instance.service_id != serviceId || instance.plan_id != planId) {
			statusCode = 400
			break
		}
	}

	if statusCode == 200 {
		c.JSON(statusCode, gin.H{
			"operation": "task10-deprovisioning",
		})
	} else {
		var description string
		if statusCode == 400 {
			description = "Bad Request. ServiceId or PlanId doesn't match"
		} else {
			description = "Instance not found"
		}
		c.JSON(statusCode, gin.H{
			"error":       statusCode,
			"description": description,
		})
	}
}

func GetInstanceLastOperation(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"state":             "in progress",
		"description":       "lastOperation",
		"instance_usable":   true,
		"update_repeatable": true,
	})
}

func GetServiceBinding(c *gin.Context) {
	instanceId := c.Param("instance_id")
	bindingId := c.Param("binding_id")
	statusCode := 404

	for _, instance := range InstanceList {
		if instance.instance_id == instanceId {
			for _, binding := range instance.bindings {
				if binding == bindingId {
					fmt.Println("Binding Found!")
					statusCode = 200
				}
			}
		}
	}

	if statusCode == 200 {
		result := JSONReader("serviceBinding.json")
		c.JSON(statusCode, result)
	} else {
		c.JSON(statusCode, gin.H{
			"error":       404,
			"description": "Service binding not found. Check instance_id and binding_id",
		})
	}
}

func GenerateServiceBinding(c *gin.Context) {
	instanceId := c.Param("instance_id")
	bindingId := c.Param("binding_id")

	byteValue, _ := ioutil.ReadAll(c.Request.Body)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	serviceId, ok := result["service_id"]
	if !ok {
		fmt.Println("service_id not found in request body")
		return
	} else {
		fmt.Printf("ServiceId: %v\n", serviceId)
	}
	planId, ok := result["plan_id"]
	if !ok {
		fmt.Println("plan_id not found in request body")
		return
	} else {
		fmt.Printf("PlanId: %v\n", planId)
	}

	service_id := fmt.Sprint(serviceId)
	plan_id := fmt.Sprint(planId)
	statusCode := 400

	fmt.Println(instanceId, bindingId)

	for i, instance := range InstanceList {
		if instance.instance_id == instanceId && instance.plan_id == plan_id && instance.service_id == service_id {
			for _, binding := range instance.bindings {
				if binding == bindingId {
					fmt.Println("Binding already exist!")
					statusCode = 200
					break
				}
			}
			if statusCode != 200 {
				InstanceList[i].bindings = append(InstanceList[i].bindings, bindingId)
				statusCode = 201
			}
		}
	}

	if statusCode == 400 {
		c.JSON(statusCode, gin.H{
			"error":       400,
			"description": "Bad Request. Parameters doesn't match",
		})
	} else {
		result := JSONReader("serviceBinding.json")
		c.JSON(statusCode, result)
	}

}

func removeBinding(bindings []string, i int) []string {
	bindings[i] = bindings[len(bindings)-1]
	return bindings[:len(bindings)-1]
}

func DeleteServiceBinding(c *gin.Context) {
	instanceId := c.Param("instance_id")
	bindingId := c.Param("binding_id")
	serviceId := c.Query("service_id")
	planId := c.Query("plan_id")

	fmt.Printf("ServiceId: %s -- PlanID: %s", serviceId, planId)
	if serviceId == "" || planId == "" {
		c.JSON(404, gin.H{
			"error":       404,
			"description": "Bad request. Missing parameter/s",
		})
		return
	}

	statusCode := 410 // gone

	for i, instance := range InstanceList {
		if instance.instance_id == instanceId && instance.service_id == serviceId && instance.plan_id == planId {
			for j, binding := range instance.bindings {
				if binding == bindingId {
					fmt.Println("Binding found!")
					fmt.Println("Deleting the binding!")
					InstanceList[i].bindings = removeBinding(InstanceList[i].bindings, j)
					statusCode = 200
				}
			}
		}
	}

	if statusCode == 200 {
		c.JSON(statusCode, gin.H{
			"operation": "DeleteServiceBinding successful",
		})
	} else {
		c.JSON(statusCode, gin.H{
			"error":       410,
			"description": "Binding gone",
		})
	}

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
