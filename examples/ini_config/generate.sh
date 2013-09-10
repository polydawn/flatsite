#!/bin/bash
rm -r conf
INPUT_DIR="conf_tmpl" \
OUTPUT_DIR="conf" \
go run ../../flatsite.go
