#!/bin/bash

DB_PATH="./storage/storage.db"

if [ -f "$DB_PATH" ]; then
    if rm "$DB_PATH"; then
        echo "Файл storage.db удалён."
    else
        echo "Ошибка при удалении файла." >&2
        exit 1
    fi
else
    echo "Файл storage.db не найден."
fi
