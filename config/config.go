package config

import (
    "fmt"
    "os"
    "github.com/spf13/viper"
)

type Config struct {
    Url string
    User string
    Password string
    ApiKey string
    Upload struct {
        Folder int
        Comment string
    }
    Ls struct {
        Folder int
    }
    Autoupload struct {
        Folder int
        Comment string
        Keep bool
    }
}

func NewConfig(file string, section string) (*Config, error) {

    fmt.Println(file)
    if file != "" {
        viper.SetConfigName(file)
        viper.SetConfigType("yaml")
        viper.AddConfigPath("$HOME")
    } else {
        viper.SetConfigName(".seeddms-client.yaml")
        viper.SetConfigType("yaml")
        viper.AddConfigPath("/etc/seeddms-cups/")
        viper.AddConfigPath("$HOME")
        viper.AddConfigPath(".")
    }
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // No config file isn't neccesarily an error
        } else {
            return nil, fmt.Errorf("fatal error in config file: %w\n", err)
        }
    }

    if section == "" {
        section = os.Getenv("SEEDDMS_CLIENT")
        if section == "" {
//            fmt.Println("Environment varialbe SEEDDMS_CLIENT not set, using defaults")
            section = "default"
        }
    }

    cfgSection := viper.Sub(section)
    if cfgSection == nil {
        return nil, fmt.Errorf("section \"%s\" not found\n", section)
    }

    cfg := Config{}
    if err := cfgSection.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("Unable to decode config section into struct: %v\n", err)
    }

//    fmt.Printf("%#v", cfg);
    return &cfg, nil
}

