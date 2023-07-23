package vault

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/vault/api"
	configdata "github.com/ijlik/dating-user/pkg/config/data"
)

func NewVaultEnv(path string, interval int) configdata.Config {
	address := os.Getenv("VAULT_ADDR")
	secret := os.Getenv("VAULT_TOKEN")
	ca := os.Getenv("VAULT_CA_CERT")
	projectEnv := os.Getenv("PROJECT_ENV")

	// create func to get vault data
	getData := getVaultEnv()

	clientConfig := &configdata.ClientConfig{
		Address:    address,
		Cert:       ca,
		Secret:     secret,
		Path:       path,
		ProjectEnv: projectEnv,
		Interval:   interval,
		Mutex:      &sync.Mutex{},
	}

	// get vault data
	data, err := getData(clientConfig)
	if err != nil {
		return nil
	}

	// set config
	config := configdata.ConfigData{
		Data:   data,
		Secret: secret,
		Path:   path,
	}

	// enable reload config
	config.Reload(getData, clientConfig)

	return &config
}

func getVaultEnv() func(cc *configdata.ClientConfig) (map[string]interface{}, error) {
	return func(cc *configdata.ClientConfig) (map[string]interface{}, error) {
		var err error
		client := cc.Client
		if client == nil {
			client, err = getClient(cc.Secret, cc.Cert, cc.Address)
			if err != nil {
				return nil, err
			}
		}

		resp, err := client.Logical().Read(cc.GetFullPath())
		if err != nil {
			return nil, err
		}

		return convertToMap(resp.Data)
	}
}

func convertToMap(data map[string]interface{}) (map[string]interface{}, error) {
	param, ok := data["data"]
	if !ok {
		return nil, errors.New("not found")
	}

	d, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	var md map[string]interface{}
	err = json.Unmarshal(d, &md)
	if err != nil {
		return nil, err
	}

	return md, nil

}

func getClient(secret, ca, address string) (*api.Client, error) {
	config := api.Config{
		Address: address,
	}

	if len(ca) > 0 {
		if err := config.ConfigureTLS(&api.TLSConfig{
			CACert: ca,
		}); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	client, err := api.NewClient(&config)
	if err != nil {
		return nil, err
	}
	client.SetToken(secret)

	return client, nil
}
