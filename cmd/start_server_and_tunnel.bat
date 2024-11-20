@echo off

REM Установить параметры для vk-tunnel
set TUNNEL_CMD=vk-tunnel --insecure=1 --http-protocol=http --ws-protocol=ws --host=localhost --port=8080 --timeout=5000

REM Установить путь к директории проекта
set PROJECT_DIR=R:\ProjectsGo\InstaSpace\cmd

REM Запуск vk-tunnel в отдельной консоли
start cmd /k "%TUNNEL_CMD%"

REM Перейти в директорию с проектом
cd /d "%PROJECT_DIR%"

REM Запуск main.go в отдельной консоли
start cmd /k "go run main.go"
