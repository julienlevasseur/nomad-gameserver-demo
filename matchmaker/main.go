package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	nomad "github.com/hashicorp/nomad/api"
)

// gameServer map the game server name with its job description file
var gameServer = map[string]string{
	"minecraft": "../minecraft-docker.nomad",
}

// RootCmd root command
var RootCmd = &cobra.Command{
	Use:   "matchmaker",
	Short: "A simple POC matcmaker tool.",
	Run: func(cmd *cobra.Command, args []string) {
		requestGameServer(args[0])
	},
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

type JobDispatchResponse struct {
	JobDispatchResponse jobDispatchResponse `json:"JobDispatchResponse"`
}

type jobDispatchResponse struct {
	LastIndex       uint64 `json:"LastIndex"`
	JobCreateIndex  uint64 `json:"JobCreateIndex"`
	EvalCreateIndex uint64 `json:"EvalCreateIndex"`
	EvalID          string `json:"EvalID"`
	DispatchedJobID string `json:"DispatchedJobID"`
	RequestTime     uint64 `json:"RequestTime"`
}

type consulServices struct {
	Name string
	Tags []string
}

func nomadClient() (*nomad.Jobs, error) {
	n, err := nomad.NewClient(nomad.DefaultConfig())
	if err != nil {
		return &nomad.Jobs{}, err
	}
	return n.Jobs(), nil
}

func registerJob(jobName string) (string, error) {
	url := "http://127.0.0.1:7070/v1/jobs/"

	data, err := ioutil.ReadFile(gameServer[jobName])
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return resp.Status, nil
}

func dispatchJob(jobName string) (string, JobDispatchResponse, error) {
	url := "http://127.0.0.1:7070/v1/jobs/" + jobName + "/dispatch"

	data := []byte(`{}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", JobDispatchResponse{}, err
	}

	client := &http.Client{Timeout: time.Second * 10}

	var dispatchResponse JobDispatchResponse

	resp, err := client.Do(req)
	if err != nil {
		return "", JobDispatchResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", JobDispatchResponse{}, err
	}

	err = json.Unmarshal(body, &dispatchResponse)
	if err != nil {
		return "", JobDispatchResponse{}, err
	}
	return resp.Status, dispatchResponse, nil
}

func consulServiceByTag(tag string) (string, error) {
	url := "http://127.0.0.1:8500/v1/catalog/services"

	data := []byte(`{}`)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var services *consulServices

	bytes, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(bytes))
	err = json.Unmarshal(bytes, &services)
	fmt.Println(services)
	return "", nil
}

//func FindTraefikServerAddress {
//
//}

// createGameServer func ask Nomad for the creation of a game server
func requestGameServer(jobName string) {
	n, err := nomadClient()
	if err != nil {
		panic(err)
	}

	// Check for minecraft-server-register job:
	jobs, _, err := n.List(nil)
	if err != nil {
		panic(err)
	}

	// Retrieving all Noamd jobs name:
	var jobNames []string
	for _, j := range jobs {
		jobNames = append(jobNames, j.Name)
	}

	// If the game server job is not registered, register it, otherwise,
	// dispatch it:
	if contains(jobNames, jobName) {
		fmt.Printf(
			"%v gameserver description found, dispatch the game server ...\n",
			jobName,
		)
		_, jdr, err := dispatchJob(jobName)
		if err != nil {
			panic(err)
		} else {
			time.Sleep(3 * time.Second)
			service, err := consulServiceByTag(jdr.JobDispatchResponse.DispatchedJobID)
			if err != nil {
				panic(err)
			}
			fmt.Println(service)
		}
	} else {
		fmt.Printf(
			"%v gameserver description not found, register the game server ...\n",
			jobName,
		)

		status, err := registerJob(jobName)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(jobName, " Register - Response Status: ", status)
		}

		_, jdr, err := dispatchJob(jobName)
		if err != nil {
			panic(err)
		} else {
			time.Sleep(3 * time.Second)
			service, err := consulServiceByTag(jdr.JobDispatchResponse.DispatchedJobID)
			if err != nil {
				panic(err)
			}
			fmt.Println(service)
		}
	}
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}
