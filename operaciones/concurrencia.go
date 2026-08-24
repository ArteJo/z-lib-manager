package operaciones

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ArteJo/z-lib-manager/inventario"
)

// RepositorioConcurrente protege el acceso multihilo a la memoria compartida mediante sync.RWMutex
type RepositorioConcurrente struct {
	mu       sync.RWMutex
	libros   map[int]*inventario.Libro
	ultimoID int
}

func NuevoRepositorioConcurrente(inicial []*inventario.Libro) *RepositorioConcurrente {
	repo := &RepositorioConcurrente{
		libros: make(map[int]*inventario.Libro),
	}
	for _, l := range inicial {
		repo.libros[l.ID()] = l
		if l.ID() > repo.ultimoID {
			repo.ultimoID = l.ID()
		}
	}
	return repo
}

// ObtenerTodos realiza lectura concurrente compartida (RLock)
func (r *RepositorioConcurrente) ObtenerTodos() []*inventario.Libro {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var lista []*inventario.Libro
	for _, l := range r.libros {
		lista = append(lista, l)
	}
	return lista
}

// ObtenerPorID
func (r *RepositorioConcurrente) ObtenerPorID(id int) (*inventario.Libro, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	l, existe := r.libros[id]
	return l, existe
}

// Guardar realiza bloqueo exclusivo de escritura (Lock)
func (r *RepositorioConcurrente) Guardar(l *inventario.Libro) *inventario.Libro {
	r.mu.Lock()
	defer r.mu.Unlock()

	if l.ID() == 0 {
		r.ultimoID++
		// Recrear con el nuevo ID
		nuevoLibro, _ := inventario.NewLibro(r.ultimoID, l.Titulo(), l.Autor(), l.Genero(), l.Anio(),
			inventario.ConPopularidad(l.Popularidad()),
			inventario.ConEditorial(l.Editorial, l.Edicion))
		r.libros[r.ultimoID] = nuevoLibro
		return nuevoLibro
	}

	r.libros[l.ID()] = l
	if l.ID() > r.ultimoID {
		r.ultimoID = l.ID()
	}
	return l
}

// Eliminar
func (r *RepositorioConcurrente) Eliminar(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, existe := r.libros[id]; !existe {
		return false
	}
	delete(r.libros, id)
	return true
}

// BusquedaParalela ejecuta la consulta distribuida a través de Goroutines y Canales
func (r *RepositorioConcurrente) BusquedaParalela(ctx context.Context, consulta string) ([]*inventario.Libro, error) {
	// Validación defensiva para evitar pánicos por contexto nulo
	if ctx == nil {
		ctx = context.Background()
	}

	libros := r.ObtenerTodos()
	chResultados := make(chan *inventario.Libro, len(libros))
	var wg sync.WaitGroup

	termino := strings.ToLower(consulta)

	for _, l := range libros {
		wg.Add(1)
		go func(libro *inventario.Libro) {
			defer wg.Done()
			time.Sleep(5 * time.Millisecond) // Simulación de procesamiento concurrente
			if strings.Contains(strings.ToLower(libro.Titulo()), termino) ||
				strings.Contains(strings.ToLower(libro.Autor()), termino) ||
				strings.Contains(strings.ToLower(libro.Genero()), termino) {
				chResultados <- libro
			}
		}(l)
	}

	go func() {
		wg.Wait()
		close(chResultados)
	}()

	var encontrados []*inventario.Libro
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case libro, ok := <-chResultados:
			if !ok {
				return encontrados, nil
			}
			encontrados = append(encontrados, libro)
		}
	}
}
