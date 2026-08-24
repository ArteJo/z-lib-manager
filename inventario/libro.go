package inventario

import (
	"encoding/json"
	"errors"
	"fmt"
)

// MetadatosPublicacion implementa incrustación de estructuras (Struct Embedding)
type MetadatosPublicacion struct {
	Editorial string `json:"editorial"`
	Edicion   int    `json:"edicion"`
}

// Libro implementa encapsulación estricta con campos privados (en minúscula)
// inaccesibles desde fuera del paquete inventario.
type Libro struct {
	id                   int
	titulo               string
	autor                string
	genero               string
	anio                 int
	popularidad          int // Calificación de 1 a 5 estrellas
	MetadatosPublicacion     // Incrustación de estructura para metadatos
}

// Getters (Métodos de lectura)
func (l *Libro) ID() int          { return l.id }
func (l *Libro) Titulo() string   { return l.titulo }
func (l *Libro) Autor() string    { return l.autor }
func (l *Libro) Genero() string   { return l.genero }
func (l *Libro) Anio() int        { return l.anio }
func (l *Libro) Popularidad() int { return l.popularidad }

// Setters con Receptor de Puntero (*Libro) y Validación de Errores
func (l *Libro) SetTitulo(titulo string) error {
	if titulo == "" {
		return errors.New("el título no puede estar vacío")
	}
	l.titulo = titulo
	return nil
}

func (l *Libro) SetAutor(autor string) error {
	if autor == "" {
		return errors.New("el autor no puede estar vacío")
	}
	l.autor = autor
	return nil
}

func (l *Libro) SetGenero(genero string) error {
	if genero == "" {
		return errors.New("el género no puede estar vacío")
	}
	l.genero = genero
	return nil
}

func (l *Libro) SetAnio(anio int) error {
	if anio < 1450 || anio > 2026 {
		return fmt.Errorf("año inválido %d: debe estar comprendido entre 1450 y 2026", anio)
	}
	l.anio = anio
	return nil
}

func (l *Libro) SetPopularidad(pop int) error {
	if pop < 1 || pop > 5 {
		return errors.New("la popularidad debe pertenecer a la escala de 1 a 5 estrellas")
	}
	l.popularidad = pop
	return nil
}

// MarshalJSON implementa la interfaz json.Marshaler para serializar campos privados
func (l *Libro) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID          int    `json:"id"`
		Titulo      string `json:"titulo"`
		Autor       string `json:"autor"`
		Genero      string `json:"genero"`
		Anio        int    `json:"anio"`
		Popularidad int    `json:"popularidad"`
		Editorial   string `json:"editorial,omitempty"`
		Edicion     int    `json:"edicion,omitempty"`
	}{
		ID:          l.id,
		Titulo:      l.titulo,
		Autor:       l.autor,
		Genero:      l.genero,
		Anio:        l.anio,
		Popularidad: l.popularidad,
		Editorial:   l.Editorial,
		Edicion:     l.Edicion,
	})
}

// UnmarshalJSON implementa la interfaz json.Unmarshaler para deserialización segura
func (l *Libro) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID          int    `json:"id"`
		Titulo      string `json:"titulo"`
		Autor       string `json:"autor"`
		Genero      string `json:"genero"`
		Anio        int    `json:"anio"`
		Popularidad int    `json:"popularidad"`
		Editorial   string `json:"editorial"`
		Edicion     int    `json:"edicion"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	l.id = aux.ID
	l.Editorial = aux.Editorial
	l.Edicion = aux.Edicion

	if err := l.SetTitulo(aux.Titulo); err != nil {
		return err
	}
	if err := l.SetAutor(aux.Autor); err != nil {
		return err
	}
	if err := l.SetGenero(aux.Genero); err != nil {
		return err
	}
	if err := l.SetAnio(aux.Anio); err != nil {
		return err
	}
	if aux.Popularidad != 0 {
		if err := l.SetPopularidad(aux.Popularidad); err != nil {
			return err
		}
	} else {
		l.popularidad = 3
	}

	return nil
}

// OpcionLibro define la firma para el Patrón Constructor Funcional
type OpcionLibro func(*Libro) error

func ConPopularidad(pop int) OpcionLibro {
	return func(l *Libro) error {
		return l.SetPopularidad(pop)
	}
}

func ConEditorial(editorial string, edicion int) OpcionLibro {
	return func(l *Libro) error {
		l.Editorial = editorial
		l.Edicion = edicion
		return nil
	}
}

// NewLibro instancia un libro con validación y opciones funcionales
func NewLibro(id int, titulo, autor, genero string, anio int, opciones ...OpcionLibro) (*Libro, error) {
	l := &Libro{
		id:          id,
		popularidad: 3, // Valor por defecto
	}

	if err := l.SetTitulo(titulo); err != nil {
		return nil, err
	}
	if err := l.SetAutor(autor); err != nil {
		return nil, err
	}
	if err := l.SetGenero(genero); err != nil {
		return nil, err
	}
	if err := l.SetAnio(anio); err != nil {
		return nil, err
	}

	for _, opt := range opciones {
		if err := opt(l); err != nil {
			return nil, err
		}
	}

	return l, nil
}

// GenerarCatalogoInicial crea una colección inicial de punteros a Libro (provee datos base)
// GenerarCatalogoInicial crea una colección inicial de punteros a Libro (provee datos base)
func GenerarCatalogoInicial() []*Libro {
	// Libros base existentes (1 - 6)
	l1, _ := NewLibro(1, "The Go Programming Language", "Alan Donovan", "Tecnología", 2015, ConPopularidad(5))
	l2, _ := NewLibro(2, "Designing Data-Intensive Applications", "Martin Kleppmann", "Tecnología", 2017, ConPopularidad(5))
	l3, _ := NewLibro(3, "Clean Code", "Robert C. Martin", "Programación", 2008, ConPopularidad(4))
	l4, _ := NewLibro(4, "1984", "George Orwell", "Sociología", 1949, ConPopularidad(5))
	l5, _ := NewLibro(5, "Sapiens: De animales a dioses", "Yuval Noah Harari", "Sociología", 2011, ConPopularidad(5))
	l6, _ := NewLibro(6, "The Pragmatic Programmer", "Andrew Hunt & David Thomas", "Programación", 1999, ConPopularidad(5))
	l7, _ := NewLibro(7, "Introduction to Cybersecurity", "Go", "Seguridad", 2026, ConPopularidad(5))
	l8, _ := NewLibro(8, "Network Defense Essentials", "Jane Smith", "Seguridad", 2023, ConPopularidad(4))
	l9, _ := NewLibro(9, "Practical Malware Analysis", "Michael Sikorski", "Seguridad", 2012, ConPopularidad(5))
	l10, _ := NewLibro(10, "The Art of Invisibility", "Kevin Mitnick", "Seguridad", 2017, ConPopularidad(4))
	l11, _ := NewLibro(11, "Refactoring", "Martin Fowler", "Programación", 1999, ConPopularidad(5))
	l12, _ := NewLibro(12, "Head First Design Patterns", "Eric Freeman", "Programación", 2004, ConPopularidad(4))
	l13, _ := NewLibro(13, "Computer Networking: A Top-Down Approach", "James Kurose", "Redes", 2021, ConPopularidad(5))
	l14, _ := NewLibro(14, "TCP/IP Illustrated", "W. Richard Stevens", "Redes", 1994, ConPopularidad(5))
	l15, _ := NewLibro(15, "Site Reliability Engineering", "Betsy Beyer", "Tecnología", 2016, ConPopularidad(4))
	l16, _ := NewLibro(16, "Building Microservices", "Sam Newman", "Tecnología", 2021, ConPopularidad(5))
	l17, _ := NewLibro(17, "Linux Bible", "Christopher Negus", "Sistemas", 2020, ConPopularidad(4))
	l18, _ := NewLibro(18, "Modern Operating Systems", "Andrew S. Tanenbaum", "Sistemas", 2014, ConPopularidad(5))
	l19, _ := NewLibro(19, "Brave New World", "Aldous Huxley", "Sociología", 1932, ConPopularidad(4))
	l20, _ := NewLibro(20, "Fahrenheit 451", "Ray Bradbury", "Sociología", 1953, ConPopularidad(4))
	l21, _ := NewLibro(21, "Ghost in the Wires", "Kevin Mitnick", "Seguridad", 2011, ConPopularidad(5))

	return []*Libro{
		l1, l2, l3, l4, l5, l6, l7, l8, l9, l10,
		l11, l12, l13, l14, l15, l16, l17, l18, l19, l20, l21,
	}
}
