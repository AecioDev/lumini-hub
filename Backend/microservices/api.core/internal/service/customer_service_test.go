package service

import (
	"errors"
	"testing"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/common/utils"
)

// Regras de negócio do Cliente. Precisam de Postgres: sem TEST_DATABASE_URL são ignoradas.

func novoCliente(t *testing.T, s *CustomerService, nome, doc string) *domain.ApiCustomer {
	t.Helper()
	c, err := s.CreateCustomer(domain.CreateCustomerRequest{FirstName: nome, PersonType: "F", DocumentNumber: doc}, 1)
	if err != nil {
		t.Fatalf("criar cliente %q: %v", nome, err)
	}
	return c
}

func TestCliente_CriacaoGravaSoDigitosEAutor(t *testing.T) {
	s := NewCustomerService(newTestUow(t))

	c, err := s.CreateCustomer(domain.CreateCustomerRequest{FirstName: "Maria", PersonType: "F", DocumentNumber: "123.456.789-09"}, 7)
	if err != nil {
		t.Fatalf("criar: %v", err)
	}
	if c.DocumentNumber != "12345678909" {
		t.Errorf("documento gravado = %q, quer só os dígitos 12345678909", c.DocumentNumber)
	}
	if !c.IsActive {
		t.Error("cliente novo deveria nascer ativo")
	}
	detalhe, err := s.GetCustomerByID(c.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if detalhe.CreatedBy == nil || detalhe.CreatedBy.ID != 7 {
		t.Errorf("CreatedBy = %+v, quer o usuário 7", detalhe.CreatedBy)
	}
}

func TestCliente_DocumentoDuplicadoComOuSemMascara(t *testing.T) {
	s := NewCustomerService(newTestUow(t))
	novoCliente(t, s, "Maria", "123.456.789-09")

	for _, repetido := range []string{"12345678909", "123.456.789-09"} {
		_, err := s.CreateCustomer(domain.CreateCustomerRequest{FirstName: "Outra", PersonType: "F", DocumentNumber: repetido}, 1)
		exigeErroNoCampo(t, err, "document")
	}
}

func TestCliente_NomeMuitoCurto(t *testing.T) {
	s := NewCustomerService(newTestUow(t))

	_, err := s.CreateCustomer(domain.CreateCustomerRequest{FirstName: "Al", PersonType: "F", DocumentNumber: "11122233344"}, 1)
	exigeErroNoCampo(t, err, "first_name")
	_, err = s.CreateCustomer(domain.CreateCustomerRequest{FirstName: "Alice", LastName: "Li", PersonType: "F", DocumentNumber: "11122233344"}, 1)
	exigeErroNoCampo(t, err, "last_name")
}

func TestCliente_Atualizacao(t *testing.T) {
	s := NewCustomerService(newTestUow(t))
	a := novoCliente(t, s, "Ana", "111.111.111-11")
	novoCliente(t, s, "Bia", "222.222.222-22")

	// manter o próprio documento (mascarado) não é duplicidade
	if _, err := s.UpdateCustomer(a.ID, domain.UpdateCustomerRequest{FirstName: "Ana Maria", DocumentNumber: "111.111.111-11", IsActive: true}); err != nil {
		t.Fatalf("manter o próprio documento: %v", err)
	}
	// o de outro cliente é
	_, err := s.UpdateCustomer(a.ID, domain.UpdateCustomerRequest{FirstName: "Ana", DocumentNumber: "222.222.222-22", IsActive: true})
	exigeErroNoCampo(t, err, "document_number")

	_, err = s.UpdateCustomer(9999, domain.UpdateCustomerRequest{FirstName: "Fantasma", DocumentNumber: "333.333.333-33", IsActive: true})
	exigeErroNoCampo(t, err, "id")

	atualizado, err := s.UpdateCustomer(a.ID, domain.UpdateCustomerRequest{FirstName: "Ana Paula", DocumentNumber: "11111111111", IsActive: false, Notes: "inativa"})
	if err != nil {
		t.Fatalf("atualizar: %v", err)
	}
	if atualizado.FirstName != "Ana Paula" || atualizado.IsActive {
		t.Errorf("atualizado = %+v, quer nome novo e inativo", atualizado)
	}
}

func TestCliente_Exclusao(t *testing.T) {
	s := NewCustomerService(newTestUow(t))
	c := novoCliente(t, s, "Ana", "111.111.111-11")

	if err := s.DeleteCustomer(c.ID); err != nil {
		t.Fatalf("excluir: %v", err)
	}
	if _, err := s.GetCustomerByID(c.ID); !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("cliente excluído: erro = %v, quer ErrNotFound", err)
	}
	if err := s.DeleteCustomer(c.ID); !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("excluir de novo: erro = %v, quer ErrNotFound", err)
	}
}

func TestCliente_BuscaFiltradaPaginadaEOrdenada(t *testing.T) {
	s := NewCustomerService(newTestUow(t))
	novoCliente(t, s, "Carlos", "333.333.333-33")
	novoCliente(t, s, "Alice", "111.111.111-11")
	juridica, err := s.CreateCustomer(domain.CreateCustomerRequest{FirstName: "Beatriz", PersonType: "J", DocumentNumber: "11222333000181", CompanyName: "Beatriz Móveis"}, 1)
	if err != nil {
		t.Fatalf("criar PJ: %v", err)
	}
	if _, err := s.UpdateCustomer(juridica.ID, domain.UpdateCustomerRequest{FirstName: "Beatriz", DocumentNumber: "11222333000181", CompanyName: "Beatriz Móveis", IsActive: false}); err != nil {
		t.Fatalf("inativar: %v", err)
	}

	nomes := func(r *domain.ApiCustomerListPaginated) []string {
		var out []string
		for _, c := range r.Customer {
			out = append(out, c.FirstName)
		}
		return out
	}
	igual := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
	boolPtr := func(b bool) *bool { return &b }
	intPtr := func(i int) *int { return &i }

	// sem filtro: o mais novo primeiro
	r, err := s.GetCustomersByFilter(domain.CustomerFilterRequest{})
	if err != nil {
		t.Fatalf("sem filtro: %v", err)
	}
	if got := nomes(r); !igual(got, []string{"Beatriz", "Alice", "Carlos"}) || r.Pagination.TotalRows != 3 {
		t.Errorf("sem filtro = %v (total %d)", got, r.Pagination.TotalRows)
	}

	// texto parcial, sem diferenciar maiúsculas
	r, _ = s.GetCustomersByFilter(domain.CustomerFilterRequest{FirstName: "ARL"})
	if got := nomes(r); !igual(got, []string{"Carlos"}) {
		t.Errorf("first_name ARL = %v", got)
	}
	r, _ = s.GetCustomersByFilter(domain.CustomerFilterRequest{PersonType: "J"})
	if got := nomes(r); !igual(got, []string{"Beatriz"}) {
		t.Errorf("person_type J = %v", got)
	}
	r, _ = s.GetCustomersByFilter(domain.CustomerFilterRequest{IsActive: boolPtr(false)})
	if got := nomes(r); !igual(got, []string{"Beatriz"}) {
		t.Errorf("inativos = %v", got)
	}

	// ordenação pedida pelo cliente da API
	r, _ = s.GetCustomersByFilter(domain.CustomerFilterRequest{OrderByColumn: "first_name", IsAsc: boolPtr(true)})
	if got := nomes(r); !igual(got, []string{"Alice", "Beatriz", "Carlos"}) {
		t.Errorf("ordenado por first_name asc = %v, quer [Alice Beatriz Carlos]", got)
	}
	r, _ = s.GetCustomersByFilter(domain.CustomerFilterRequest{OrderByColumn: "first_name", IsAsc: boolPtr(false)})
	if got := nomes(r); !igual(got, []string{"Carlos", "Beatriz", "Alice"}) {
		t.Errorf("ordenado por first_name desc = %v, quer [Carlos Beatriz Alice]", got)
	}

	// paginação: o total é o da busca, não o da página
	r, _ = s.GetCustomersByFilter(domain.CustomerFilterRequest{OrderByColumn: "first_name", IsAsc: boolPtr(true), PageNo: intPtr(2), PageSize: intPtr(2)})
	if got := nomes(r); !igual(got, []string{"Carlos"}) || r.Pagination.TotalRows != 3 || r.Pagination.TotalPages != 2 {
		t.Errorf("página 2 de 2 = %v (total %d, páginas %d)", got, r.Pagination.TotalRows, r.Pagination.TotalPages)
	}

	// coluna desconhecida ou maliciosa cai na ordem padrão, sem erro
	r, err = s.GetCustomersByFilter(domain.CustomerFilterRequest{OrderByColumn: "id; DROP TABLE customers"})
	if err != nil || r.Pagination.TotalRows != 3 {
		t.Errorf("coluna inválida: erro = %v, total = %d", err, r.Pagination.TotalRows)
	}
}
