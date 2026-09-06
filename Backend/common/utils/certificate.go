package utils

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

// CertificateInfo representa os dados extraídos de um certificado digital
// PKCS#12 (.pfx/.p12) — usado pra popular a Configuração Fiscal da empresa
// automaticamente a partir do próprio arquivo, em vez do usuário digitar
// (o vencimento vem do certificado; nome/documento do titular servem pro
// usuário conferir visualmente se subiu o certificado certo).
type CertificateInfo struct {
	SubjectName     string // nome/razão social do titular do certificado
	SubjectDocument string // CPF/CNPJ extraído do Subject, quando o certificado segue o padrão ICP-Brasil
	NotBefore       time.Time
	NotAfter        time.Time
}

// ParsePKCS12Certificate decodifica um certificado .pfx/.p12 com a senha
// informada e extrai os dados relevantes. Erros de senha incorreta ou
// arquivo inválido são reportados de forma genérica — não vaza detalhe
// críptico de biblioteca pro usuário final.
func ParsePKCS12Certificate(fileBytes []byte, password string) (*CertificateInfo, error) {
	_, cert, err := pkcs12.Decode(fileBytes, password)
	if err != nil {
		return nil, errors.New("não foi possível ler o certificado: senha incorreta ou arquivo inválido")
	}
	if cert == nil {
		return nil, errors.New("nenhum certificado encontrado no arquivo")
	}

	name, document := splitICPBrasilSubject(cert.Subject.CommonName)
	return &CertificateInfo{
		SubjectName:     name,
		SubjectDocument: document,
		NotBefore:       cert.NotBefore,
		NotAfter:        cert.NotAfter,
	}, nil
}

// splitICPBrasilSubject separa o padrão comum de certificados e-CNPJ/e-CPF
// da ICP-Brasil ("RAZAO SOCIAL:12345678000199") em nome e documento. Se o
// certificado não seguir esse padrão, devolve o nome inteiro e documento
// vazio — nunca falha por causa disso, é só um extra pro usuário conferir.
func splitICPBrasilSubject(commonName string) (name string, document string) {
	idx := strings.LastIndex(commonName, ":")
	if idx == -1 {
		return commonName, ""
	}
	doc := commonName[idx+1:]
	if len(doc) != 11 && len(doc) != 14 { // nem CPF (11) nem CNPJ (14)
		return commonName, ""
	}
	return commonName[:idx], doc
}
