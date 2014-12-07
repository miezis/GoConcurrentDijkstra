/*
============================================================================
Mantas Miežinas, IFF-2
Individualus darbas (main.go)
============================================================================
*/

package main

import (
	"./graph"
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

//maksimali integer tipo reikšmė
const MaxInt = int(^uint(0) >> 1)

var DataFile string //duomenų failo pavadinimas
var coreNum int     //procesų kiekis
var gen int         //keliu viršūnių grafą generuoti
var showRes bool    //ar spausdinti rezultatus

//Funkcija, ivykdoma pati pirma, pasiima komandines eilutes parametrus
func init() {
	//nustatome kiek procesu maksimaliai gales dirbti vienu metu
	flag.IntVar(&coreNum, "cores", runtime.NumCPU(), "Define number of maximum processes.")
	//nustatome pradiniu duomenu failo pavadinima
	flag.StringVar(&DataFile, "f", "MiezinasM_IND.txt", "Filename, where data is stored.")
	//nustatome generuojamo grafo virsuniu kieki, jei 0 - grafo negeneruojame, o skaitome is duomenu failo
	flag.IntVar(&gen, "g", 0, "If it is more than 0, then a graph with specified vertices count will be generated.")
	//nustato ar spausdinti rezultatus i terminala, ar ne
	flag.BoolVar(&showRes, "p", false, "If defined, results will be printed to your console.")
	//nuskaitome komandines eilutes parametrus
	flag.Parse()
}

func main() {
	//nustatome maksimalų procesų kiekį
	runtime.GOMAXPROCS(coreNum)
	//kintamasis sinchronizacijai
	var done sync.WaitGroup
	//kintamieji, į kuriuos dėsime rezultatus
	dist := map[string]map[string]int{}
	prev := map[string]map[string]string{}
	//inicializuojam naują grafą, saugome nuorodą į jį
	graph := graph.New()

	//jei nurodyta, tai generuojame gretimumo matricą
	//sugeneruota, ji įrašoma į failą
	if gen > 0 {
		generateMatrix(gen)
	}

	//nuskaitome duomenis, sudedame į grafo struktūra
	readData(graph)

	fmt.Println("Procesų kiekis:", coreNum)
	fmt.Println("Viršūnių kiekis:", graph.Len())

	//keys kintamajame saugome visų viršūnių pavadinimus
	keys := graph.GetKeys()

	//pradedame skaičiuoti vykdymo laiką
	t1 := time.Now()
	//einame per viršūnes, visoms pritaikome Dijkstra
	for _, key := range keys {
		//inicializuojame rezultatų struktūras
		dist[key] = map[string]int{}
		prev[key] = map[string]string{}
		//į waitgroup struktūrą pridedame vieną veiksmą
		done.Add(1)
		//kviečiame Dijkstra f-ja kaip go routine
		go Dijkstra(graph, key, dist[key], prev[key], &done)
	}
	//laukiame, kol visi procesai atliks savo darbą
	done.Wait()
	//apskaičiuojame kiek truko kelių radimas
	t2 := time.Since(t1)

	//jei buvo nurodyta, išvedame rezultatus į terminalą
	if showRes {
		resultsPrinter(dist, prev, keys)
	}
	fmt.Println("Vykdymo laikas:", t2)
}

//Pagalbinė f-ja, gražinanti atsitiktinį sk. iš [min, max] intervalo
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

/*
============================================================================
generateMatrix
	Sugeneruoja grafo viršūnių gretimumo matricą. Funkcijos argumentas yra
	viršūnių kiekis. Sugeneravus kreipiamasi į writeToFile funkciją ir
	matrica įrašoma į failą.
============================================================================
*/
func generateMatrix(vertexes int) {
	rand.Seed(time.Now().UTC().UnixNano())

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

	writeToFile(a)
}

/*
============================================================================
writeToFile
	Priima grafo gretimumo matricą ir ją išspausdina į failą.
============================================================================
*/
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

/*
============================================================================
readData
	Nuskaito duomenis iš failo, surašo juos į grafo struktūrą ir gražina
	per nuorodą.
============================================================================
*/
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

/*
============================================================================
Dijkstra
	Realizuotas Dijkstra algoritmas, rezultatus surašo į map struktūras,
	kurios pagal nutylėjimą yra nuorodos tipo, sinchronizacijai naudojame
	WaitGroup struktūra iš sync paketo.
============================================================================
*/
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

/*
============================================================================
minDist
	Priima atstumų map struktūrą ir likusių viršūnių slice struktūrą,
	randą viršūnę iki kurios kelias trumpiausias ir gražina jos pavadinimą.
============================================================================
*/
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

/*
============================================================================
resultsPrinter
	Išspausdina rezultatus į terminalo langą.
============================================================================
*/
func resultsPrinter(dist map[string]map[string]int, prev map[string]map[string]string, keys []string) {
	for _, key := range keys {
		fmt.Println("Trumpiausi keliai iš", key, "viršūnės:")
		for index, _ := range dist[key] {
			if index != key {
				if prev[key][index] != "undefined" {
					fmt.Println("Į viršūnę", index, "atkeliavome iš", prev[key][index], ". Atstumas:", dist[key][index])
				} else {
					fmt.Println("Į viršūnę", index, "kelio nėra.")
				}
			}
		}
	}
}
