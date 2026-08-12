package operaciones

import (
	"fmt"

	"github.com/ArteJo/z-lib-manager/inventario"
)

// CriterioBusqueda es una Interfaz que define el contrato para el filtrado polimórfico
type CriterioBusqueda interface {
	Cumple(l *inventario.Libro) bool
	Descripcion() string
}

// Implementación 1: Filtro por Género
type FiltroGenero struct {
	GeneroBuscado string
}

func (f FiltroGenero) Cumple(l *inventario.Libro) bool {
	return l.Genero() == f.GeneroBuscado
}

func (f FiltroGenero) Descripcion() string {
	return "Filtrado por Género: '" + f.GeneroBuscado + "'"
}

// Implementación 2: Filtro por Año Mínimo
type FiltroAnioMinimo struct {
	AnioMinimo int
}

func (f FiltroAnioMinimo) Cumple(l *inventario.Libro) bool {
	return l.Anio() >= f.AnioMinimo
}

func (f FiltroAnioMinimo) Descripcion() string {
	return fmt.Sprintf("Filtrado por Año de Publicación >= %d", f.AnioMinimo)
}

// Implementación 3: Filtro por Popularidad
type FiltroPopularidad struct {
	PopularidadMinima int
}

func (f FiltroPopularidad) Cumple(l *inventario.Libro) bool {
	return l.Popularidad() >= f.PopularidadMinima
}

func (f FiltroPopularidad) Descripcion() string {
	return fmt.Sprintf("Filtrado por Popularidad >= %d estrellas", f.PopularidadMinima)
}

// Filtrar procesa polimórficamente cualquier estructura que satisfaga CriterioBusqueda
func Filtrar(catalogo []*inventario.Libro, criterio CriterioBusqueda) []*inventario.Libro {
	var resultado []*inventario.Libro
	for _, libro := range catalogo {
		if criterio.Cumple(libro) {
			resultado = append(resultado, libro)
		}
	}
	return resultado
}

// GenerarIndicePorGenero utiliza una Tabla Hash (map) para indexación en O(1)
func GenerarIndicePorGenero(catalogo []*inventario.Libro) map[string][]*inventario.Libro {
	indice := make(map[string][]*inventario.Libro)
	for _, libro := range catalogo {
		genero := libro.Genero()
		indice[genero] = append(indice[genero], libro)
	}
	return indice
}
