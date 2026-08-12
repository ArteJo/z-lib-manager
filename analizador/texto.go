package analizador

import (
	"bytes"  // Paquete estándar para un uso eficiente de búferes de bytes
	"errors" // Paquete estándar para manejo de errores
	"fmt"    // Paquete estándar de formato de entrada/salida

	"github.com/mariomac/analizador" // Importamos una dependencia externa de terceros (un contador de caracteres y palabras)
)

// AuditorTexto define el contrato para la revisión de metadatos y resúmenes
type AuditorTexto interface {
	Auditar(texto string) error
}

// AuditorResumen implementa AuditorTexto mediante búferes optimizados
type AuditorResumen struct{}

func (a AuditorResumen) Auditar(texto string) (err error) {
	// Manejo excepcional con defer y recover para evitar caídas catastróficas
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recuperado de un error crítico de procesamiento: %v", r)
		}
	}()

	if texto == "" {
		return errors.New("el texto a auditar no puede estar vacío")
	}

	var buf bytes.Buffer
	buf.WriteString("Auditando metadatos de texto: ")
	fmt.Fprintf(&buf, "\"%s\"", texto)
	fmt.Println(buf.String())

	// Paquete externo de terceros
	analizador.PrintEstadistica(texto)
	return nil
}
