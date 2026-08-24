package analizador

import (
	"bytes"  // Paquete estándar para un uso eficiente de búferes de bytes
	"errors" // Paquete estándar para manejo de errores
	"fmt"    // Paquete estándar de formato de entrada/salida

	"github.com/mariomac/analizador" // Importamos una dependencia externa de terceros (un contador de caracteres y palabras)
)

// AuditorTexto define el contrato para la revisión de metadatos y resúmenes
type AuditorTexto interface {
	Auditar(texto string) (string, error)
}

// AuditorResumen implementa AuditorTexto mediante búferes optimizados
type AuditorResumen struct{}

// Auditar procesa el título o texto utilizando bytes.Buffer y el paquete externo
func (a AuditorResumen) Auditar(texto string) (resumen string, err error) {
	// Captura defensiva de pánicos en tiempo de ejecución con defer y recover
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error crítico en el análisis de texto: %v", r)
		}
	}()

	if texto == "" {
		return "", errors.New("el texto a auditar no puede estar vacío")
	}

	var buf bytes.Buffer
	buf.WriteString("\n[Auditoría de Texto]:\n")
	fmt.Fprintf(&buf, "Cadena evaluada: \"%s\"\n", texto)

	// Llamada a la función del paquete externo de GitHub
	analizador.PrintEstadistica(texto)
	return buf.String(), nil
}

