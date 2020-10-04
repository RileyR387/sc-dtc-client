package main

import (
    "fmt"
    "log"
    "os"
//    "os/user"
    "path/filepath"
    "strings"
    "github.com/pborman/getopt/v2"
    "github.com/spf13/viper"

    "github.com/RileyR387/sc-dtc-client/dtc"
)

var (
    me = filepath.Base(os.Args[0])
    yamlFile = fmt.Sprintf("%s.yaml", me)
    envPrefix = "DTCCLIENT"
    configSearchPaths = []string {".", "./etc", "$HOME/.dtc-client-go/", "$HOME/etc", "/etc"}
    version = "undefined"

    //dtcServerHost = getopt.StringLong("host", 'h', "", "DTC Server Host")
    //dtcServerPort = getopt.Int("port", 'p', "", "DTC Server Port")
    //dtcServerUsername = getopt.StringLong("username", 'u', "", "DTC Username")
    //dtcServerPassword = getopt.StringLong("password", 'p', "", "DTC Password")
    //singleThread  = getopt.BoolLong("singlethreaded", 's', "Disable threading")
    genConfig = getopt.BoolLong("genconfig", 'x', "Write example config to \"./" + yamlFile + "\"")
)

func init() {
    log.SetPrefix("[main] ")
    viper.SetConfigName(yamlFile)
    viper.SetConfigType("yaml")
    viper.SetEnvPrefix(envPrefix)
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()
    for _, p := range configSearchPaths {
        viper.AddConfigPath(p)
    }
    /*
    user, err := user.Current()
    if err != nil {
        panic( err )
    }
    */
    viper.SetDefault("dtc.Host", "127.0.0.1")
    viper.SetDefault("dtc.Port", "11099")
    viper.SetDefault("dtc.HistPort", "11098")
    viper.SetDefault("dtc.Username", "")
    viper.SetDefault("dtc.Password", "")
}

func usage(msg ...string) {
    if len(msg) > 0 {
        fmt.Fprintf(os.Stderr, "%s\n", msg[0])
    }
    // strip off the first line of generated usage
    b := &strings.Builder{}
    getopt.PrintUsage(b)
    u := strings.SplitAfterN(b.String(), "\n", 2)
    fmt.Printf(`Usage: %s

Activity log is written to STDERR.

OPTIONS
%s
`, me, u[1])

    os.Exit(1)
}

func main() {
    getopt.SetUsage(func() { usage() })
    getopt.Parse()
    if *genConfig {
        configWrite()
    }
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            log.Printf("%v\n", err)
            usage( fmt.Sprintf("\nTry: %s --genconfig\n", me) )
            os.Exit(1)
        } else {
            log.Fatalf("Failed to parse config: %v\n", err)
        }
    }
    args := dtc.ConnectArgs {
        viper.GetString("dtc.Host"),
        viper.GetString("dtc.Port"),
        viper.GetString("dtc.HistPort"),
        viper.GetString("dtc.Username"),
        viper.GetString("dtc.Password"),
    }

    dtc.Connect(args)
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
