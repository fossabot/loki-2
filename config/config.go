package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cjburchell/uatu-go"

	"github.com/pkg/errors"
)

// Endpoint configuration
type Endpoint struct {
	Path         string            `json:"path"`
	Method       string            `json:"method"`
	ResponseBody json.RawMessage   `json:"response_body"`
	StringBody   string            `json:"string_body"`
	ContentType  string            `json:"content_type"`
	Response     int               `json:"response"`
	Header       map[string]string `json:"header"`
	Name         string            `json:"-"`
	ReplyDelay   int               `json:"reply_delay"`
}

// GetEndpoints configuration
func GetEndpoints(log log.ILog) ([]Endpoint, error) {
	results, err := load(log)
	if err != nil {
		return nil, err
	}

	endpoints := make([]Endpoint, len(results))
	index := 0
	for name, value := range results {
		value.Name = name
		endpoints[index] = value
		index++
	}

	return endpoints, nil
}

// GetEndpoint with given ID
func GetEndpoint(id string, log log.ILog) (*Endpoint, error) {
	results, err := load(log)
	if err != nil {
		return nil, err
	}

	if item, ok := results[id]; ok {
		return &item, nil
	}

	return nil, nil
}

// DeleteEndpoint in configuration
func DeleteEndpoint(id string, log log.ILog) error {
	endpoints, err := load(log)
	if err != nil {
		return err
	}

	if _, ok := endpoints[id]; ok {
		delete(endpoints, id)
		return save(endpoints)
	}

	return errors.WithStack(fmt.Errorf("unable to find logger with ID %s", id))
}

// Setup the configuration
func Setup(file string) error {
	configFileName = file
	return nil
}

var configFileName string

func load(log log.ILog) (map[string]Endpoint, error) {
	loggers := make(map[string]Endpoint)
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		log.Warnf("Config file %s not found", configFileName)
		return loggers, nil
	}

	log.Printf("loading config file %s", configFileName)
	fileData, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return loggers, errors.WithStack(err)
	}

	err = json.Unmarshal(fileData, &loggers)
	return loggers, errors.WithStack(err)
}

func save(config map[string]Endpoint) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return errors.WithStack(err)
	}

	return ioutil.WriteFile(configFileName, configJSON, 0644)
}
