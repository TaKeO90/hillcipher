package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO: Implement encrypt function

const letters string = "abcdefghijklmnopqrstuvwxyz_"

type Key [][]int

func (k *Key) Gen(n int) {

	*k = make(Key, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			rand.Seed(time.Now().UnixNano())
			// FIXME: hardcoded max rand range
			(*k)[i] = append((*k)[i], rand.Intn(40)%26)
		}
	}

}

func (k *Key) String() string {
	var res string
	for _, n := range *k {
		res += strings.Join(Stack(n).ToLetters(), "")
	}
	return res
}

type Stack []int

func (s Stack) Sum() (sum int) {
	for _, n := range s {
		sum += n
	}
	return
}

func (s Stack) ToLetters() (out []string) {
	for _, n := range s {
		var index int = 0
		strings.Map(func(r rune) rune {
			if n == index {
				out = append(out, string(r))
			}
			index++
			return 0
		}, letters)
	}

	return
}

func (k *Key) Mult(pair []int) (result Stack) {
	x := 0
	stack := Stack{}
	for i := range *k {
		for j := range (*k)[i] {
			y := j % len((*k)[i])
			pairIndex := j % len(pair)
			stack = append(stack, (*k)[x][y]*pair[pairIndex])
			if j == len((*k)[i])-1 {
				result = append(result, stack.Sum()%26)
				stack = Stack{}
				x++
			}
		}
	}
	return
}

func wordToPairs(word string, keySize int) (pairs [][]int) {

	getLetterIndex := func(p string) (index int) {
		count := 0
		strings.Map(func(r rune) rune {
			if string(r) == p {
				index = count
				return 0
			}
			count++
			return 0
		}, letters)
		return
	}

	tmp := []int{}
	for _, l := range word {
		if len(tmp) == keySize {
			pairs = append(pairs, tmp)
			tmp = []int{}
		}
		idx := getLetterIndex(strings.ToLower(string(l)))
		tmp = append(tmp, idx)
	}
	pairs = append(pairs, tmp)
	return
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("USAGE")
		os.Exit(1)
	}

	n, _ := strconv.Atoi(os.Args[1])

	key := new(Key)

	key.Gen(n)
	fmt.Println("Generated the key => ", key.String())

	word := os.Args[2]
	pairs := wordToPairs(word, n)

	fmt.Println("Pairs", pairs)

	var str string
	for _, p := range pairs {
		str += strings.Join(key.Mult(p).ToLetters(), "")
	}

	fmt.Println(str)
}
