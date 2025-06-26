#!/bin/bash

# Setup database
createdb library_db

# Run the application
go run cmd/server/main.go
