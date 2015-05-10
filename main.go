package main

import (
    "log"
    "sync"
    "github.com/vaughan0/go-ini"
    "github.com/biogo/store/interval"
    "os"
    "os/signal"
    "syscall"
    "net"
    "fmt"
)

var mutex sync.RWMutex = sync.RWMutex{}
var Data interval.IntTree
var generalConfig GeneralConfig
var driver Driver

func initialLoad() error {
    log.Print("Loading config file... ")
    file, err := ini.LoadFile("settings.conf")
    if err != nil {
        return err
    }
    log.Println("done")
    log.Print("Parsing general config... ")
    generalConfig, err = NewGeneralConfig(file)
    if err != nil {
        return err
    }
    log.Println("done")

    log.Printf("Got driver %s\n", generalConfig.Driver)
    log.Print("Driver initialization... ")
    driver, err = InitDriver(generalConfig, file)
    if err != nil {
        return err
    }
    log.Println("done");
    log.Print("Loading data... ")
    tree, err := driver.Load()
    if err != nil {
        return err
    }
    mutex.Lock()
    defer mutex.Unlock()
    Data = tree
    log.Println("done")
    return nil
}

func main() {

    log.Println("Initial load")
    err := initialLoad()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Initial load done")

    usr1_channel := make(chan os.Signal, 1)
    signal.Notify(usr1_channel, syscall.SIGUSR1)

    log.Println("Listening to SIGUSR1 for reload")
    go func() {
        for {
            <-usr1_channel
            log.Println("Got SIGUSR1. Starting reload")
            err := initialLoad()
            if err != nil {
                log.Print(err)
            }
            log.Println("Reload done")
        }
    }()
    log.Printf("Listening tcp socket on: %s:%d\n", generalConfig.Host, generalConfig.Port)
    ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", generalConfig.Host, generalConfig.Port))
    if err != nil {
        log.Fatal(err)
    }
    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go handleConnection(conn)
    }
}