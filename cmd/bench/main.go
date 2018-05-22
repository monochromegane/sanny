package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/monochromegane/sanny"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	yaml "gopkg.in/yaml.v2"
)

var (
	configPath string
	dataPath   string
	testSize   int
	outPath    string
	runs       int
	count      int
	algo       string
	seed       int64
	components int
)

func init() {
	flag.StringVar(&configPath, "config", "algos.yaml", "Algorithm definitions file path")
	flag.StringVar(&algo, "algo", "sanny", "brute_force|brute_force_blas|annoy|sanny")
	flag.StringVar(&dataPath, "data", "", "Input data file path")
	flag.IntVar(&testSize, "test-size", 500, "test size")
	flag.IntVar(&runs, "runs", 3, "Run each algorithm")
	flag.IntVar(&count, "count", 10, "Number of near neighbours to search for")
	flag.StringVar(&outPath, "out", "results", "Output directory path")
	flag.Int64Var(&seed, "seed", 1, "Random seed for shuffle. If seed < 0 then use random seed")
	flag.IntVar(&components, "components", 0, "Compress dimension by PCA.")
	flag.Parse()
}

func main() {
	fmt.Printf("Loading data... (%s)\n", dataPath)
	x, err := loadFrom(dataPath)
	if err != nil {
		panic(err)
	}
	shuffle(x)

	queries := x[:testSize]
	data := x[testSize:]

	fmt.Printf("Computing truth...\n")
	truth := computeTruth(queries, data)

	if components > 0 {
		fmt.Printf("Compressing...\n")
		x = pca(x)
		queries = x[:testSize]
		data = x[testSize:]
	}

	fmt.Printf("Running benchmarks...\n")
	switch algo {
	case "brute_force":
		benchBruteForce(queries, data, truth)
	case "brute_force_blas":
		benchBruteForceBLAS(queries, data, truth)
	case "annoy":
		benchAnnoy(queries, data, truth)
	case "sanny":
		benchSanny(queries, data, truth)
	default:
		benchBruteForce(queries, data, truth)
	}
}

func benchBruteForce(queries, data [][]float32, truth [][]int) {
	algo := sanny.NewBruteForce()
	runner := Runner{
		Name: "BruteForce",
		Algo: algo,
	}
	fmt.Printf("Building %s\n", runner.Name)
	algo.Build(data)
	fmt.Printf("%s\n", fmt.Sprintf("data: %d", len(data)))
	recall, qps := runner.Run(truth, queries, data)
	writeTo(runner.Name, recall, qps)
}

func benchBruteForceBLAS(queries, data [][]float32, truth [][]int) {
	algo := sanny.NewBruteForceBLAS()
	runner := Runner{
		Name: "BruteForceBLAS",
		Algo: algo,
	}
	fmt.Printf("Building %s\n", runner.Name)
	algo.Build(data)
	fmt.Printf("%s\n", fmt.Sprintf("data: %d", len(data)))
	recall, qps := runner.Run(truth, queries, data)
	writeTo(runner.Name, recall, qps)
}

func benchAnnoy(queries, data [][]float32, truth [][]int) {
	config, _ := loadConfig(configPath)
	for _, tree := range config.Args[0] {
		algo := sanny.NewAnnoy(tree, 0)

		runner := Runner{
			Name: "Annoy",
			Algo: algo,
		}

		fmt.Printf("Building %s\n", runner.Name)
		algo.Build(data)

		for _, searchK := range config.Args[1] {
			fmt.Printf("%s\n", fmt.Sprintf("tree: %d, searchK: %d", tree, searchK))
			algo.SearchK = searchK
			recall, qps := runner.Run(truth, queries, data)
			writeTo(runner.Name, recall, qps)
		}
	}
}

func benchSanny(queries, data [][]float32, truth [][]int) {
	config, _ := loadConfig(configPath)
	for _, splitNum := range config.Args[0] {
		indecies := make([][]int, splitNum)
		for i, _ := range indecies {
			indecies[i] = []int{i}
		}
		for _, tree := range config.Args[2] {
			searchers := make([]sanny.Searcher, splitNum)
			for i, _ := range searchers {
				searchers[i] = sanny.NewAnnoy(tree, 0)
			}
			algo := sanny.NewSanny(splitNum, 0, true, searchers, indecies)
			runner := Runner{
				Name: "Sanny",
				Algo: algo,
			}
			fmt.Printf("Building %s\n", runner.Name)
			algo.Build(data)
			for _, top := range config.Args[1] {
				for _, searchK := range config.Args[3] {
					fmt.Printf("%s\n", fmt.Sprintf("split: %d, top: %d, tree: %d, searchK: %d", splitNum, top, tree, searchK))
					algo.(*sanny.Sanny).Top = top
					for i, _ := range searchers {
						searchers[i].(*sanny.Annoy).SearchK = searchK
					}
					recall, qps := runner.Run(truth, queries, data)
					writeTo(runner.Name, recall, qps)
				}
			}
		}
	}
}

func computeTruth(queries, data [][]float32) [][]int {
	bf := sanny.NewBruteForce()
	bf.Build(data)

	truth := make([][]int, testSize)
	for i, q := range queries {
		if i%50 == 0 {
			fmt.Printf("%d\n", i)
		}
		truth[i] = bf.SearchConcurrent(q, count)
	}
	return truth
}

func loadFrom(path string) ([][]float32, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var x [][]float32
	reader := csv.NewReader(f)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		item := make([]float32, len(record))
		for i, c := range record {
			value, err := strconv.ParseFloat(c, 32)
			if err != nil {
				return nil, err
			}
			item[i] = float32(value)
		}
		x = append(x, item)
	}
	return x, nil
}

func shuffle(data [][]float32) {
	rnd := rand.New(rand.NewSource(seed))
	if seed < 0 {
		rnd.Seed(time.Now().UnixNano())
	}
	n := len(data)
	for i := n - 1; i >= 0; i-- {
		j := rnd.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

type Runner struct {
	Name        string
	Description string
	Algo        sanny.Searcher
}

func (r *Runner) Run(truth [][]int, queries, data [][]float32) (float64, time.Duration) {
	bestSearchTime := time.Duration(math.MaxInt32)
	match := 0
	for i := 0; i < runs; i++ {
		fmt.Printf("Running queries %d\n", i)
		start := time.Now()
		results := make([][]int, len(queries))
		for i, q := range queries {
			results[i] = r.Algo.Search(q, count)
		}
		searchTime := (time.Now().Sub(start)) / time.Duration(len(queries))
		if searchTime < bestSearchTime {
			bestSearchTime = searchTime
		}
		for j, result := range results {
			for _, id := range result {
				for _, tid := range truth[j] {
					if id == tid {
						match += 1
						break
					}
				}
			}
		}
	}
	return float64(match) / float64(len(queries)*runs*count), bestSearchTime
}

func writeTo(name string, recall float64, qps time.Duration) {
	os.Mkdir(outPath, 0755)
	file, err := os.OpenFile(fmt.Sprintf("%s/%s.csv", outPath, name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{fmt.Sprintf("%f", recall), fmt.Sprintf("%f", float64(qps)/1000000000)})
	writer.Flush()
}

type Config struct {
	Algos []Algo `yaml:"algos"`
}

type Algo struct {
	Name string  `yaml:"name"`
	Args [][]int `yaml:"args"`
}

func loadConfig(path string) (Algo, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Algo{}, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Algo{}, err
	}

	for _, c := range config.Algos {
		if c.Name == algo {
			return c, nil
		}
	}
	return Algo{}, fmt.Errorf("Unknown algorism")
}

func pca(data [][]float32) [][]float32 {
	y := mat.NewDense(len(data), len(data[0]), nil)
	for i := 0; i < len(data); i++ {
		r := make([]float64, len(data[i]))
		for j, _ := range data[i] {
			r[j] = float64(data[i][j])
		}
		row := mat.NewVecDense(len(data[i]), r)
		y.SetRow(i, mat.Col(nil, 0, row))
	}

	var pc stat.PC
	ok := pc.PrincipalComponents(y, nil)
	if !ok {
		panic("PCA fails")
	}

	k := components
	var proj mat.Dense
	proj.Mul(y, pc.VectorsTo(nil).Slice(0, len(data[0]), 0, k))

	x := make([][]float32, len(data))
	for i := 0; i < len(data); i++ {
		x[i] = make([]float32, k)
		for j := 0; j < k; j++ {
			x[i][j] = float32(proj.At(i, j))
		}
	}
	return x
}
