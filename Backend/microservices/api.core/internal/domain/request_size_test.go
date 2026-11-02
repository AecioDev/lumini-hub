package domain

import (
	"strings"
	"testing"

	"lumini-hub/common/testutil"

	"github.com/gin-gonic/gin/binding"
)

// A coluna tem tamanho no banco, então a request precisa recusar o excesso antes do banco: todo campo
// com gorm "size:N" que aparece numa request de criar/editar deve ter binding "max=N". Ao criar uma
// entidade nova no api.core, acrescente-a aqui.
func TestRequests_MaxAlinhadoAoTamanhoDaColuna(t *testing.T) {
	testutil.AssertRequestMaxMatchesColumnSize(t, Customer{}, CreateCustomerRequest{}, UpdateCustomerRequest{})
	testutil.AssertRequestMaxMatchesColumnSize(t, Supplier{}, CreateSupplierRequest{}, UpdateSupplierRequest{})
	testutil.AssertRequestMaxMatchesColumnSize(t, Company{}, CreateCompanyRequest{}, UpdateCompanyRequest{})
	testutil.AssertRequestMaxMatchesColumnSize(t, CompanyFiscalConfig{}, CreateCompanyFiscalConfigRequest{}, UpdateCompanyFiscalConfigRequest{})
	// Servidor SMTP: as colunas vivem em SMTPSettings, embutida no sistema e na configuração de cada empresa.
	testutil.AssertRequestMaxMatchesColumnSize(t, SMTPSettings{}, SaveSMTPRequest{})
	testutil.AssertRequestMaxMatchesColumnSize(t, SystemSettings{}, SaveSystemEmailRequest{})
}

// Comportamento de verdade: 100 caracteres passam, 101 são recusados (conta caracteres, não bytes).
func TestRequests_ExcessoDeTamanhoEhRecusadoAntesDoBanco(t *testing.T) {
	longText := func(n int) string { return strings.Repeat("ã", n) }

	validCustomer := func() CreateCustomerRequest {
		return CreateCustomerRequest{FirstName: "Ana", PersonType: "F", DocumentNumber: "52998224725"}
	}
	validSupplier := func() CreateSupplierRequest {
		return CreateSupplierRequest{FirstName: "Ana", PersonType: "F", DocumentNumber: "52998224725", IsActive: true}
	}

	type fieldCase struct {
		name string
		set  func(text string) (customer any, supplier any)
		max  int
	}
	fields := []fieldCase{
		{"first_name", func(s string) (any, any) {
			c, f := validCustomer(), validSupplier()
			c.FirstName, f.FirstName = s, s
			return c, f
		}, 100},
		{"last_name", func(s string) (any, any) {
			c, f := validCustomer(), validSupplier()
			c.LastName, f.LastName = s, s
			return c, f
		}, 100},
		{"company_name", func(s string) (any, any) {
			c, f := validCustomer(), validSupplier()
			c.CompanyName, f.CompanyName = s, s
			return c, f
		}, 100},
		{"document_number", func(s string) (any, any) {
			c, f := validCustomer(), validSupplier()
			c.DocumentNumber, f.DocumentNumber = s, s
			return c, f
		}, 20},
		{"notes", func(s string) (any, any) {
			c, f := validCustomer(), validSupplier()
			c.Notes, f.Notes = s, s
			return c, f
		}, 255},
	}

	for _, fc := range fields {
		customerOK, supplierOK := fc.set(longText(fc.max))
		for name, req := range map[string]any{"cliente": customerOK, "fornecedor": supplierOK} {
			if err := binding.Validator.ValidateStruct(req); err != nil {
				t.Errorf("%s %s com %d caracteres deveria passar: %v", name, fc.name, fc.max, err)
			}
		}
		customerBad, supplierBad := fc.set(longText(fc.max + 1))
		for name, req := range map[string]any{"cliente": customerBad, "fornecedor": supplierBad} {
			if err := binding.Validator.ValidateStruct(req); err == nil {
				t.Errorf("%s %s com %d caracteres deveria ser recusado", name, fc.name, fc.max+1)
			}
		}
	}
}
