# ZKits Logger Library #

[![ZKits](https://img.shields.io/badge/ZKits-Library-f3c)](https://github.com/edoger/zkits-logger)
[![Build Status](https://travis-ci.org/edoger/zkits-logger.svg?branch=master)](https://travis-ci.org/edoger/zkits-logger)
[![Coverage Status](https://coveralls.io/repos/github/edoger/zkits-logger/badge.svg?branch=master)](https://coveralls.io/github/edoger/zkits-logger?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/cb497bf703f44950afb43b51b3a0e581)](https://www.codacy.com/manual/edoger/zkits-logger?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=edoger/zkits-logger&amp;utm_campaign=Badge_Grade)
[![Golang Version](https://img.shields.io/badge/golang-1.13+-orange)](https://github.com/edoger/zkits-logger)

## About ##

This package is a library of ZKits project. 
This library provides structured JSON format logger.
At the same time, it supports 7 levels of logs and provides complete hook function to assist log scheduling.

## Usage ##

 1. Import package.
 
    ```sh
    go get -u -v github.com/edoger/zkits-logger
    ```

 2. Example.
    ```go
    package main
    
    import (
        "github.com/edoger/zkits-logger"
    )
    
    func main() {
        // Creates a logger instance with the specified name.
        // By default, logs are output to standard output.
        // In addition, logger only record INFO or higher level logs.
        // You can modify these settings by calling methods.
        // log.SetOutput()
        // log.SetLevel()
        log := logger.New("test")
    
        // Record an info level log.
        // {"level":"info","logger":"test","message":"It's log!","time":"2006-01-02 15:04:05","fields":{}}
        log.Info("It's log!")
    }
    ```

## License ##

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)
