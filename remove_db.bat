@echo off
chcp 65001 > nul

set "file=.\storage\storage.db"

if exist "%file%" (
    del "%file%"
    echo File deleted:  %file%
) else (
    echo File not found: %file%
)

pause