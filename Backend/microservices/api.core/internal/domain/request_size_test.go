package domain

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

// columnSize devolve o N de uma tag gorm "size:N", se houver.
func columnSize(field reflect.StructField) (int, bool) {
	for _, part := range strings.Split(field.Tag.Get("gorm"), ";") {
		if value, found := strings.CutPrefix(strings.TrimSpace(part), "size:"); found {
			if size, err := strconv.Atoi(value); err == nil {
				return size, true
			}
		}
	}
	return 0, false
}

func jsonName(field reflect.StructField) string {
	return strings.Split(field.Tag.Get("json"), ",")[0]
}

// A coluna tem tamanho no banco, então a request precisa recusar o excesso antes do banco: todo campo
// com gorm "size:N" que aparece numa request de criar/editar deve ter binding "max=N".
func TestRequests_MaxAlinhadoAoTamanhoDaColuna(t *testing.T) {
	cases := []struct {
		entity   any
		requests []any
	}{
		{Customer{}, []any{CreateCustomerRequest{}, UpdateCustomerRequest{}}},
		{Supplier{}, []any{CreateSupplierRequest{}, UpdateSupplierRequest{}}},
	}

	for _, tc := range cases {
		entityType := reflect.TypeOf(tc.entity)
		sizes := map[string]int{}
		for i := 0; i < entityType.NumField(); i++ {
			field := entityType.Field(i)
			if size, ok := columnSize(field); ok {
				sizes[jsonName(field)] = size
			}
		}
		if len(sizes) == 0 {
			t.Fatalf("%s: nenhum campo com size no banco, o teste não está olhando nada", entityType.Name())
		}

		for _, request := range tc.requests {
			requestType := reflect.TypeOf(request)
			covered := 0
			for i := 0; i < requestType.NumField(); i++ {
				field := requestType.Field(i)
				size, ok := sizes[jsonName(field)]
				if !ok {
					continue
				}
				covered++
				want := "max=" + strconv.Itoa(size)
				if !strings.Contains(field.Tag.Get("binding"), want) {
					t.Errorf("%s.%s: binding %q deveria conter %q (o tamanho da coluna em %s)",
						requestType.Name(), field.Name, field.Tag.Get("binding"), want, entityType.Name())
				}
			}
			if covered == 0 {
				t.Errorf("%s não tem nenhum campo em comum com %s", requestType.Name(), entityType.Name())
			}
		}
	}
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
