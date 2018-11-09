package configmanager

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"encoding/json"
)

func ReadFromYaml( t interface{}, filename string )  error {

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(yamlFile), t)
	if err != nil {
		log.Fatalf("error: %v", err)
		return err
	}

	return nil
}

func ReadFromJson( t interface{}, filename string )  error {

	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonFile), t)
	if err != nil {
		log.Fatalf("error: %v", err)
		return err
	}

	return nil
}
