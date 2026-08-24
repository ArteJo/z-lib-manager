package operaciones

import (
	"fmt"
	"strings"

	"github.com/ArteJo/z-lib-manager/inventario"
)

// CriterioBusqueda define el contrato polimórfico
type CriterioBusqueda interface {
	Cumple(l *inventario.Libro) bool
	Descripcion() string
}

// FiltroGenero
type FiltroGenero struct {
	GeneroBuscado string
}

func (f FiltroGenero) Cumple(l *inventario.Libro) bool {
	return strings.EqualFold(l.Genero(), f.GeneroBuscado)
}

func (f FiltroGenero) Descripcion() string {
	return fmt.Sprintf("Género = '%s'", f.GeneroBuscado)
}

// FiltroAnioMinimo
type FiltroAnioMinimo struct {
	AnioMinimo int
}

func (f FiltroAnioMinimo) Cumple(l *inventario.Libro) bool {
	return l.Anio() >= f.AnioMinimo
}

func (f FiltroAnioMinimo) Descripcion() string {
	return fmt.Sprintf("Año >= %d", f.AnioMinimo)
}

// FiltroPopularidad
type FiltroPopularidad struct {
	PopularidadMinima int
}

func (f FiltroPopularidad) Cumple(l *inventario.Libro) bool {
	return l.Popularidad() >= f.PopularidadMinima
}

func (f FiltroPopularidad) Descripcion() string {
	return fmt.Sprintf("Popularidad >= %d estrellas", f.PopularidadMinima)
}

// FiltroTextoGeneral
type FiltroTextoGeneral struct {
	Termino string
}

func (f FiltroTextoGeneral) Cumple(l *inventario.Libro) bool {
	t := strings.ToLower(f.Termino)
	return strings.Contains(strings.ToLower(l.Titulo()), t) || strings.Contains(strings.ToLower(l.Autor()), t)
}

func (f FiltroTextoGeneral) Descripcion() string {
	return fmt.Sprintf("Texto libre = '%s'", f.Termino)
}

// Filtrar procesa polimórficamente cualquier estructura que satisfaga CriterioBusqueda
func Filtrar(catalogo []*inventario.Libro, criterio CriterioBusqueda) []*inventario.Libro {
	var res []*inventario.Libro
	for _, libro := range catalogo {
		if criterio.Cumple(libro) {
			res = append(res, libro)
		}
	}
	return res
}
