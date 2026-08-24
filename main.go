package main

import (
	"context"  // Paquete estándar para manejo de contexto y cancelación
	"fmt"      // Paquete estándar para interactuar con entradas del teclado
	"log"      // Paquete estándar para manejo de errores y logging
	"net/http" // Paquete estándar para manejo de solicitudes HTTP
	"sort"     // Paquete estándar para ordenamiento de slices y arrays
	"strings"  // Paquete estándar para manipulación de cadenas de texto
	"time"     // Paquete estándar para manejo de tiempo y fechas

	// Rutas explícitas locales

	"github.com/ArteJo/z-lib-manager/analizador"
	"github.com/ArteJo/z-lib-manager/inventario"
	"github.com/ArteJo/z-lib-manager/operaciones"
	"github.com/ArteJo/z-lib-manager/reportes"
	"github.com/ArteJo/z-lib-manager/servicios"
)

func main() {
	catalogo := inventario.GenerarCatalogoInicial()
	repo := operaciones.NuevoRepositorioConcurrente(catalogo)
	srv := servicios.NuevoServidorAPI(repo)
	auditor := analizador.AuditorResumen{}

	// Iniciar el servidor Web en una Goroutine concurrente en segundo plano
	go func() {
		servidorHTTP := &http.Server{
			Addr:         ":8080",
			Handler:      srv.Router(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		fmt.Println("[Servidor Activo] Escuchando en http://localhost:8080")
		if err := servidorHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error en el servidor HTTP: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond) // Pequeña pausa para permitir que el mensaje del servidor no choque con el menú

	// Menú CLI interactivo en la Goroutine principal
	for {
		fmt.Println("Z-LIB-MANAGER: SISTEMA DE GESTIÓN DE LIBROS ELECTRÓNICOS")
		fmt.Println("=======================================================")
		fmt.Println("1. Ver Catálogo de Libros Completos (Orden A-Z & Año)")
		fmt.Println("2. Filtrar Libros por Género (Polimórfismo)")
		fmt.Println("3. Filtrar Libros por Año de Publicación (Polimorfismo)")
		fmt.Println("4. Filtrar Libros por su Género y cuantos libros hay en cada genero (Tabla Hash)")
		fmt.Println("5. Comprueba una validación de datos y manejo de errores (Encapsulación)")
		fmt.Println("6. Búsqueda Rápida en Paralelo (Multihilo / Goroutines)")
		fmt.Println("7. Ver Lista de Endpoints de la API") // Muestra la lista URL disponibles que el server HTTP esta escuchando
		fmt.Println("8. Analizar y Contar Letras de un Título (Librería Externa)")
		fmt.Println("9. Salir")
		fmt.Print("Seleccione una opción (1-9): ")

		var op int
		if _, err := fmt.Scan(&op); err != nil {
			fmt.Println("Por favor ingrese un número válido.")
			var descartar string
			fmt.Scanln(&descartar)
			continue
		}

		switch op {
		case 1:
			// OPCIÓN 1: Ver catálogo con submenú de ordenamiento
			fmt.Println("\nVISTA DE CATÁLOGO")
			fmt.Println("a. Orden original")
			fmt.Println("b. Ordenar alfabéticamente por Título")
			fmt.Println("c. Ordenar cronológicamente por Año (Más reciente primero)")
			fmt.Print("Seleccione el tipo de vista: ")
			var subOp string
			fmt.Scan(&subOp)

			lista := repo.ObtenerTodos()
			copiaLista := make([]*inventario.Libro, len(lista))
			copy(copiaLista, lista)

			switch strings.ToLower(subOp) {
			case "b":
				sort.Slice(copiaLista, func(i, j int) bool {
					return copiaLista[i].Titulo() < copiaLista[j].Titulo()
				})
				_ = reportes.RenderizarConsola(copiaLista, "Catálogo: Orden Alfabético por Título")
			case "c":
				sort.Slice(copiaLista, func(i, j int) bool {
					return copiaLista[i].Anio() > copiaLista[j].Anio()
				})
				_ = reportes.RenderizarConsola(copiaLista, "Catálogo: Orden Cronológico (Más Recientes)")
			default:
				_ = reportes.RenderizarConsola(copiaLista, "Catálogo Completo")
			}

		case 2:
			// OPCIÓN 2: Filtrado por géneros disponibles detectados automáticamente
			indice := operaciones.GenerarIndicePorGenero(repo.ObtenerTodos())
			var generosDisponibles []string
			for g := range indice {
				generosDisponibles = append(generosDisponibles, g)
			}
			sort.Strings(generosDisponibles)

			fmt.Println("\n--- GÉNEROS DISPONIBLES EN EL SISTEMA ---")
			for i, g := range generosDisponibles {
				fmt.Printf("%d. %s (%d libros)\n", i+1, g, len(indice[g]))
			}
			fmt.Print("Seleccione el número de género a filtrar: ")
			var numGen int
			if _, err := fmt.Scan(&numGen); err == nil && numGen >= 1 && numGen <= len(generosDisponibles) {
				seleccionado := generosDisponibles[numGen-1]
				filtro := operaciones.FiltroGenero{GeneroBuscado: seleccionado}
				res := operaciones.Filtrar(repo.ObtenerTodos(), filtro)
				_ = reportes.RenderizarConsola(res, filtro.Descripcion())
			} else {
				fmt.Println("Selección no válida.")
			}

		case 3:
			// OPCIÓN 3: Filtro por año ingresado por el usuario
			fmt.Println("\nFILTRAR POR AÑO DE PUBLICACIÓN")
			fmt.Print("Ingrese el año mínimo a consultar (ej. 2010): ")
			var anioIngresado int
			if _, err := fmt.Scan(&anioIngresado); err == nil {
				filtro := operaciones.FiltroAnioMinimo{AnioMinimo: anioIngresado}
				res := operaciones.Filtrar(repo.ObtenerTodos(), filtro)
				_ = reportes.RenderizarConsola(res, filtro.Descripcion())
			} else {
				fmt.Println("Año ingresado no válido.")
			}

		case 4:
			// OPCIÓN 4: Explorar categorías vía Hash Map e inspeccionar
			indice := operaciones.GenerarIndicePorGenero(repo.ObtenerTodos())
			fmt.Println("\n--- RESUMEN EN TIEMPO CONSTANTE O(1) (TABLA HASH) ---")
			var categorias []string
			for cat := range indice {
				categorias = append(categorias, cat)
			}
			sort.Strings(categorias)

			for i, cat := range categorias {
				fmt.Printf("[%d] %s -> %d libros\n", i+1, cat, len(indice[cat]))
			}

			fmt.Print("\n¿Desea ver los libros de alguna categoría? Ingrese el número (o 0 para volver): ")
			var opcionCat int
			if _, err := fmt.Scan(&opcionCat); err == nil && opcionCat >= 1 && opcionCat <= len(categorias) {
				catElegida := categorias[opcionCat-1]
				_ = reportes.RenderizarConsola(indice[catElegida], fmt.Sprintf("Categoría: %s", catElegida))
			}

		case 5:
			// OPCIÓN 5: Comprobación clara de validaciones y errores
			fmt.Println("\nPRUEBA DE SEGURIDAD: VALIDACIÓN Y ERRORES")
			fmt.Println("Intentando registrar un libro con año inválido (Año: 1200)...")
			_, errAnio := inventario.NewLibro(99, "Libro Antiguo Inválido", "Autor Test", "Historia", 1200)
			if errAnio != nil {
				fmt.Printf(" [Validación Exitosa]: %v\n", errAnio)
			}

			fmt.Println("\nIntentando asignar una calificación fuera de rango (Popularidad: 10 estrellas)...")
			libroPrueba, _ := inventario.NewLibro(100, "Libro Test", "Autor", "General", 2023)
			errPop := libroPrueba.SetPopularidad(10)
			if errPop != nil {
				fmt.Printf(" [Validación Exitosa]: %v\n", errPop)
			}

		case 6:
			// OPCIÓN 6: Búsqueda paralela con Goroutines (Solucionado el Panic)
			fmt.Println("\nBÚSQUEDA RÁPIDA EN PARALELO (MULTIHILO)")
			fmt.Print("Ingrese palabra clave (Título, Autor o Género): ")
			var termino string
			fmt.Scan(&termino)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			inicio := time.Now()
			res, err := repo.BusquedaParalela(ctx, termino)
			duracion := time.Since(inicio)

			if err != nil {
				fmt.Printf("Error durante la búsqueda: %v\n", err)
			} else {
				_ = reportes.RenderizarConsola(res, fmt.Sprintf("Resultados en Paralelo para '%s' (Tiempo: %v)", termino, duracion))
			}

		case 7:
			// OPCIÓN 7: Lista de Endpoints REST
			fmt.Println("\n=======================================================")
			fmt.Println("             LISTA DE ENDPOINTS DE LA API REST        ")
			fmt.Println("=======================================================")
			fmt.Println("Servidor local activo en: http://localhost:8080")
			fmt.Println("1. GET    /api/v1/libros                   -> Obtener todos los libros")
			fmt.Println("2. GET    /api/v1/libros/{id}              -> Obtener libro por ID (ej. /1)")
			fmt.Println("3. POST   /api/v1/libros                   -> Crear nuevo libro (JSON)")
			fmt.Println("4. PUT    /api/v1/libros/{id}              -> Actualizar libro (JSON)")
			fmt.Println("5. DELETE /api/v1/libros/{id}              -> Eliminar libro por ID")
			fmt.Println("6. GET    /api/v1/libros/filtrar?genero=.. -> Filtrado polimórfico")
			fmt.Println("7. GET    /api/v1/libros/indice/genero     -> Índice agrupado en O(1)")
			fmt.Println("8. POST   /api/v1/libros/analizar          -> Auditoría de texto (JSON)")
			fmt.Println("9. GET    /api/v1/concurrencia/buscar?q=.. -> Búsqueda concurrente")
			fmt.Println("=======================================================")

		case 8:
			// OPCIÓN 8: Auditoría y conteo de letras con la librería de GitHub
			fmt.Println("\nAUDITAR Y CONTAR LETRAS DE UN TÍTULO")
			todos := repo.ObtenerTodos()
			for i, l := range todos {
				fmt.Printf("%2d. %s\n", i+1, l.Titulo())
			}
			fmt.Print("Seleccione el número de libro a auditar: ")
			var numLibro int
			if _, err := fmt.Scan(&numLibro); err == nil && numLibro >= 1 && numLibro <= len(todos) {
				tituloElegido := todos[numLibro-1].Titulo()
				_, err := auditor.Auditar(tituloElegido)
				if err != nil {
					fmt.Printf("Error en la auditoría: %v\n", err)
				}
			} else {
				fmt.Println("Selección no válida.")
			}

		case 9:
			fmt.Println("Cerrando z-lib-manager. ¡Hasta pronto!")
			return

		default:
			fmt.Println("Opción no válida. Por favor elija un número entre 1 y 9.")
		}
	}
}

