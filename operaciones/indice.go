package operaciones

import "github.com/ArteJo/z-lib-manager/inventario"

// GenerarIndicePorGenero indexa el catálogo en una Tabla Hash para búsquedas en O(1)
func GenerarIndicePorGenero(catalogo []*inventario.Libro) map[string][]*inventario.Libro {
	indice := make(map[string][]*inventario.Libro)
	for _, l := range catalogo {
		indice[l.Genero()] = append(indice[l.Genero()], l)
	}
	return indice
}
