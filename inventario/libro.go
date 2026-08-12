package inventario

import (
	"errors"
	"fmt"
)

// Libro implementa encapsulación estricta con campos privados (en minúscula)
// inaccesibles desde fuera del paquete inventario.
type Libro struct {
	id          int
	titulo      string
	autor       string
	genero      string
	anio        int
	popularidad int // Calificación de 1 a 5 estrellas
}

// Getters (Métodos de lectura)
func (l *Libro) ID() int          { return l.id }
func (l *Libro) Titulo() string   { return l.titulo }
func (l *Libro) Autor() string    { return l.autor }
func (l *Libro) Genero() string   { return l.genero }
func (l *Libro) Anio() int        { return l.anio }
func (l *Libro) Popularidad() int { return l.popularidad }

// Setters con Receptor de Puntero (*Libro) y Validación de Errores
func (l *Libro) SetAnio(anio int) error {
	if anio < 1450 || anio > 2026 {
		return fmt.Errorf("Año inválido %d: debe estar entre 1450 y 2026", anio)
	}
	l.anio = anio
	return nil
}

func (l *Libro) SetPopularidad(pop int) error {
	if pop < 1 || pop > 5 {
		return errors.New("la popularidad debe estar en la escala de 1 a 5 estrellas")
	}
	l.popularidad = pop
	return nil
}

// OpcionLibro define un tipo funcional para aplicar el patrón Constructor Funcional
type OpcionLibro func(*Libro) error

// ConPopularidad asigna una calificación opcional al construir el objeto
func ConPopularidad(pop int) OpcionLibro {
	return func(l *Libro) error {
		return l.SetPopularidad(pop)
	}
}

// NewLibro emula un constructor de objeto con inicialización segura y opciones dinámicas
func NewLibro(id int, titulo, autor, genero string, anio int, opciones ...OpcionLibro) (*Libro, error) {
	if titulo == "" || autor == "" {
		return nil, errors.New("el título y el autor son obligatorios")
	}

	l := &Libro{
		id:          id,
		titulo:      titulo,
		autor:       autor,
		genero:      genero,
		popularidad: 3, // Valor por defecto
	}

	if err := l.SetAnio(anio); err != nil {
		return nil, err
	}

	for _, opcion := range opciones {
		if err := opcion(l); err != nil {
			return nil, err
		}
	}

	return l, nil
}

// GenerarCatalogoInicial crea una colección inicial de punteros a Libro
func GenerarCatalogoInicial() []*Libro {
	l1, _ := NewLibro(1, "The Go Programming Language", "Alan Donovan", "Tecnología", 2015, ConPopularidad(5))
	l2, _ := NewLibro(2, "Designing Data-Intensive Applications", "Martin Kleppmann", "Tecnología", 2017, ConPopularidad(5))
	l3, _ := NewLibro(3, "Clean Code", "Robert C. Martin", "Programación", 2008, ConPopularidad(4))
	l4, _ := NewLibro(4, "1984", "George Orwell", "Sociología", 1949, ConPopularidad(5))
	l5, _ := NewLibro(5, "Sapiens: De animales a dioses", "Yuval Noah Harari", "Sociología", 2011, ConPopularidad(5))
	l6, _ := NewLibro(6, "The Pragmatic Programmer", "Andrew Hunt & David Thomas", "Programación", 1999, ConPopularidad(5))

	return []*Libro{l1, l2, l3, l4, l5, l6}
}
