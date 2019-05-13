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
	Trained map[[SIZE]int]strAndProb
	Count   int
	Encoder map[string]int
	Decoder map[int]string
}

type strAndProb struct {
	Str  int
	Prob float64
}

type trainingMatrix map[[SIZE]int]map[int]float64

var model stateMatrix

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())
	input := loadText("source/myemails.txt")
	// for i := 0; i < len(input); i++ {
	// 	fmt.Println(input[i])
	// }
	train(&model, input)
	fmt.Println("trained")
	// fmt.Println(generate(&model, [SIZE]int{encode("implied"), encode("your")}))
	http.HandleFunc("/", requestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func encode(s string) int {
	if model.Encoder[s] == 0 {
		model.Count++
		model.Encoder[s] = model.Count
		model.Decoder[model.Count] = s
		return model.Count
	}
	return model.Encoder[s]
}

func decode(num int) string {
	return model.Decoder[num]
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	processedURI := strings.Replace(r.RequestURI[1:], "%20", " ", -1)
	fmt.Println(processedURI)

	splitString := strings.Fields(processedURI)
	input := [SIZE]int{}
	if len(splitString) < SIZE {
		fmt.Println("invalid request")
		return
	}
	for i := 0; i < SIZE; i++ {
		input[i] = encode(splitString[i])
	}
	output := generate(&model, input)
	fmt.Println(output)
	strOutput := []string{}
	for i := 0; i < len(output); i++ {
		strOutput = append(strOutput, decode(output[i]))
	}
	fmt.Fprintf(w, strings.Join(strOutput, " "))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
}

func generate(model *stateMatrix, input [SIZE]int) []int {
	ret := []int{}
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
	model.Encoder = make(map[string]int)
	model.Decoder = make(map[int]string)
	matrix := make(trainingMatrix)
	for i := 0; i < len(input)-SIZE; i++ {
		key := [SIZE]int{}
		for x := 0; x < SIZE; x++ {
			key[x] = encode(input[i+x])
		}
		val := encode(input[i+SIZE])
		if matrix[key] == nil {
			matrix[key] = make(map[int]float64)
		}
		if matrix[key][val] == 0 {
			matrix[key][val] = 1
		} else {
			matrix[key][val]++
		}
	}
	model.Trained = make(map[[SIZE]int]strAndProb)
	for key := range matrix {
		var sum float64
		var biggestSkey int
		var biggestVal float64
		for skey := range matrix[key] {
			sum += matrix[key][skey]
			if matrix[key][skey] > biggestVal {
				biggestVal = matrix[key][skey]
				biggestSkey = skey
			}
		}
		if biggestVal > 0 {
			model.Trained[key] = strAndProb{biggestSkey, biggestVal / sum}
		}
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
