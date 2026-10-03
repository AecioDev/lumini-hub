package main

import (
	"os"
	"os/exec"
	"testing"
)

// Os serviços são iniciados com `go run main.go` (run_services.bat), que compila SÓ o main.go.
// Código movido para outro arquivo do pacote passa em `go build .` e nos demais testes, mas
// quebra esse comando com "undefined". Este teste roda o comando de verdade.
func TestMainCompilaSozinho(t *testing.T) {
	out, err := exec.Command("go", "build", "-o", os.DevNull, "main.go").CombinedOutput()
	if err != nil {
		t.Fatalf("`go run main.go` (run_services.bat) compila só o main.go e falhou: %v\n%s\n"+
			"Mantenha tudo que o main.go usa dentro dele (ou em pacotes importados), não em outro arquivo do package main.", err, out)
	}
}
