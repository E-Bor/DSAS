#!/bin/bash

# base compile dir
BASE_DIR="internal/data_sources"

# find files by mask *_report.go
find "$BASE_DIR" -name "*_report.go" | while read -r report_file; do
    dir=$(dirname "$report_file")
    filename=$(basename "$report_file" .go)
    output_file="$dir/compiled/${filename}.so"
    echo "Compiling $report_file в $output_file"
    go build -buildmode=plugin -o "$output_file" "$report_file"
    if [[ $? -eq 0 ]]; then
        echo "Compiling success: $output_file"
    else
        echo "Error while compiling: $report_file"
    fi
done
