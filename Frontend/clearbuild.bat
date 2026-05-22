@echo off
setlocal
echo ===============================
echo  Limpando node_modules, lock e .next
echo ===============================

REM Verifica e remove a pasta node_modules
if exist "node_modules" (
    echo Removendo node_modules...
    rmdir /s /q "node_modules"
    echo node_modules removida.
) else (
    echo Pasta node_modules nao encontrada.
)

REM Verifica e remove o arquivo pnpm-lock.yaml
if exist "pnpm-lock.yaml" (
    echo Removendo pnpm-lock.yaml...
    del /f /q "pnpm-lock.yaml"
    echo Arquivo pnpm-lock.yaml removido.
) else (
    echo Arquivo pnpm-lock.yaml nao encontrado.
)

REM Verifica e remove a pasta .next
if exist ".next" (
    echo Removendo .next...
    rmdir /s /q ".next"
    echo .next removida.
) else (
    echo Pasta .next nao encontrada.
)

echo ===============================
echo       Operacao concluida!
echo ===============================
pause
endlocal
