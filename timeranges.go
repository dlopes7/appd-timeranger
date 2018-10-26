package main

import (
	"fmt"
	"time"

	"encoding/json"
	"github.com/dlopes7/go-appdynamics-rest-api/appdrest"
	"github.com/jinzhu/now"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Controller struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Account  string `json:"account"`
}

func readControllerConf() (*Controller, error) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	jsonPath := filepath.Join(filepath.Dir(ex), "controller.json")

	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var controller *Controller
	json.Unmarshal(byteValue, &controller)

	return controller, nil
}

var controller, _ = readControllerConf()
var client, _ = appdrest.NewClient(controller.Protocol, controller.Host, controller.Port, controller.User, controller.Password, controller.Account)

func updateTimeRange(name string, startTime time.Time, endTime time.Time) {
	timeRange, err := client.TimeRange.GetTimeRangeByName(name)

	if err != nil {
		fmt.Printf("Não foi possível obter o Time Range: %s\n", err.Error())
	}

	timeRange.ModifiedOn = time.Now().UnixNano() / 1000000
	timeRange.TimeRange.StartTime = startTime.UnixNano() / 1000000
	timeRange.TimeRange.EndTime = endTime.UnixNano() / 1000000

	_, err = client.TimeRange.UpdateTimeRange(*timeRange)
	if err != nil {
		fmt.Printf("Não foi possível alterar o Time Range %s: %s\n", name, err.Error())
	}

	fmt.Printf("Time Range '%s' atualizado\n", name)

}

func main() {

	updateTimeRange("Auto - Mes Atual", now.BeginningOfMonth(), time.Now())
	updateTimeRange("Auto - Semana Atual", now.BeginningOfWeek(), time.Now())
	updateTimeRange("Auto - Mes Passado", now.BeginningOfMonth().AddDate(0, -1, 0), now.BeginningOfMonth())
	updateTimeRange("Auto - Semana Passada", now.BeginningOfWeek().AddDate(0, 0, -7), now.BeginningOfWeek())
}
