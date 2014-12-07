/*
============================================================================
Mantas Miežinas, IFF-2
Individualus darbas (graph.go)
============================================================================
*/

package graph

import (
	"errors"
	"sync"
)

//Viršūnės struktūra
type Vertex struct {
	key       string
	neighbors map[string]int
	sync.RWMutex
}

//Gražina viršūnės visas kaimynines viršūnes
func (v *Vertex) Neighbors() map[string]int {
	if v == nil {
		return nil
	}

	v.RLock()
	neighbors := v.neighbors
	v.RUnlock()

	return neighbors
}

//Gražina viršūnės pavadinimą
func (v *Vertex) Key() string {
	if v == nil {
		return ""
	}

	v.RLock()
	key := v.key
	v.RUnlock()

	return key
}

//Grafo struktūra
type Graph struct {
	vertexes map[string]*Vertex
	sync.RWMutex
}

//Sukuria naują grafą ir gražina nuorodą į jį
func New() *Graph {
	return &Graph{map[string]*Vertex{}, sync.RWMutex{}}
}

//Gražina viršūnių kiekį
func (g *Graph) Len() int {
	return len(g.vertexes)
}

//Prideda nauja viršūnę, kurios pavadinimą ir kaimynines
//viršūnes paduodame per parametrus, jei pridejo gražina true
func (g *Graph) Add(key string, adjacents map[string]int) bool {
	g.Lock()
	defer g.Unlock()

	if g.get(key) != nil {
		return false
	}

	g.vertexes[key] = &Vertex{key, adjacents, sync.RWMutex{}}

	return true
}

//Trina nurodytą viršūnę, jei pavyko gražina true
func (g *Graph) Delete(key string) bool {
	g.Lock()
	defer g.Unlock()

	v := g.get(key)
	if v == nil {
		return false
	}

	delete(g.vertexes, key)

	return true
}

//Gražina viršūnę, pagal jos pavadinimą, jei nepavyko, tai
//gražina error'ą
func (g *Graph) Get(key string) (v *Vertex, err error) {
	g.RLock()
	v = g.get(key)
	g.RUnlock()

	if v == nil {
		err = errors.New("graph: invalid key")
	}

	return
}

//Gražina viršūnių pavadinimų masyvą
func (g *Graph) GetKeys() (keys []string) {
	g.Lock()
	defer g.Unlock()

	for k := range g.vertexes {
		keys = append(keys, k)
	}

	return keys
}

//vidinė funkcija, gražinanti viršūnę
//naudoti tik RLock(), RUnlock() kontekste
func (g *Graph) get(key string) *Vertex {
	return g.vertexes[key]
}
