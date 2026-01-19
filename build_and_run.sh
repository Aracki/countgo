#!/bin/bash
pkill -f countgo 2>/dev/null
go build -o bin/countgo ./cmd/aracki
nohup ./bin/countgo -y=false > logs/countgo.log 2>&1 &
echo "Started countgo (PID: $!)"
