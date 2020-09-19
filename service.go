package main

import (
	"fmt"
	"github.com/s8sg/goflow-dashboard/lib"
	redis "gopkg.in/redis.v5"
	"os"
	"strings"
)

var rdb *redis.Client

// listGoFLows get list of go-flows
func listGoFLows() ([]*Flow, error) {
	rdb = getRDB()
	command := rdb.Keys("goflow-flow:*")
	rdb.Process(command)
	flowKeys, err := command.Result()
	if err != nil {
		return nil, nil
	}
	flows := make([]*Flow, 0)
	for _, key := range flowKeys {
		flowName := strings.Split(key, ":")[1]
		if flowName == "" {
			continue
		}
		flow := &Flow{
			Name: flowName,
		}
		flows = append(flows, flow)
	}
	return flows, nil
}

// getDot request to dot-generator for the dag dot graph
func getDot(flowName string) (string, error) {
	rdb = getRDB()
	command := rdb.Get("goflow-flow:"+flowName)
	rdb.Process(command)
	definition, err := command.Result()
	if err != nil {
		return "", nil
	}
	dot, err := lib.MakeDotFromDefinitionString(definition)
	return dot, err
}

// listFlowRequests get list of request for a goflow
func listFlowRequests(flow string) (map[string]string, error) {
	return lib.ListRequests(flow)
}

// buildFlowDesc get a flow details
func buildFlowDesc(functions []*Flow, flowName string) (*FlowDesc, error) {

	var functionObj *Flow
	for _, functionObj = range functions {
		if functionObj.Name == flowName {
			break
		}
	}

	dot, dErr := getDot(flowName)
	if dErr != nil {
		return nil, fmt.Errorf("failed to get dot, %v", dErr)
	}

	flowDesc := &FlowDesc{
		Name:            functionObj.Name,
		Dot:             dot,
	}

	return flowDesc, nil
}

// listRequestTraces get list of traces for a request traceID
func listRequestTraces(requestId string,  requestTraceId string) (*RequestTrace, error) {
	requestTraceResponse, err := lib.ListTraces(requestTraceId)
	if err == nil {
		return nil, err
	}
	requestTrace := &RequestTrace{
		RequestID: requestId,
		TraceId: requestTraceId,
		StartTime: requestTraceResponse.StartTime,
		NodeTraces: make(map[string]*NodeTrace, 0),
		Duration: requestTraceResponse.Duration,
	}
	for id, nodeTrace := range requestTraceResponse.NodeTraces {
		nodeTraceObj := &NodeTrace{
			StartTime: nodeTrace.StartTime,
			Duration: nodeTrace.Duration,
		}
		requestTrace.NodeTraces[id] = nodeTraceObj
	}

	return requestTrace, nil
}

// getRequestStatus request the flow for the request status
func getRequestStatus(flow, requestTraceId string) (string, error) {
	rdb = getRDB()
	return "", nil
}

func getRDB() *redis.Client{
	addr := os.Getenv("redis_url")
	if addr == "" {
		addr = "localhost:6379"
	}
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   0,
		})
	}
	return rdb
}