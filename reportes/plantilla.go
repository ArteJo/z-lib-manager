package reportes

import (
	"os"            // Permite interactuar con flujos de salida del sistema operativo
	"text/template" // Motor estándar de sustitución de variables independientes en texto

	"github.com/ArteJo/z-lib-manager/inventario" // Referencia local externa
)

// Plantilla para el reporte de libros electrónicos filtrados
const PlantillaConsola = `
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
Editorial:   {{.Editorial}} (Edición {{.Edicion}})
------------------------------------------------------
{{end}}======================================================
`

type DatosReporte struct {
	Criterio string
	Libros   []*inventario.Libro
}

func RenderizarConsola(libros []*inventario.Libro, tituloReporte string) error {
	t, err := template.New("rep").Parse(PlantillaConsola)
	if err != nil {
		return err
	}
	return t.Execute(os.Stdout, DatosReporte{Criterio: tituloReporte, Libros: libros})
}
