@echo off
echo ===================================================
echo Iniciando Microsservicos e API Gateway - Lumini Hub
echo ===================================================
echo.

echo [1/4] Iniciando api.auth na porta 4001...
start "Lumini Hub - api.auth (Porta 4001)" cmd /k "cd microservices\api.auth && go run main.go"

echo [2/4] Iniciando api.core na porta 4002...
start "Lumini Hub - api.core (Porta 4002)" cmd /k "cd microservices\api.core && go run main.go"

echo [3/4] Iniciando api.integrations na porta 4007...
start "Lumini Hub - api.integrations (Porta 4007)" cmd /k "cd microservices\api.integrations && go run main.go"

echo [4/4] Iniciando api.gateway (Porta 4000)...
start "Lumini Hub - api.gateway (Porta 4000)" cmd /k "cd microservices\api.gateway && go run main.go"

echo.
echo ===================================================
echo Todos os servicos foram abertos em janelas separadas.
echo API Gateway rodando e escutando em http://localhost:4000
echo ===================================================
echo.
pause
