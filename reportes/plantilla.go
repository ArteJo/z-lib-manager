package reportes

import (
	"os"            // Permite interactuar con flujos de salida del sistema operativo
	"text/template" // Motor estándar de sustitución de variables independientes en texto

	"github.com/ArteJo/z-lib-manager/inventario" // Referencia local externa
)

// Renderizador abstrae la generación visual del reporte
type Renderizador interface {
	Renderizar(libros []*inventario.Libro, criterioNombre string) error
}

type ReporteConsola struct{}

// Plantilla para el reporte de libros electrónicos filtrados
const plantillaFormato = `
REPORTE DE LIBROS: {{.Criterio}}
======================================================
Total encontrados: {{len .Libros}}
------------------------------------------------------
{{range .Libros}}ID:          {{.ID}}
Título:      {{.Titulo}}
Autor:       {{.Autor}}
Género:      {{.Genero}}
Año:         {{.Anio}}
Popularidad: {{.Popularidad}} / 5 estrellas
------------------------------------------------------
{{end}}======================================================
`

type DatosReporte struct {
	Criterio string
	Libros   []*inventario.Libro
}

func (r ReporteConsola) Renderizar(libros []*inventario.Libro, criterioNombre string) error {
	tmpl, err := template.New("reporte").Parse(plantillaFormato)
	if err != nil {
		return err
	}

	datos := DatosReporte{
		Criterio: criterioNombre,
		Libros:   libros,
	}

	return tmpl.Execute(os.Stdout, datos)
}
