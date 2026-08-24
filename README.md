# z-lib-manager: Sistema de Gestión de Libros Electrónicos

**Carrera:** Ingeniería en Ciberseguridad

**Institución:** Universidad Internacional del Ecuador (UIDE)

**Asignatura:** Programación Orientada a Objetos en Go

**Fecha:** 23/08/2026

## Descripción Conceptual
`z-lib-manager` es una plataforma "backend" que gestiona un catálogo digital de libros electrónicos aplicando los conceptos fundamentales de las 4 unidades del curso que cursamos: Programación Funcional, POO en Go (Structs, Embedding, Constructores), Encapsulación, Interfaces Implícitas, Manejo de Errores, Concurrencia (Goroutines/Channels) y Servicios Web RESTful con serialización en JSON.

### Prerrequisitos
* Go 1.20 o superior instalado.
* Git.

### Ejecutar Servidor y Menú Interactivo
`go run main.go`

El servidor HTTP iniciará en http://localhost:8080 de manera concurrente con la consola interactiva.

El menú interactivo es el siguiente:

1. Ver Catálogo de Libros Completos (Orden A-Z & Año)
2. Filtrar Libros por Género
3. Filtrar Libros por Año de Publicación
4. Filtrar Libros por su Género y cuantos libros hay en cada genero (Tabla Hash)
5. Comprueba una validación de datos y manejo de errores (Encapsulación)
6. Búsqueda Rápida en Paralelo (Multihilo / Goroutines)
7. Ver Lista de Endpoints de la API
8. Analizar y Contar Letras de un Título (Librería Externa)
9. Salir

