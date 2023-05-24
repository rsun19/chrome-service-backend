package main

import (
	"fmt"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
	"github.com/xeipuuv/gojsonschema"
)

var duplicateCounter = make(map[string]int)

func validateNavigation(cwd string) {
	schemaLoader := gojsonschema.NewReferenceLoader("file://./static/navigationSchema.json")
	modulesFiles, err := filepath.Glob(cwd + "/static/**/**/**/*-navigation.json")
	handleErr(err)

	for _, file := range modulesFiles {
		documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", file))
		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		
		if err != nil {
			fmt.Println("File", file)
			panic(err.Error())
		}

		if !result.Valid() {
			for _, desc := range result.Errors() {
				fmt.Printf("- %s\n", desc)
			}
			panic(fmt.Sprintf("The %s is not valid. see errors :\n", file))
		}

		var data map[string]interface {}
		
		var arrayData [] map[string]interface {};

		fileContent, status := ioutil.ReadFile(file)
		handleErr(status)

		ok:= json.Unmarshal(fileContent, &data)
		if ok == nil{
			//Note, only works if there are no errors unmarshalling!
			duplicateCounter = make(map[string]int)
			parseJSONIDs(data, file)
		} else {
			ok := json.Unmarshal(fileContent, &arrayData)
			if ok == nil {
				duplicateCounter = make(map[string]int)
				jsonArrayData := make([]interface{}, len(arrayData))
				for k,v := range arrayData {
					jsonArrayData[k] = v
				}
				loopOverFields(jsonArrayData, file)
				//parsingSetup(rootJSONFile, file)
			} else {
				panic(ok.Error())
			}
		}
	}
}

func parseJSONIDs(data map[string]interface{}, file string) {
	if idValue, ok := data["id"]; ok {
		if idMap, ok := idValue.(string); ok {
			// duplicateCounter[idMap] = 1
			if _, exists := duplicateCounter[idMap]; exists {
				//panic(fmt.Sprintf("The id %s in %s is not valid because it is duplicated\n", id, file))
			} else {
				duplicateCounter[idMap] = 1
			}
			//FOR DEBUGGING PURPOSES
			fmt.Sprintf(idMap)
		}
	}
	if navItems, ok := data["navItems"].([]interface{}); ok {
		loopOverFields(navItems, file)
	} else {
		if routeItems, ok := data["routes"].([]interface{}); ok {
			loopOverFields(routeItems, file)
		}
	}
}

func loopOverFields(navItems []interface{}, file string) {
	for i := 0; i < len(navItems); i++ {
		if navItem, ok := navItems[i].(map[string]interface{}); ok {
				//if value, ok := navItem["id"]; ok {
					// if id, ok := value.(string); ok {
					// 	if _, exists := duplicateCounter[id]; exists {
					// 		//panic(fmt.Sprintf("The id %s in %s is not valid because it is duplicated\n", id, file))
					// 	} else {
					// 		duplicateCounter[id] = 1
					// 	}
					// 	//FOR DEBUGGING PURPOSES
					// 	fmt.Println("----")
					// 	fmt.Println(id)
					// 	fmt.Println(file)
					// }
				//}
				parseJSONIDs(navItem, file)
			} else {
				panic(fmt.Sprintf("Invalid format. The 'navItems' field MUST be a map"))
			}
	}
}