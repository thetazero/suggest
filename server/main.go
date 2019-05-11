package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const SIZE int = 5

type stateMatrix struct {
	state map[[SIZE]byte]map[byte]float64
}

var model stateMatrix

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	input := loadText("source.txt")
	model = stateMatrix{state: make(map[[SIZE]byte]map[byte]float64)}
	train(model.state, input)
	// fmt.Println(model.generate())
	// ioutil.WriteFile("out.txt", []byte(model.generate()), 0644)
	http.HandleFunc("/", requestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	fmt.Println(r.RequestURI)
	querybytes := []byte(r.RequestURI)
	var query []byte
	for i := 0; i < SIZE; i++ {
		query = append(query, querybytes[i+1])
	}
	strinky := model.generate(string(query))
	fmt.Fprintf(w, strinky)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
}

func (m stateMatrix) generate(inputStr string) string {
	nxt := []byte(inputStr)
	seed := [SIZE]byte{} //{69, 109}
	for i := 0; i < SIZE; i++ {
		seed[i] = nxt[i]
	}
	ret := []byte{}
	// for i := 0; i < SIZE; i++ {
	// 	ret = append(ret, seed[i])
	// }

	var done bool
	for i := 0; i < 10000-SIZE; i++ {
		getTo := rand.Float64()
		var at float64
		if done {
			break
		}
		for key, value := range m.state[seed] {
			at += value
			if at > getTo {
				if key == 32 || key == 10 {
					done = true
					break
				}
				ret = append(ret, key)
				// seed = [SIZE]byte{seed[1], seed[2], key}
				cseed := [SIZE]byte{}
				for s := 1; s < SIZE; s++ {
					cseed[s-1] = seed[s]
				}
				cseed[SIZE-1] = key
				seed = cseed
				break
			}
		}
	}
	return string(ret)
}

func train(state map[[SIZE]byte]map[byte]float64, input []byte) {
	norm := make(map[[SIZE]byte]float64)
	for i := 0; i < len(input)-SIZE; i++ {
		key := [SIZE]byte{}
		for x := 0; x < SIZE; x++ {
			key[x] = input[i+x]
		}
		val := input[i+SIZE]
		if state[key] == nil {
			state[key] = make(map[byte]float64)
		}
		if state[key][val] == 0 {
			state[key][val] = 1
			norm[key] = 1
		} else {
			state[key][val]++
			norm[key]++
		}
		// fmt.Println(state[key][val])
	}
	for key := range state {
		for skey := range state[key] {
			state[key][skey] /= norm[key]
		}
	}

}

func loadText(path string) []byte {
	ret, _ := ioutil.ReadFile(path)
	return ret
}
