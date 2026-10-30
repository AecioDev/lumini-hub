package testutil

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
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

// AssertRequestMaxMatchesColumnSize confere a regra de tamanho dos campos: todo campo de texto da
// entidade com gorm "size:N" que aparece numa request de criar/editar (mesmo nome no JSON) precisa ter
// binding "max=N". Sem o max=, o backend aceita o excesso e só o banco recusa.
//
// Use um por entidade: testutil.AssertRequestMaxMatchesColumnSize(t, Customer{}, CreateCustomerRequest{}, UpdateCustomerRequest{}).
func AssertRequestMaxMatchesColumnSize(t testing.TB, entity any, requests ...any) {
	t.Helper()

	entityType := reflect.TypeOf(entity)
	sizes := map[string]int{}
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if size, ok := columnSize(field); ok && jsonName(field) != "-" && jsonName(field) != "" {
			sizes[jsonName(field)] = size
		}
	}
	if len(sizes) == 0 {
		t.Fatalf("%s: nenhum campo com size no banco, o teste não está olhando nada", entityType.Name())
	}

	for _, request := range requests {
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
