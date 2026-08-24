package servicios

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ArteJo/z-lib-manager/analizador"
	"github.com/ArteJo/z-lib-manager/inventario"
	"github.com/ArteJo/z-lib-manager/operaciones"
)

type ServidorAPI struct {
	repo    *operaciones.RepositorioConcurrente
	auditor analizador.AuditorTexto
	mux     *http.ServeMux
}

func NuevoServidorAPI(repo *operaciones.RepositorioConcurrente) *ServidorAPI {
	s := &ServidorAPI{
		repo:    repo,
		auditor: analizador.AuditorResumen{},
		mux:     http.NewServeMux(),
	}
	s.configurarRutas()
	return s
}

func (s *ServidorAPI) Router() http.Handler {
	return s.mux
}

func (s *ServidorAPI) configurarRutas() {
	s.mux.HandleFunc("/api/v1/libros", s.manejadorColeccionLibros)                  // Servicio 1 (GET) y Servicio 3 (POST)
	s.mux.HandleFunc("/api/v1/libros/", s.manejadorElementoLibro)                   // Servicio 2 (GET), Servicio 4 (PUT), Servicio 5 (DELETE)
	s.mux.HandleFunc("/api/v1/libros/filtrar", s.manejadorFiltrado)                 // Servicio 6 (GET Filtrado Polimórfico)
	s.mux.HandleFunc("/api/v1/libros/indice/genero", s.manejadorIndice)             // Servicio 7 (GET Map Hash O(1))
	s.mux.HandleFunc("/api/v1/libros/analizar", s.manejadorAnalizarTexto)           // Servicio 8 (POST Auditoría de Texto)
	s.mux.HandleFunc("/api/v1/concurrencia/buscar", s.manejadorBusquedaConcurrente) // Servicio 9 (GET Paralelismo/Canales)
}

func responderJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func responderError(w http.ResponseWriter, statusCode int, mensaje string) {
	responderJSON(w, statusCode, map[string]string{"error": mensaje})
}

// SERVICIO 1 (GET /api/v1/libros) & SERVICIO 3 (POST /api/v1/libros)
func (s *ServidorAPI) manejadorColeccionLibros(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		libros := s.repo.ObtenerTodos()
		responderJSON(w, http.StatusOK, libros)

	case http.MethodPost:
		var nuevoLibro inventario.Libro
		if err := json.NewDecoder(r.Body).Decode(&nuevoLibro); err != nil {
			responderError(w, http.StatusBadRequest, fmt.Sprintf("JSON inválido: %v", err))
			return
		}
		guardado := s.repo.Guardar(&nuevoLibro)
		responderJSON(w, http.StatusCreated, guardado)

	default:
		responderError(w, http.StatusMethodNotAllowed, "Método HTTP no permitido")
	}
}

// SERVICIO 2 (GET), SERVICIO 4 (PUT), SERVICIO 5 (DELETE) en /api/v1/libros/{id}
func (s *ServidorAPI) manejadorElementoLibro(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(partes) < 4 {
		responderError(w, http.StatusBadRequest, "ID de libro no proporcionado")
		return
	}

	id, err := strconv.Atoi(partes[3])
	if err != nil {
		responderError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	switch r.Method {
	case http.MethodGet:
		libro, existe := s.repo.ObtenerPorID(id)
		if !existe {
			responderError(w, http.StatusNotFound, "Libro no encontrado")
			return
		}
		responderJSON(w, http.StatusOK, libro)

	case http.MethodPut:
		var libroActualizado inventario.Libro
		if err := json.NewDecoder(r.Body).Decode(&libroActualizado); err != nil {
			responderError(w, http.StatusBadRequest, fmt.Sprintf("Datos JSON inválidos: %v", err))
			return
		}
		_, existe := s.repo.ObtenerPorID(id)
		if !existe {
			responderError(w, http.StatusNotFound, "Libro no encontrado para actualizar")
			return
		}
		// Forzar el ID correspondiente
		l, _ := inventario.NewLibro(id, libroActualizado.Titulo(), libroActualizado.Autor(),
			libroActualizado.Genero(), libroActualizado.Anio(),
			inventario.ConPopularidad(libroActualizado.Popularidad()),
			inventario.ConEditorial(libroActualizado.Editorial, libroActualizado.Edicion))

		guardado := s.repo.Guardar(l)
		responderJSON(w, http.StatusOK, guardado)

	case http.MethodDelete:
		eliminado := s.repo.Eliminar(id)
		if !eliminado {
			responderError(w, http.StatusNotFound, "Libro no encontrado para eliminación")
			return
		}
		responderJSON(w, http.StatusOK, map[string]string{"mensaje": "Libro eliminado correctamente"})

	default:
		responderError(w, http.StatusMethodNotAllowed, "Método HTTP no permitido")
	}
}

// SERVICIO 6: GET /api/v1/libros/filtrar?genero=...&anioMin=...&popMin=...
func (s *ServidorAPI) manejadorFiltrado(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	libros := s.repo.ObtenerTodos()
	q := r.URL.Query()

	if genero := q.Get("genero"); genero != "" {
		libros = operaciones.Filtrar(libros, operaciones.FiltroGenero{GeneroBuscado: genero})
	}
	if anioStr := q.Get("anioMin"); anioStr != "" {
		if anio, err := strconv.Atoi(anioStr); err == nil {
			libros = operaciones.Filtrar(libros, operaciones.FiltroAnioMinimo{AnioMinimo: anio})
		}
	}
	if popStr := q.Get("popMin"); popStr != "" {
		if pop, err := strconv.Atoi(popStr); err == nil {
			libros = operaciones.Filtrar(libros, operaciones.FiltroPopularidad{PopularidadMinima: pop})
		}
	}

	responderJSON(w, http.StatusOK, libros)
}

// SERVICIO 7: GET /api/v1/libros/indice/genero
func (s *ServidorAPI) manejadorIndice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}
	indice := operaciones.GenerarIndicePorGenero(s.repo.ObtenerTodos())
	responderJSON(w, http.StatusOK, indice)
}

// SERVICIO 8: POST /api/v1/libros/analizar
func (s *ServidorAPI) manejadorAnalizarTexto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	var req struct {
		Texto string `json:"texto"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Texto == "" {
		responderError(w, http.StatusBadRequest, "Cuerpo JSON inválido o 'texto' vacío")
		return
	}

	resumen, err := s.auditor.Auditar(req.Texto)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{
		"estado":  "Auditado con éxito",
		"resumen": resumen,
	})
}

// SERVICIO 9 (Extra Concurrencia): GET /api/v1/concurrencia/buscar?q=...
func (s *ServidorAPI) manejadorBusquedaConcurrente(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	termino := r.URL.Query().Get("q")
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	resultados, err := s.repo.BusquedaParalela(ctx, termino)
	if err != nil {
		responderError(w, http.StatusGatewayTimeout, "Tiempo de búsqueda excedido")
		return
	}

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"termino":   termino,
		"total":     len(resultados),
		"resultado": resultados,
	})
}
