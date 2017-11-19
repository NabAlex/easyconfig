package easyconfig

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const yamlExample = `
example:
  hello: "Hello, 世界"
  hello_int: 10
  hello_array:
    - "test1"
    - "test2"
    - "test3"
`

/* instead args */
func beforeInitZero() int {
	log.Printf("start before init test")
	reader := bytes.NewReader([]byte(yamlExample))
	yamlObj = parseFile(reader)
	if yamlObj == nil {
		log.Panicf("cannot parse test yaml \n%s", yamlExample)
	}

	configname = "testing"
	return 0
}

var (
	testZero = beforeInitZero()
	testVar  = GetInt("example.hello_int", testZero)
)

func TestBeforeInit(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(testVar, 10)
}

/* init() will be use default */
func TestParseInt(t *testing.T) {
	assert := assert.New(t)
	EnableWorkAfterInit()

	assert.Equal(GetInt("example.hello_int", 0), 10)
}

func TestParseInt64(t *testing.T) {
	assert := assert.New(t)
	EnableWorkAfterInit()

	assert.Equal(GetInt64("example.hello_int", 0), int64(10))
}

func TestParseString(t *testing.T) {
	assert := assert.New(t)
	EnableWorkAfterInit()

	assert.Equal(GetString("example.hello", ""), "Hello, 世界")
}

func TestParseStringArray(t *testing.T) {
	assert := assert.New(t)
	EnableWorkAfterInit()

	arr := GetArrayString("example.hello_array", nil)
	assert.Equal(arr, []string{"test1", "test2", "test3"})
}

func TestDefaultFromConfig(t *testing.T) {
	assert := assert.New(t)
	EnableWorkAfterInit()
	UseOnlyDefault(true)

	/* with out yaml file */
	assert.Equal(GetInt("path1.path2.maybe", -1), -1)
}
