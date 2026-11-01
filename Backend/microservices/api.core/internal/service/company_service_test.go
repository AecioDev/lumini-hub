package service

import (
	"errors"
	"testing"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/common/utils"
)

// Regras de negócio da Empresa (matriz e filiais). Precisam de Postgres: sem TEST_DATABASE_URL são
// ignoradas.

func criarMatriz(t *testing.T, s *CompanyService, nome, cnpj string) uint {
	t.Helper()
	c, err := s.CreateCompany(domain.CreateCompanyRequest{LegalName: nome, CNPJ: cnpj}, 1)
	if err != nil {
		t.Fatalf("criar matriz %q: %v", nome, err)
	}
	return c.ID
}

func TestEmpresa_MatrizExigeCNPJ(t *testing.T) {
	s := NewCompanyService(newTestUow(t))

	_, err := s.CreateCompany(domain.CreateCompanyRequest{LegalName: "Matriz Sem CNPJ"}, 1)
	exigeErroNoCampo(t, err, "cnpj")
}

func TestEmpresa_CNPJGravadoSoComDigitosEUnico(t *testing.T) {
	s := NewCompanyService(newTestUow(t))

	c, err := s.CreateCompany(domain.CreateCompanyRequest{LegalName: "Alfa Ltda", CNPJ: "12.345.678/0001-99"}, 1)
	if err != nil {
		t.Fatalf("criar: %v", err)
	}
	if c.CNPJ == nil || *c.CNPJ != "12345678000199" {
		t.Fatalf("CNPJ gravado = %v, quer só os dígitos 12345678000199", c.CNPJ)
	}
	if !c.IsActive {
		t.Error("empresa nova deveria nascer ativa")
	}

	// o mesmo CNPJ, mascarado ou não, é duplicidade
	for _, repetido := range []string{"12345678000199", "12.345.678/0001-99"} {
		_, err = s.CreateCompany(domain.CreateCompanyRequest{LegalName: "Outra Ltda", CNPJ: repetido}, 1)
		if msg := exigeErroNoCampo(t, err, "cnpj"); msg != "CNPJ já está em uso" {
			t.Errorf("mensagem = %q", msg)
		}
	}
}

func TestEmpresa_FilialPodeNaoTerCNPJ(t *testing.T) {
	s := NewCompanyService(newTestUow(t))
	matriz := criarMatriz(t, s, "Matriz Ltda", "11222333000181")

	// duas filiais sem CNPJ não podem colidir no índice único (NULL nunca é igual a NULL)
	for _, nome := range []string{"Filial Um", "Filial Dois"} {
		f, err := s.CreateCompany(domain.CreateCompanyRequest{ParentID: &matriz, LegalName: nome}, 1)
		if err != nil {
			t.Fatalf("criar %s: %v", nome, err)
		}
		if f.CNPJ != nil {
			t.Errorf("%s: CNPJ = %q, quer nil (nunca string vazia)", nome, *f.CNPJ)
		}
		if f.ParentID == nil || *f.ParentID != matriz {
			t.Errorf("%s: ParentID = %v, quer %d", nome, f.ParentID, matriz)
		}
	}
}

func TestEmpresa_VinculoComEmpresaInexistenteOuCiclo(t *testing.T) {
	s := NewCompanyService(newTestUow(t))

	_, err := s.CreateCompany(domain.CreateCompanyRequest{ParentID: uintPtr(9999), LegalName: "Filial Órfã"}, 1)
	exigeErroNoCampo(t, err, "parent_id")

	a := criarMatriz(t, s, "Empresa A", "11111111000111")
	b, err := s.CreateCompany(domain.CreateCompanyRequest{ParentID: &a, LegalName: "Empresa B"}, 1)
	if err != nil {
		t.Fatalf("criar B: %v", err)
	}
	c, err := s.CreateCompany(domain.CreateCompanyRequest{ParentID: &b.ID, LegalName: "Empresa C"}, 1)
	if err != nil {
		t.Fatalf("criar C: %v", err)
	}

	// ela mesma como pai
	_, err = s.UpdateCompany(a, domain.UpdateCompanyRequest{ParentID: &a, LegalName: "Empresa A", CNPJ: "11111111000111", IsActive: true}, 1)
	if msg := exigeErroNoCampo(t, err, "parent_id"); msg != "uma empresa não pode ser vinculada a ela mesma" {
		t.Errorf("mensagem = %q", msg)
	}

	// A -> C fecha o ciclo A <- B <- C
	_, err = s.UpdateCompany(a, domain.UpdateCompanyRequest{ParentID: &c.ID, LegalName: "Empresa A", CNPJ: "11111111000111", IsActive: true}, 1)
	if msg := exigeErroNoCampo(t, err, "parent_id"); msg != "vínculo formaria um ciclo entre empresas" {
		t.Errorf("mensagem = %q", msg)
	}
}

func TestEmpresa_AlterarCNPJ(t *testing.T) {
	s := NewCompanyService(newTestUow(t))
	a := criarMatriz(t, s, "Empresa A", "11111111000111")
	criarMatriz(t, s, "Empresa B", "22222222000122")

	// o próprio CNPJ, mascarado, não é duplicidade
	if _, err := s.UpdateCompany(a, domain.UpdateCompanyRequest{LegalName: "Empresa A Nova", CNPJ: "11.111.111/0001-11", IsActive: true}, 2); err != nil {
		t.Fatalf("manter o próprio CNPJ: %v", err)
	}
	// o de outra empresa é
	_, err := s.UpdateCompany(a, domain.UpdateCompanyRequest{LegalName: "Empresa A", CNPJ: "22.222.222/0001-22", IsActive: true}, 2)
	exigeErroNoCampo(t, err, "cnpj")

	// matriz não pode ficar sem CNPJ ao editar
	_, err = s.UpdateCompany(a, domain.UpdateCompanyRequest{LegalName: "Empresa A", IsActive: true}, 2)
	exigeErroNoCampo(t, err, "cnpj")

	// quem alterou fica registrado
	detalhe, err := s.GetCompanyByID(a)
	if err != nil {
		t.Fatalf("buscar: %v", err)
	}
	if detalhe.LegalName != "Empresa A Nova" || detalhe.UpdatedBy == nil || detalhe.UpdatedBy.ID != 2 {
		t.Errorf("detalhe = %+v, quer nome novo e UpdatedBy 2", detalhe)
	}
	if detalhe.CreatedBy == nil || detalhe.CreatedBy.ID != 1 {
		t.Errorf("CreatedBy = %+v, quer 1", detalhe.CreatedBy)
	}

	_, err = s.UpdateCompany(9999, domain.UpdateCompanyRequest{LegalName: "Fantasma", CNPJ: "33333333000133", IsActive: true}, 2)
	exigeErroNoCampo(t, err, "id")
}

func TestEmpresa_Exclusao(t *testing.T) {
	s := NewCompanyService(newTestUow(t))
	matriz := criarMatriz(t, s, "Matriz", "11111111000111")
	filial, err := s.CreateCompany(domain.CreateCompanyRequest{ParentID: &matriz, LegalName: "Filial"}, 1)
	if err != nil {
		t.Fatalf("criar filial: %v", err)
	}

	if err := s.DeleteCompany(matriz); err == nil {
		t.Fatal("matriz com filial vinculada não pode ser excluída")
	} else if msg := exigeErroNoCampo(t, err, "id"); msg == "" {
		t.Error("mensagem vazia")
	}

	if err := s.DeleteCompany(filial.ID); err != nil {
		t.Fatalf("excluir filial: %v", err)
	}
	if err := s.DeleteCompany(matriz); err != nil {
		t.Fatalf("excluir matriz depois da filial: %v", err)
	}
	if _, err := s.GetCompanyByID(matriz); !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("empresa excluída: erro = %v, quer ErrNotFound", err)
	}
	if err := s.DeleteCompany(matriz); !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("excluir de novo: erro = %v, quer ErrNotFound", err)
	}
}

func TestEmpresa_Logo(t *testing.T) {
	s := NewCompanyService(newTestUow(t))
	id := criarMatriz(t, s, "Com Logo", "11111111000111")

	if _, _, err := s.GetLogo(id); !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("sem logo: erro = %v, quer ErrNotFound", err)
	}

	_, err := s.SetLogo(id, []byte("%PDF-1.4"), "application/pdf")
	exigeErroNoCampo(t, err, "logo")

	dto, err := s.SetLogo(id, []byte{0x89, 'P', 'N', 'G'}, "image/png")
	if err != nil {
		t.Fatalf("SetLogo: %v", err)
	}
	if !dto.HasLogo {
		t.Error("HasLogo deveria ser true depois do upload")
	}
	bytes, mime, err := s.GetLogo(id)
	if err != nil || mime != "image/png" || len(bytes) != 4 {
		t.Errorf("GetLogo = %v %q %v", bytes, mime, err)
	}

	dto, err = s.ClearLogo(id)
	if err != nil {
		t.Fatalf("ClearLogo: %v", err)
	}
	if dto.HasLogo {
		t.Error("HasLogo deveria ser false depois de remover")
	}
	if dto.LegalName != "Com Logo" {
		t.Error("remover o logo não pode mexer no resto do cadastro")
	}

	if _, err := s.SetLogo(9999, []byte("x"), "image/png"); !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("logo de empresa inexistente: erro = %v, quer ErrNotFound", err)
	}
}

// Visibilidade: master (sem empresa) vê tudo; usuário vinculado vê só a própria empresa, e com
// companies.hierarchy.view vê também as subsidiárias dela (mas nunca a matriz nem as irmãs).
func TestEmpresa_VisibilidadePorUsuario(t *testing.T) {
	uow := newTestUow(t)
	s := NewCompanyService(uow)
	db := uow.GetDB()

	grupo := criarMatriz(t, s, "Grupo", "11111111000111")
	filialA, _ := s.CreateCompany(domain.CreateCompanyRequest{ParentID: &grupo, LegalName: "Filial A"}, 1)
	netaA1, _ := s.CreateCompany(domain.CreateCompanyRequest{ParentID: &filialA.ID, LegalName: "Neta A1"}, 1)
	filialB, _ := s.CreateCompany(domain.CreateCompanyRequest{ParentID: &grupo, LegalName: "Filial B"}, 1)
	outroGrupo := criarMatriz(t, s, "Outro Grupo", "22222222000122")

	novoUsuario := func(companyID *uint) uint {
		row := struct{ ID uint }{}
		if err := db.Raw(`INSERT INTO users (company_id) VALUES (?) RETURNING id`, companyID).Scan(&row).Error; err != nil {
			t.Fatalf("usuário: %v", err)
		}
		return row.ID
	}
	darPermissao := func(userID uint, codigo string) {
		row := struct{ ID uint }{}
		if err := db.Raw(`INSERT INTO permissions (permission) VALUES (?) RETURNING id`, codigo).Scan(&row).Error; err != nil {
			t.Fatalf("permissão: %v", err)
		}
		if err := db.Exec(`INSERT INTO user_permissions (user_id, permission_id) VALUES (?, ?)`, userID, row.ID).Error; err != nil {
			t.Fatalf("vínculo: %v", err)
		}
	}
	ids := func(userID uint) map[uint]bool {
		lista, err := s.GetCompanies(userID)
		if err != nil {
			t.Fatalf("GetCompanies(%d): %v", userID, err)
		}
		out := map[uint]bool{}
		for _, c := range lista.Companies {
			out[c.ID] = true
		}
		return out
	}

	master := novoUsuario(nil)
	if got := ids(master); len(got) != 5 {
		t.Errorf("master vê %d empresas, quer 5", len(got))
	}

	semHierarquia := novoUsuario(&filialA.ID)
	if got := ids(semHierarquia); len(got) != 1 || !got[filialA.ID] {
		t.Errorf("vinculado sem hierarquia vê %v, quer só a Filial A", got)
	}

	comHierarquia := novoUsuario(&filialA.ID)
	darPermissao(comHierarquia, "companies.hierarchy.view")
	got := ids(comHierarquia)
	if len(got) != 2 || !got[filialA.ID] || !got[netaA1.ID] {
		t.Errorf("com hierarquia vê %v, quer a Filial A e a Neta A1", got)
	}
	for _, proibida := range []uint{grupo, filialB.ID, outroGrupo} {
		if got[proibida] {
			t.Errorf("com hierarquia não pode ver a empresa %d (matriz, irmã ou outro grupo)", proibida)
		}
	}
}

// CFG-050: o CNPJ de uma empresa excluída volta a ficar livre (índice único parcial).
func TestEmpresa_RecadastroDepoisDaExclusao(t *testing.T) {
	s := NewCompanyService(newTestUow(t))
	antiga := criarMatriz(t, s, "Matriz Antiga", "11111111000111")
	if err := s.DeleteCompany(antiga); err != nil {
		t.Fatalf("excluir: %v", err)
	}
	criarMatriz(t, s, "Matriz Nova", "11111111000111")
}
