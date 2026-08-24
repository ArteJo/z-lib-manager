package servicios_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArteJo/z-lib-manager/inventario"
	"github.com/ArteJo/z-lib-manager/operaciones"
	"github.com/ArteJo/z-lib-manager/servicios"
)

func prepararServidorPrueba() *servicios.ServidorAPI {
	catalogo := inventario.GenerarCatalogoInicial()
	repo := operaciones.NuevoRepositorioConcurrente(catalogo)
	return servicios.NuevoServidorAPI(repo)
}

func TestServicioObtenerLibros(t *testing.T) {
	srv := prepararServidorPrueba()

	req, _ := http.NewRequest("GET", "/api/v1/libros", nil)
	rr := httptest.NewRecorder()
	srv.Router().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Se esperaba código %d, pero se obtuvo %d", http.StatusOK, status)
	}

	var libros []map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &libros); err != nil {
		t.Fatalf("Respuesta no es un JSON válido: %v", err)
	}

	if len(libros) < 5 {
		t.Errorf("Se esperaban al menos 5 libros en el catálogo inicial, obtenidos %d", len(libros))
	}
}

func TestServicioCrearLibro(t *testing.T) {
	srv := prepararServidorPrueba()

	payload := []byte(`{
		"titulo": "Ciberseguridad Ofensiva",
		"autor": "Wilmer UIDE",
		"genero": "Seguridad",
		"anio": 2026,
		"popularidad": 5,
		"editorial": "UIDE Labs",
		"edicion": 1
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/libros", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Router().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Código esperado %d, obtenido %d", http.StatusCreated, status)
	}
}

func TestServicioValidacionError(t *testing.T) {
	srv := prepararServidorPrueba()

	// Envío de año inválido para provocar error de validación en encapsulación
	payload := []byte(`{
		"titulo": "Libro Erróneo",
		"autor": "Tester",
		"genero": "Test",
		"anio": 1200
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/libros", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Router().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Se esperaba Bad Request (%d) ante año inválido, obtenido %d", http.StatusBadRequest, status)
	}
}
