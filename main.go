package main

import "./graph"
import "fmt"
import "time"
import "strconv"
import "os"
import "bufio"
import "strings"
import "runtime"
import "flag"
import "math/rand"
import "sync"

const MaxInt = int(^uint(0) >> 1)

var DataFile string
var coreNum int
var gen int
var showRes bool

func init() {
	flag.IntVar(&coreNum, "cores", runtime.NumCPU(), "Define number of maximum processes.")
	flag.StringVar(&DataFile, "f", "MiezinasM_IND.txt", "Filename, where data is stored.")
	flag.IntVar(&gen, "g", 0, "If it is more than 0, then a graph with specified vertices count will be generated.")
	flag.BoolVar(&showRes, "p", false, "If defined, results will be printed to your console.")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(coreNum)
	var done sync.WaitGroup

	dist := map[string]map[string]int{}
	prev := map[string]map[string]string{}
	graph := graph.New()

	if gen > 0 {
		generateMatrix(gen)
	}

	readData(graph)
	fmt.Println(graph.Len())

	keys := graph.GetKeys()

	t1 := time.Now()
	for _, key := range keys {
		dist[key] = map[string]int{}
		prev[key] = map[string]string{}
		done.Add(1)
		go Dijkstra(graph, key, dist[key], prev[key], &done)
	}

	done.Wait()

	t2 := time.Since(t1)
	fmt.Println(t2)
	if showRes {
		resultsPrinter(dist, prev, keys)
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func generateMatrix(vertexes int) {
	rand.Seed(time.Now().UTC().UnixNano())

	// allocate composed 2d array
	a := make([][]int, vertexes)
	for i := range a {
		a[i] = make([]int, vertexes)
	}

	for i, el := range a {
		for j, _ := range el {
			el[j] = randInt(0, 9)
			a[j][i] = el[j]
		}
	}

	writeTofile(a)
}

func writeToFile(graph [][]int) {
	file, _ := os.Create(DataFile)
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i, line := range graph {
		vertex := ""
		for j, _ := range line {
			vertex += strconv.Itoa(graph[i][j]) + "\t"
		}
		writer.WriteString(vertex + "\n")
	}

	writer.Flush()
}

func readData(graf *graph.Graph) {

	i := 0

	file, _ := os.Open(DataFile)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		weights := strings.Split(line, "\t")
		neighbors := map[string]int{}
		for i, el := range weights {
			if weight, _ := strconv.Atoi(el); weight != 0 {
				neighbors[strconv.Itoa(i)] = weight
			}
		}
		graf.Add(strconv.Itoa(i), neighbors)
		i++
	}
}

func Dijkstra(graf *graph.Graph, source string, dist map[string]int, previous map[string]string, done *sync.WaitGroup) {
	Q := graph.New()
	dist[source] = 0
	for _, key := range graf.GetKeys() {
		if key != source {
			dist[key] = MaxInt
			previous[key] = "undefined"
		}
		vertex, _ := graf.Get(key)
		Q.Add(key, vertex.Neighbors())
	}
	j := Q.Len()
	for i := 0; i < j; i++ {
		u := minDist(dist, Q.GetKeys())
		vertex, _ := graf.Get(u)
		Q.Delete(u)
		neighbors := vertex.Neighbors()
		for key, distance := range neighbors {
			alt := dist[u] + distance
			if alt < dist[key] {
				dist[key] = alt
				previous[key] = u
			}
		}
	}
	done.Done()
}

func minDist(dist map[string]int, leftVert []string) string {
	var minKey string
	min := MaxInt
	for _, key := range leftVert {
		if dist[key] < min {
			min = dist[key]
			minKey = key
		}
	}
	return minKey
}

func resultsPrinter(dist map[string]map[string]int, prev map[string]map[string]string, keys []string) {
	for _, key := range keys {
		fmt.Println("Trumpiausi keliai iš", key, "viršūnės:")
		for index, _ := range dist[key] {
			if index != key {
				fmt.Println("Į viršūnę", index, "atkeliavome iš", prev[key][index], ". Atstumas:", dist[key][index])
			}
		}
	}
}
