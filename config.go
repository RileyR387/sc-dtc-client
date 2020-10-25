package main

import (
    "os"
    "strings"
    "fmt"
    "path/filepath"
    log "github.com/sirupsen/logrus"
    "github.com/pborman/getopt/v2"
    "github.com/spf13/viper"
)

var (
    me = filepath.Base(os.Args[0])
    yamlFile = fmt.Sprintf("%s.yaml", me)
    envPrefix = "DTCCLIENT"
    configSearchPaths = []string {".", "./etc", "$HOME/.dtc-client-go/", "$HOME/etc", "/etc"}
    genConfig = getopt.BoolLong("genconfig", 'x', "Write example config to \"./" + yamlFile + "\"")
)

func init() {
    viper.SetConfigName(yamlFile)
    viper.SetConfigType("yaml")
    viper.SetEnvPrefix(envPrefix)
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()

    for _, p := range configSearchPaths {
        viper.AddConfigPath(p)
    }

    viper.SetDefault("dtc.Host", "127.0.0.1")
    viper.SetDefault("dtc.Port", "11099")
    viper.SetDefault("dtc.HistPort", "11098")
    viper.SetDefault("dtc.Username", "")
    viper.SetDefault("dtc.Password", "")
    viper.SetDefault("dtc.Reconnect", false)
    viper.SetDefault("log.level", "TRACE")

    getopt.SetUsage(func() { usage() })
    getopt.Parse()

    if *genConfig {
        configWrite()
        os.Exit(1)
        return
    }

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            log.Printf("%v\n", err)
            usage( fmt.Sprintf("\nTry: %s --genconfig\n", me) )
            os.Exit(1)
        } else {
            log.Fatalf("Failed to parse config: %v\n", err)
            os.Exit(0)
        }
    }

    initLogger( viper.GetString("log.level") )
}

func configWrite(){
    viper.SafeWriteConfigAs( fmt.Sprintf("./%s", yamlFile) )
    log.Printf("Wrote example config to: \"./%s\", feel free to move to: %v", yamlFile, configSearchPaths[1:])
    log.Println(`

Alternatively, set the following environment variables:

export ` + envPrefix + `_HOST='127.0.0.1'
export ` + envPrefix + `_PORT='11099'
export ` + envPrefix + `_HISTPORT='11098'
export ` + envPrefix + `_USERNAME=''
export ` + envPrefix + `_PASSWORD=''
`)
    os.Exit(0)
}
