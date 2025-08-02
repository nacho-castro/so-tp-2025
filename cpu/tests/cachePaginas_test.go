package tests

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/traducciones"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuscarPagina(t *testing.T) {
	err := CargarConfigCPU("../configs/PLANI.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}
	traducciones.Max = 4

	cache := traducciones.NuevaCachePaginas()
	cache.Agregar(1, "Test", true)

	contenido, ok := cache.Buscar(1)

	assert.True(t, ok)
	assert.Equal(t, "Test", contenido)
}

func TestActivacionCache(t *testing.T) {
	err := CargarConfigCPU("../configs/config1.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}
	traducciones.Max = 4

	cache := traducciones.NuevaCachePaginas()

	bool := cache.EstaActiva()

	assert.True(t, bool)
}

func TestAgregarEntrada(t *testing.T) {
	err := CargarConfigCPU("../configs/PLANI.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}
	traducciones.Max = 4

	cache := traducciones.NuevaCachePaginas()

	cache.Agregar(1, "Test", true)

	contenido, ok := cache.Buscar(1)

	assert.True(t, ok)
	assert.Equal(t, "Test", contenido)
}

func TestReemplazoClock(t *testing.T) {
	err := CargarConfigCPU("../configs/CP.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}
	traducciones.Max = 4

	cache := traducciones.NuevaCachePaginas()
	cache.Agregar(1, "", false)
	cache.Agregar(2, "Reemplazada", false)
	cache.Agregar(3, "", false)
	cache.Agregar(4, "", false)

	cache.MarcarUso(1)

	cache.Agregar(5, "", false)

	contenido, ok := cache.Buscar(2) //Reemplazada por (0,0)

	assert.False(t, ok)
	assert.Equal(t, "", contenido)
}

func TestReemplazoClockM(t *testing.T) {
	err := CargarConfigCPU("../configs/CP.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}
	traducciones.Max = 4

	cache := traducciones.NuevaCachePaginas()
	cache.Entradas = []traducciones.EntradaCache{
		{NroPagina: 1, Contenido: "A", Usado: true, Modificado: true},
		{NroPagina: 2, Contenido: "B", Usado: false, Modificado: false},
		{NroPagina: 3, Contenido: "C", Usado: true, Modificado: true},
		{NroPagina: 4, Contenido: "D", Usado: true, Modificado: true},
	}

	cache.Agregar(5, "", false)

	contenido, ok := cache.Buscar(2) //Reemplazada por (0,0)

	assert.False(t, ok)
	assert.Equal(t, "", contenido)
}

func TestReemplazoClockM2(t *testing.T) {
	err := CargarConfigCPU("../configs/CP.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}
	traducciones.Max = 4

	cache := traducciones.NuevaCachePaginas()
	cache.Entradas = []traducciones.EntradaCache{
		{NroPagina: 1, Contenido: "A", Usado: true, Modificado: true},
		{NroPagina: 2, Contenido: "B", Usado: false, Modificado: true},
		{NroPagina: 3, Contenido: "C", Usado: true, Modificado: true},
		{NroPagina: 4, Contenido: "D", Usado: true, Modificado: true},
	}

	cache.Agregar(5, "", false)

	contenido, ok := cache.Buscar(2) //Reemplazada por (0,1)

	assert.False(t, ok)
	assert.Equal(t, "", contenido)
}

func TestReemplazoClockM3(t *testing.T) {
	err := CargarConfigCPU("../configs/CP.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}
	traducciones.Max = 4

	cache := traducciones.NuevaCachePaginas()
	cache.Entradas = []traducciones.EntradaCache{
		{NroPagina: 1, Contenido: "A", Usado: true, Modificado: true},
		{NroPagina: 2, Contenido: "B", Usado: true, Modificado: true},
		{NroPagina: 3, Contenido: "C", Usado: true, Modificado: true},
		{NroPagina: 4, Contenido: "D", Usado: true, Modificado: true},
	}

	cache.Agregar(5, "", false)

	contenido, ok := cache.Buscar(1) //Reemplazada por (0,0) desp de marcar U=0

	assert.False(t, ok)
	assert.Equal(t, "", contenido)
}

func TestLimpiarCache(t *testing.T) {
	err := CargarConfigCPU("../configs/CP.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}
	traducciones.Max = 4
	globals.TamanioPagina = 2

	cache := traducciones.NuevaCachePaginas()
	cache.Entradas = []traducciones.EntradaCache{
		{NroPagina: 1, Contenido: "A", Usado: true, Modificado: true},
		{NroPagina: 2, Contenido: "B", Usado: true, Modificado: true},
		{NroPagina: 3, Contenido: "C", Usado: true, Modificado: true},
		{NroPagina: 4, Contenido: "D", Usado: true, Modificado: true},
	}

	cache.LimpiarCache()

	assert.Equal(t, 0, len(cache.Entradas))
}
