#!/bin/sh

# Run migrations
go run main.go db:migrate up

# Start the application
go run main.go serve
