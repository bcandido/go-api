#!/usr/bin/env bash

function go_get() {
    local dependency=$1

    echo "getting $dependency"
    go get $dependency

}

echo "==== building dependencies..."

go_get "github.com/gima/govalid/v1"
go_get "github.com/gorilla/mux"
go_get "github.com/lib/pq"
go_get "github.com/op/go-logging"
go_get "github.com/spf13/viper"

echo "==== dependencies built"
