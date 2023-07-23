package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/hashicorp/vault/api"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ClientConfig struct {
	Secret     string
	Address    string
	Cert       string
	Path       string
	ProjectEnv string
	Interval   int
	Mutex      *sync.Mutex
	Client     *api.Client
}

func (cc *ClientConfig) GetFullPath() string {
	return fmt.Sprintf("%s/%s", cc.ProjectEnv, cc.Path)
}

type ConfigFunc = func(cc *ClientConfig) (map[string]interface{}, error)

type Config interface {
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetArray(key string) []string
	GetMap(key string) map[string]string
}

type ConfigData struct {
	Data   map[string]interface{}
	Mutex  *sync.Mutex
	Secret string
	Path   string
	Client *api.Client
}

// get string config
func (cd *ConfigData) GetString(key string) string {
	item, ok := cd.Data[key]
	if ok {
		return fmt.Sprint(item)
	}

	return ""
}

// get boolean config
// boolean must be on string format "true" or "false"
func (cd *ConfigData) GetBool(key string) bool {
	item, ok := cd.Data[key]
	if !ok {
		return false
	}

	switch fmt.Sprint(item) {
	case "true":
		return true
	default:
		return false
	}
}

// get int config
func (cd *ConfigData) GetInt(key string) int {
	item, ok := cd.Data[key]
	if !ok {
		return 0
	}

	number, err := strconv.Atoi(fmt.Sprint(item))
	if err != nil {
		return 0
	}

	return number
}

// config array must sparated by comma
// example "omama,olala,omini"
func (cd *ConfigData) GetArray(key string) []string {
	item, ok := cd.Data[key]
	if !ok {
		return nil
	}

	return strings.Split(fmt.Sprint(item), ",")
}

// config array must be on json string format
// example {"name":"chapis","address":"bekasi"}
func (cd *ConfigData) GetMap(key string) map[string]string {
	item, ok := cd.Data[key]
	if !ok {
		return nil
	}

	var data = make(map[string]string)

	err := json.Unmarshal([]byte(fmt.Sprint(item)), &data)
	if err != nil {
		return nil
	}

	return data
}

// for reload purpose
func (cd *ConfigData) Reload(fn ConfigFunc, cc *ClientConfig) {
	s := gocron.NewScheduler(time.UTC)

	cc.Mutex.Lock()
	defer cc.Mutex.Unlock()

	go func() {
		if _, err := s.Every(cc.Interval).Seconds().Do(func() {
			// get config data
			data, err := fn(cc)
			if err != nil {
				log.Println("FAILED TO GET CONFIG: ", err)
				return
			}

			// set new config data
			cd.Data = data

		}); err != nil {
			fmt.Println("scheduler specify jobFunc: ", err)
			return
		}

		s.StartAsync()
	}()
}
