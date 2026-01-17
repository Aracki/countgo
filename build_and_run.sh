#!/bin/bash
go build -o bin/countgo ./cmd/aracki
nohup ./bin/countgo -y=false > logs/countgo.log 2>&1 &
echo "Started countgo (PID: $!)"
