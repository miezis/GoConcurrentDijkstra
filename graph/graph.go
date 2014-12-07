package graph

import "sync"
import "errors"

type Vertex struct {
	key       string
	neighbors map[string]int
	sync.RWMutex
}

func (v *Vertex) Neighbors() map[string]int {
	if v == nil {
		return nil
	}

	v.RLock()
	neighbors := v.neighbors
	v.RUnlock()

	return neighbors
}

func (v *Vertex) Key() string {
	if v == nil {
		return ""
	}

	v.RLock()
	key := v.key
	v.RUnlock()

	return key
}

type Graph struct {
	vertexes map[string]*Vertex
	sync.RWMutex
}

func New() *Graph {
	return &Graph{map[string]*Vertex{}, sync.RWMutex{}}
}

func (g *Graph) Len() int {
	return len(g.vertexes)
}

func (g *Graph) Add(key string, adjacents map[string]int) bool {
	g.Lock()
	defer g.Unlock()

	if g.get(key) != nil {
		return false
	}

	g.vertexes[key] = &Vertex{key, adjacents, sync.RWMutex{}}

	return true
}

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

func (g *Graph) Get(key string) (v *Vertex, err error) {
	g.RLock()
	v = g.get(key)
	g.RUnlock()

	if v == nil {
		err = errors.New("graph: invalid key")
	}

	return
}

func (g *Graph) GetKeys() (keys []string) {
	g.Lock()
	defer g.Unlock()

	for k := range g.vertexes {
		keys = append(keys, k)
	}

	return keys
}

func (g *Graph) get(key string) *Vertex {
	return g.vertexes[key]
}
