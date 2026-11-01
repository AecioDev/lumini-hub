package service

import (
	"encoding/json"
	"errors"
	"testing"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/common/utils"
)

// Regras de negócio do Fornecedor. Precisam de Postgres: sem TEST_DATABASE_URL são ignoradas.

func novoFornecedor(t *testing.T, s *SupplierService, nome, doc string) *domain.ApiSupplier {
	t.Helper()
	f, err := s.CreateSupplier(domain.CreateSupplierRequest{FirstName: nome, PersonType: "J", DocumentNumber: doc, IsActive: true}, 1)
	if err != nil {
		t.Fatalf("criar fornecedor %q: %v", nome, err)
	}
	return f
}

func TestFornecedor_CriacaoGravaSoDigitos(t *testing.T) {
	s := NewSupplierService(newTestUow(t))

	f, err := s.CreateSupplier(domain.CreateSupplierRequest{FirstName: "Madeireira", PersonType: "J", DocumentNumber: "11.222.333/0001-81", IsActive: true}, 7)
	if err != nil {
		t.Fatalf("criar: %v", err)
	}
	if f.DocumentNumber != "11222333000181" {
		t.Errorf("documento gravado = %q, quer só os dígitos 11222333000181", f.DocumentNumber)
	}
}

// A autoria vem sempre do usuário autenticado: um "created_by" no corpo da requisição é ignorado.
func TestFornecedor_AutoriaVemDoLogin(t *testing.T) {
	s := NewSupplierService(newTestUow(t))

	var req domain.CreateSupplierRequest
	corpo := `{"first_name":"Madeireira","person_type":"J","document_number":"11222333000181","is_active":true,"created_by":999}`
	if err := json.Unmarshal([]byte(corpo), &req); err != nil {
		t.Fatalf("corpo: %v", err)
	}
	f, err := s.CreateSupplier(req, 7)
	if err != nil {
		t.Fatalf("criar: %v", err)
	}
	detalhe, err := s.GetSupplierByID(f.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if detalhe.CreatedBy == nil || detalhe.CreatedBy.ID != 7 {
		t.Errorf("CreatedBy = %+v, quer o usuário autenticado (7), nunca o 999 do corpo", detalhe.CreatedBy)
	}
}

func TestFornecedor_DocumentoDuplicadoENomeCurto(t *testing.T) {
	s := NewSupplierService(newTestUow(t))
	novoFornecedor(t, s, "Madeireira", "11.222.333/0001-81")

	for _, repetido := range []string{"11222333000181", "11.222.333/0001-81"} {
		_, err := s.CreateSupplier(domain.CreateSupplierRequest{FirstName: "Outra", PersonType: "J", DocumentNumber: repetido, IsActive: true}, 1)
		if err == nil {
			t.Fatalf("documento %q repetido deveria ser recusado", repetido)
		}
		if !hasAnyValidationError(err) {
			t.Errorf("erro de %q = %v, quer erro de validação (mensagem para o usuário)", repetido, err)
		}
	}

	_, err := s.CreateSupplier(domain.CreateSupplierRequest{FirstName: "Al", PersonType: "J", DocumentNumber: "99999999000199", IsActive: true}, 1)
	exigeErroNoCampo(t, err, "first_name")
}

func TestFornecedor_AtualizacaoEExclusao(t *testing.T) {
	s := NewSupplierService(newTestUow(t))
	a := novoFornecedor(t, s, "Madeireira", "11.222.333/0001-81")
	novoFornecedor(t, s, "Ferragens", "22.333.444/0001-92")

	// o próprio documento (mascarado) não é duplicidade; o de outro fornecedor é
	if _, err := s.UpdateSupplier(a.ID, domain.UpdateSupplierRequest{FirstName: "Madeireira Sul", DocumentNumber: "11.222.333/0001-81", IsActive: true}, 5); err != nil {
		t.Fatalf("manter o próprio documento: %v", err)
	}
	_, err := s.UpdateSupplier(a.ID, domain.UpdateSupplierRequest{FirstName: "Madeireira", DocumentNumber: "22.333.444/0001-92", IsActive: true}, 5)
	if err == nil || !hasAnyValidationError(err) {
		t.Errorf("documento de outro fornecedor: erro = %v, quer erro de validação", err)
	}

	detalhe, err := s.GetSupplierByID(a.ID)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if detalhe.UpdatedBy == nil || detalhe.UpdatedBy.ID != 5 {
		t.Errorf("UpdatedBy = %+v, quer o usuário 5", detalhe.UpdatedBy)
	}

	if err := s.DeleteSupplier(a.ID); err != nil {
		t.Fatalf("excluir: %v", err)
	}
	if _, err := s.GetSupplierByID(a.ID); !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("fornecedor excluído: erro = %v, quer ErrNotFound", err)
	}
	if err := s.DeleteSupplier(a.ID); !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("excluir de novo: erro = %v, quer ErrNotFound", err)
	}
}

func TestFornecedor_BuscaFiltradaEPaginada(t *testing.T) {
	s := NewSupplierService(newTestUow(t))
	novoFornecedor(t, s, "Carvalho Madeiras", "33.333.333/0001-33")
	novoFornecedor(t, s, "Alfa Ferragens", "11.111.111/0001-11")
	novoFornecedor(t, s, "Beta Tintas", "22.222.222/0001-22")

	asc := true
	pagina2, tamanho := 2, 2
	r, err := s.GetSuppliersByFilter(domain.SupplierFilterRequest{OrderByColumn: "first_name", IsAsc: &asc, PageNo: &pagina2, PageSize: &tamanho})
	if err != nil {
		t.Fatalf("busca: %v", err)
	}
	if len(r.Supplier) != 1 || r.Supplier[0].FirstName != "Carvalho Madeiras" {
		t.Errorf("página 2 (de 2 por página) = %+v, quer só Carvalho Madeiras", r.Supplier)
	}
	if r.Pagination.TotalRows != 3 || r.Pagination.TotalPages != 2 {
		t.Errorf("total = %d, páginas = %d, quer 3 e 2", r.Pagination.TotalRows, r.Pagination.TotalPages)
	}

	r, _ = s.GetSuppliersByFilter(domain.SupplierFilterRequest{FirstName: "FERRAG"})
	if len(r.Supplier) != 1 || r.Supplier[0].FirstName != "Alfa Ferragens" {
		t.Errorf("first_name FERRAG = %+v", r.Supplier)
	}
}

// CFG-050: o documento de um fornecedor excluído volta a ficar livre (índice único parcial).
func TestFornecedor_RecadastroDepoisDaExclusao(t *testing.T) {
	s := NewSupplierService(newTestUow(t))
	antigo := novoFornecedor(t, s, "Casa Alfa", "11.222.333/0001-81")
	if err := s.DeleteSupplier(antigo.ID); err != nil {
		t.Fatalf("excluir: %v", err)
	}
	if _, err := s.CreateSupplier(domain.CreateSupplierRequest{FirstName: "Casa Beta", PersonType: "J", DocumentNumber: "11.222.333/0001-81"}, 1); err != nil {
		t.Fatalf("recadastrar o documento de um fornecedor excluído: %v", err)
	}
}
