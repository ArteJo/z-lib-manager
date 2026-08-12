package main

import (
	"fmt" // Paquete estándar para interactuar con entradas del teclado
	"log" // Paquete estándar para manejo de errores y logging

	// Rutas explícitas locales
	"github.com/ArteJo/z-lib-manager/analizador"
	"github.com/ArteJo/z-lib-manager/inventario"
	"github.com/ArteJo/z-lib-manager/operaciones"
	"github.com/ArteJo/z-lib-manager/reportes"
)

func main() {
	catalogo := inventario.GenerarCatalogoInicial()
	reportador := reportes.ReporteConsola{}
	auditor := analizador.AuditorResumen{}

	var opcion int
	fmt.Println("SISTEMA DE GESTIÓN Z-LIB-MANAGER")
	fmt.Println("1. Filtrar por Género 'Tecnología'")
	fmt.Println("2. Filtrar por Año Mínimo (>= 2024)")
	fmt.Println("3. Filtrar por Popularidad (>= 4 estrellas)")
	fmt.Println("4. Búsqueda directa en Tabla Hash (Map por Género)")
	fmt.Println("5. Probar Validación de Errores al crear un Libro")
	fmt.Print("Seleccione una opción: ")

	if _, err := fmt.Scan(&opcion); err != nil {
		log.Fatalf("Error al leer la entrada: %v", err)
	}

	switch opcion {
	case 1:
		filtro := operaciones.FiltroGenero{GeneroBuscado: "Tecnología"}
		resultados := operaciones.Filtrar(catalogo, filtro)
		_ = reportador.Renderizar(resultados, filtro.Descripcion())

	case 2:
		filtro := operaciones.FiltroAnioMinimo{AnioMinimo: 2024}
		resultados := operaciones.Filtrar(catalogo, filtro)
		_ = reportador.Renderizar(resultados, filtro.Descripcion())

	case 3:
		filtro := operaciones.FiltroPopularidad{PopularidadMinima: 4}
		resultados := operaciones.Filtrar(catalogo, filtro)
		_ = reportador.Renderizar(resultados, filtro.Descripcion())

	case 4:
		fmt.Println("\n--- BÚSQUEDA RÁPIDA VÍA TABLA HASH (MAPS) ---")
		indiceMap := operaciones.GenerarIndicePorGenero(catalogo)
		librosSeguridad, existe := indiceMap["Seguridad"]
		if existe {
			_ = reportador.Renderizar(librosSeguridad, "Índice Map: Categoría 'Seguridad'")
		}

	case 5:
		fmt.Println("\n--- PRUEBA DE MANEJO DE ERRORES Y ENCAPSULACIÓN ---")
		_, err := inventario.NewLibro(99, "Libro Inválido", "Autor Test", "Pruebas", 1200)
		if err != nil {
			fmt.Printf(" [Error Capturado Exitosamente]: %v\n", err)
		}

	default:
		fmt.Println("Opción no válida. Mostrando catálogo completo.")
		_ = reportador.Renderizar(catalogo, "Catálogo Completo")
	}

	if len(catalogo) > 0 {
		fmt.Println("\n--- AUDITORÍA DE METADATOS DE TEXTO ---")
		if err := auditor.Auditar(catalogo[0].Titulo()); err != nil {
			fmt.Printf("Error durante la auditoría: %v\n", err)
		}
	}
}
