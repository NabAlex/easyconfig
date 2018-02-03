package easyconfig

import (
	"errors"
	"flag"
	"github.com/beego/goyaml2"
	"io"
	"log"
	"os"
	"strings"
	"runtime"
)

const module = "easyconfig"

var (
	checkInit  = true
	useDefault = false
	configname = ""
)

var yamlObj interface{} = nil

func parseFile(r io.Reader) interface{} {
	_yaml, err := goyaml2.Read(r)
	if err != nil {
		log.Panicf(module + ": cannot parse file [%s]", configname)
	}

	return _yaml
}

func initYaml() (interface{}, bool) {
	filename := flag.String("c", "", "use [-c filename]")
	flag.Parse()
	if *filename == "" {
		log.Printf(module + ": use default values (no flag \"conf\")")
		return nil, false
	}

	f, err := os.Open(*filename)
	if err != nil {
		log.Panicf(module + ": file does not exist [%s]", configname)
	}

	result := parseFile(f)
	if result == nil {
		log.Panicf(module + ": cannot parse config")
	}

	configname = *filename
	return result, true
}

func toMap(obj interface{}) (map[string]interface{}, error) {
	if obj == nil {
		return nil, errors.New(module + ": key is nil")
	}

	m, ok := obj.(map[string]interface{})
	if !ok {
		return nil, errors.New(module + ": cannot cast to map")
	}

	return m, nil
}

func toList(obj interface{}) ([]interface{}, error) {
	if obj == nil {
		return nil, errors.New(module + ": key is nil")
	}

	l, ok := obj.([]interface{})
	if !ok {
		return nil, errors.New(module + ": cannot cast to list")
	}

	return l, nil
}

func _getCallerIsFromInit() bool {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(4, fpcs)
	if n == 0 {
		return false
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return false
	}

	modulePath := strings.Split(fun.Name(), ".")
	return modulePath[len(modulePath) - 1] == "init"
}

func getVar(pathToValue string) (interface{}, bool) {
	if checkInit && !_getCallerIsFromInit() {
		log.Panicf(module + ": call function after init, use EnableWorkWithConfig()")
	}

	if useDefault {
		return nil, false
	}

	if yamlObj == nil {
		var ok bool
		yamlObj, ok = initYaml()
		if !ok {
			useDefault = true
			return nil, false
		}
	}

	tmp := yamlObj
	paths := strings.Split(pathToValue, ".")
	for _, path := range paths {
		newMap, err := toMap(tmp)
		if err != nil {
			log.Println(err)
			log.Panicf(module + ": cannot key: [%s] path: [%s] file: [%s]", path, pathToValue, configname)
		}

		var ok bool
		tmp, ok = newMap[path]
		if !ok {
			log.Panicf(module + ": cannot find key: [%s] path: [%s] file: [%s]", path, pathToValue, configname)
		}
	}

	return tmp, true
}

/* debug */
func EnableWorkAfterInit() {
	checkInit = true
}

/* debug */
func UseOnlyDefault(flag bool) {
	useDefault = flag
}

func GetInt64(pathToValue string, defaultValue int64) int64 {
	el, ok := getVar(pathToValue)
	if !ok {
		return defaultValue
	}

	result, ok := el.(int64)
	if !ok {
		log.Panicf(module + ": value is not \"int64\" path: [%s] file: [%s]", pathToValue, configname)
	}

	return result
}

func GetInt(pathToValue string, defaultValue int) int {
	el, ok := getVar(pathToValue)
	if !ok {
		return defaultValue
	}

	result, ok := el.(int64)
	if !ok {
		log.Panicf(module + ": value is not \"int64\" path: [%s] file: [%s]", pathToValue, configname)
	}

	return int(result)
}

func GetString(pathToValue string, defaultValue string) string {
	el, ok := getVar(pathToValue)
	if !ok {
		return defaultValue
	}

	result, ok := el.(string)
	if !ok {
		log.Panicf(module + ": value is not \"string\" path: [%s] file: [%s]", pathToValue, configname)
	}

	return result
}

func GetArrayString(pathToValue string, defaultValue []string) []string {
	el, ok := getVar(pathToValue)
	if !ok {
		return defaultValue
	}

	arr, err := toList(el)
	if err != nil {
		log.Println(err)
		log.Panicf(module + ": value is not \"[]interface{}\" path: [%s] file: [%s]", pathToValue, configname)
	}

	arrString := make([]string, 0, len(arr))
	for _, val := range arr {
		str, ok := val.(string)
		if !ok {
			log.Panicf(module + ": value is not \"[]string\" path: [%s] file: [%s]", pathToValue, configname)
		}

		arrString = append(arrString, str)
	}

	return arrString
}
