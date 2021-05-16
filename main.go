package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const letters string = "abcdefghijklmnopqrstuvwxyz "

// InfoLogger ...
var InfoLogger *log.Logger = log.New(os.Stdout, "[INFO] ", log.Ldate)

// HillCipher ...
type HillCipher interface {
	Encrypt() string
	Decrypt() string
}

// Cipher structure holds the items
// that we need to generate our cipher key and encrypt or decrypt
type Cipher struct {
	keysize int
	word    string
	genKey  Key
}

// Encrypt hillcipher encryption.
func (c *Cipher) Encrypt() (result string) {
	pairs := wordToPairs(c.word, c.keysize)
	for _, pair := range pairs {
		result += strings.Join(c.genKey.Mult(pair).ToLetters(), "")
	}
	return
}

// Decrypt hillcipher decryption.
func (c *Cipher) Decrypt() (result string) {
	//TODO: Find Determinant of Key.
	//TODO: Transpose key matrix.
	//TODO: Find Minor.
	//TODO: Find Co-Factor.
	return
}

// NewCipher ....
func NewCipher(size int, word string) (hillCipher HillCipher) {

	key := new(Key)
	key.Gen(size)

	InfoLogger.Printf("Generated Key %s", key.String())

	cipher := &Cipher{
		keysize: size,
		word:    word,
		genKey:  *key,
	}
	hillCipher = cipher

	return
}

// Key ...
type Key [][]int

// Gen generate new key with a given size.
func (k *Key) Gen(n int) {

	*k = make(Key, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			rand.Seed(time.Now().UnixNano())
			// FIXME: hardcoded max rand range
			(*k)[i] = append((*k)[i], rand.Intn(26))
		}
	}

}

// String key into string.
func (k *Key) String() string {
	var res string
	for _, n := range *k {
		res += strings.Join(Stack(n).ToLetters(), "")
	}
	return res
}

// Stack type alias of array of int.
type Stack []int

// Sum sum stack items.
func (s Stack) Sum() (sum int) {
	for _, n := range s {
		sum += n
	}
	return
}

// GetDet get determinant using the integers that are stored in the stack like type.
func (s Stack) GetDet() int {
	return 0
}

// ToLetters Stack type into array of strings.
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

// Mult Key method for multiplying matrices.
func (k *Key) Mult(pair []int) (result Stack) {

	if len(pair) < len(*k) {
		rng := len(*k) - len(pair)
		for i := 0; i < rng; i++ {
			pair = append(pair, 26)
		}
	}

	x := 0
	stack := Stack{}
	for i := range *k {
		for j := range (*k)[i] {
			y := j % len((*k)[i])
			pairIndex := j % len(pair)
			stack = append(stack, (*k)[x][y]*pair[pairIndex])
			if j == len((*k)[i])-1 {
				result = append(result, stack.Sum()%len(letters))
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
		// TODO: update USAGE after introducing decryption
		//       at the moment we only support encryption.
		InfoLogger.Fatal("USAGE ./hillcipher <key size> <word>")
	}

	n, _ := strconv.Atoi(os.Args[1])
	word := os.Args[2]

	// Encryption test.
	fmt.Println(NewCipher(n, word).Encrypt())
}
