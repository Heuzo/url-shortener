#!/bin/bash

DB_PATH="./storage/storage.db"

if [ -f "$DB_PATH" ]; then
    if rm "$DB_PATH"; then
        echo "File deleted:  $DB_PATH"
    else
        echo "File delete error" >&2
        exit 1
    fi
else
    echo "File not found: $DB_PATH"
fi
