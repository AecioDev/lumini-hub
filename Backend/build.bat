@echo off
setlocal

echo O que você deseja fazer?
echo 1 - Compilar
echo 2 - Executar
echo 3 - Compilar e Executar
set /p opcao=Digite o número da opção desejada: 

if "%opcao%"=="1" goto compilar
if "%opcao%"=="2" goto executar
if "%opcao%"=="3" goto compilar_executar

echo Opcao invalida.
goto fim

:compilar
echo Compilando...
go build -o bin/main.exe ./cmd/api
if %errorlevel% neq 0 (
    echo Erro ao compilar o programa.
    exit /b %errorlevel%
)
echo Build concluído!
goto fim

:executar
echo Executando...
bin\main.exe
goto fim

:compilar_executar
echo Compilando...
go build -o bin/main.exe ./cmd/api
if %errorlevel% neq 0 (
    echo Erro ao compilar o programa.
    exit /b %errorlevel%
)
echo Build concluído!
echo Executando...
bin\main.exe
goto fim

:fim
endlocal
