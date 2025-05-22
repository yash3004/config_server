package cmd

import (
	"flag"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

type BindOptions struct {
	HTTP int `yaml:"http"`
	GRPC int `yaml:"grpc"`
}

type Configurations struct {
	MongoURI string      `yaml:"mongoURI"`
	Bind     BindOptions `yaml:"bind"`
	UseFile  bool        `yaml:"use_file"`
}

var (
	configOnce sync.Once
	config     Configurations
	klogOnce   sync.Once
)

func GetConfigurations() Configurations {
	configOnce.Do(func() {
		klogOnce.Do(func() {
			klog.InitFlags(nil)
			klog.EnableContextualLogging(true)
		})

		var configPath string

		flagSet := flag.NewFlagSet("config", flag.ContinueOnError)
		flagSet.StringVar(&configPath, "cfg", "config.yaml", "Configuration File")

		if flag.Parsed() {
			if cfgFlag := flag.Lookup("cfg"); cfgFlag != nil {
				configPath = cfgFlag.Value.String()
			}
		} else {
			flag.StringVar(&configPath, "cfg", "config.yaml", "Configuration File")
			flag.Parse()
		}

		file, err := os.Open(configPath)
		if err != nil {
			klog.Fatalf("cannot read config file:%v", err)
		}
		defer file.Close()

		if err := yaml.NewDecoder(file).Decode(&config); err != nil {
			klog.Fatalf("cannot unmarshal the yaml file %v", err)
		}
	})

	return config
}