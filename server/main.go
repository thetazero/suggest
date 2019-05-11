package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const SIZE int = 2

type stateMatrix struct {
	Trained map[[SIZE]string]strAndProb
}

type strAndProb struct {
	Str  string
	Prob float64
}

type trainingMatrix map[[SIZE]string]map[string]float64

var model stateMatrix

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())
	input := loadText("source/news.txt")
	// for i := 0; i < len(input); i++ {
	// 	fmt.Println(input[i])
	// }
	train(&model, input)
	fmt.Println(generate(&model, [SIZE]string{"implied", "your"}))
	http.HandleFunc("/", requestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	processedURI := strings.Replace(r.RequestURI[1:], "%20", " ", -1)
	fmt.Println(processedURI)

	splitString := strings.Fields(processedURI)
	input := [SIZE]string{}
	if len(splitString) < SIZE {
		fmt.Println("invalid request")
		return
	}
	for i := 0; i < SIZE; i++ {
		input[i] = splitString[i]
	}
	output := generate(&model, input)
	fmt.Println(output)
	fmt.Fprintf(w, strings.Join(output, " "))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
}

func generate(model *stateMatrix, input [SIZE]string) []string {
	ret := []string{}
	key := input
	var probability float64 = 1
	for true {
		if len(ret) > 6 {
			break
		}
		if _, exists := model.Trained[key]; exists {
			probability *= model.Trained[key].Prob
			fmt.Println(model.Trained[key])

			if probability > 0.1 {
				ret = append(ret, model.Trained[key].Str)
				last := model.Trained[key].Str
				for i := 0; i < len(key)-1; i++ {
					key[i] = key[i+1]
				}
				key[len(key)-1] = last
			} else {
				break
			}
		} else {
			break
		}
	}
	return ret
}

func train(model *stateMatrix, input []string) {
	matrix := make(trainingMatrix)
	for i := 0; i < len(input)-SIZE; i++ {
		key := [SIZE]string{}
		for x := 0; x < SIZE; x++ {
			key[x] = input[i+x]
		}
		val := input[i+SIZE]
		if matrix[key] == nil {
			matrix[key] = make(map[string]float64)
		}
		if matrix[key][val] == 0 {
			matrix[key][val] = 1
		} else {
			matrix[key][val]++
		}
	}
	model.Trained = make(map[[SIZE]string]strAndProb)
	for key := range matrix {
		var sum float64
		var biggestSkey string
		var biggestVal float64
		for skey := range matrix[key] {
			sum += matrix[key][skey]
			if matrix[key][skey] > biggestVal {
				biggestVal = matrix[key][skey]
				biggestSkey = skey
			}
		}
		model.Trained[key] = strAndProb{biggestSkey, biggestVal / sum}
		// fmt.Println(key, sum, biggestVal, biggestSkey)
	}
	// fmt.Println(model.Trained)
}

func loadText(path string) []string {
	unProcessed, _ := ioutil.ReadFile(path)
	ret := strings.Fields(string(unProcessed))
	for i := 0; i < len(ret); i++ {
		ret[i] = strings.ToLower(ret[i])
	}
	return ret
}
