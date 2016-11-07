package parseConfig

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func ParseYML(yamlfile string) map[string]string {
	yamlFile, _ := ioutil.ReadFile(yamlfile)
	any := map[string]interface{}{}
	err := yaml.Unmarshal(yamlFile, &any)
	if err != nil {
		log.Fatal(err)
	}

	flatmap := map[string]string{}
	for k, v := range any {
		flatten(k, v, flatmap)
	}
	return flatmap
}

func flatten(prefix string, value interface{}, flatmap map[string]string) {
	submap, ok := value.(map[interface{}]interface{})
	if ok {
		for k, v := range submap {
			flatten(prefix+"."+k.(string), v, flatmap)
		}
		return
	}
	stringlist, ok := value.([]interface{})
	if ok {
		flatten(fmt.Sprintf("%s.size", prefix), len(stringlist), flatmap)
		for i, v := range stringlist {
			flatten(fmt.Sprintf("%s.%d", prefix, i), v, flatmap)
		}
		return
	}
	flatmap[prefix] = fmt.Sprintf("%v", value)
}
