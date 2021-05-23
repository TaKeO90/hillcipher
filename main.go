/*
NOTE: for this implementation we only support matrices with size 2 and 3
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
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
	// TODO: need to separate the process of decryption for 2x2 and 3x3
	// TODO: Find Co-Factor.
	key := c.genKey
	InfoLogger.Printf("decrpytion key %v", key)
	det := key.Determinant(c.keysize)
	InfoLogger.Printf("Det %v", det)
	transpKey := key.Trans()
	InfoLogger.Printf("Trans key %v", transpKey)
	tK := Key(transpKey)
	InfoLogger.Printf("Minor %v", tK.GetKeyMinor())
	return
}

// NewCipher ....
func NewCipher(size int, word string) (hillCipher HillCipher) {

	key := new(Key)
	key.Gen(size)

	InfoLogger.Printf("Generated Key as Matrix %v", *key)
	InfoLogger.Printf("Generated Key %s", key.String())

	cipher := &Cipher{
		keysize: size,
		word:    word,
		genKey:  *key,
	}
	hillCipher = cipher

	return
}

// KeyAsString
type KeyAsString string

// ToMatrix ...
func (ks KeyAsString) ToMatrix(size int) *Key {
	key := Key{}
	index := 0
	row := []int{}
	count := 0
	for _, l := range ks {
		strings.Map(func(r rune) rune {
			if r == l {
				row = append(row, index)
			}
			index++
			return 0
		}, letters)
		index = 0
		count++
		if count == size {
			key = append(key, row)
			row = []int{}
			count = 0
		}
	}

	return &key
}

// Minor ...
type Minor [][]int

// Cal ...
func (m Minor) Cal() int {
	xIndex, yIndex := 0, len(m)-1
	return (m[xIndex][xIndex] * m[yIndex][yIndex]) - (m[xIndex][yIndex] * m[yIndex][xIndex])
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

func (k *Key) Trans() (Tkey [][]int) {

	Tkey = make([][]int, len(*k))
	for i := range *k {
		index := 0
		for _, n := range (*k)[i] {
			Tkey[index] = append(Tkey[index], n)
			index++
		}
	}
	return
}

func (k *Key) GetKeyMinor() (r [][]int) {
	for i := range *k {
		row := []int{}
		for j := range (*k)[i] {
			row = append(row, k.findMinor(i, j).Cal())
		}
		r = append(r, row)
	}
	return
}

func (k *Key) findMinor(x, y int) (m Minor) {
	for i := range *k {
		row := []int{}
		for j := range (*k)[i] {
			if i != x && j != y {
				row = append(row, (*k)[i][j])
			}
		}
		if len(row) != 0 {
			m = append(m, row)
		}
	}
	return
}

// String key into string.
func (k *Key) String() string {
	var res string
	for _, n := range *k {
		res += strings.Join(Stack(n).ToLetters(), "")
	}
	return res
}

func (k *Key) Determinant(size int) (det int) {
	switch size {
	case 2:
		det = k.detOf2x2()
	case 3:
		det = k.detOf3x3()
	}
	return
}

// get the determinant of matrices of size 3
func (k *Key) detOf3x3() (result int) {

	const firstRowIndex = 0

	stk := Stack{}
	for i := range *k {
		// twobytwo
		tbt := TwoByTwo{}
		x, y := 0, 0
		switch i {
		case 0:
			x = i + 1
			y = len(*k)
			tbt = append(tbt, (*k)[x][x:y]...)
			tbt = append(tbt, (*k)[x+1][x:y]...)
			stk = append(stk, tbt.Cal((*k)[firstRowIndex][i]))
		case 1:
			x = i
			y = len(*k) - 1
			tbt = append(tbt, (*k)[x][x-1], (*k)[x+1][x-1], (*k)[x][y], (*k)[x+1][y])
			stk = append(stk, tbt.Cal((*k)[firstRowIndex][i]))
		case 2:
			for _ = range *k {
				if x == len(*k)-1 {
					break
				}
				x++
				if x == 1 {
					tbt = append(tbt, (*k)[x][x-1:x+1]...)
				} else {
					tbt = append(tbt, (*k)[x][x-2:x]...)
				}
			}

			stk = append(stk, tbt.Cal((*k)[firstRowIndex][i]))
		}
	}

	d := stk.GetDet()
	result = MultModular(d, 26)
	return
}

// get the determinant of matrices of size 2
func (k *Key) detOf2x2() int {
	var tb TwoByTwo
	for _, n := range *k {
		for _, m := range n {
			tb = append(tb, m)
		}
	}
	return tb.Cal(-1)
}

// TwoByTwo ...
type TwoByTwo []int

func (tb TwoByTwo) Cal(currentNum int) (result int) {
	xIndex, yIndex := 0, len(tb)-1
	if currentNum != -1 {
		result = currentNum * ((tb[xIndex] * tb[yIndex]) - (tb[xIndex+1] * tb[yIndex-1]))
		return
	}
	result = (tb[xIndex] * tb[yIndex]) - (tb[xIndex+1] * tb[yIndex-1])
	return
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

	res := 0
	for i, n := range s {
		if 1&i == 0 {
			res += n
		} else {
			res -= n
		}
	}

	return s.Mod(res, 26)
}

func (s Stack) Mod(n, m int) int {
	if n < 0 && m > 0 || n > 0 && m < 0 {
		return (n % m) + m
	}

	return n % m
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

func MultModular(n, m int) int {
	for i := 0; i < m; i++ {
		if (i*n)%m == 1 {
			return i
		}
	}
	return 0
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

type CipherMode int

const (
	encryption CipherMode = iota
	decryption
)

func (cm CipherMode) String() string {
	switch cm {
	case encryption:
		return "encryption"
	case decryption:
		return "decryption"
	default:
		return ""
	}
}

func main() {

	var (
		mode    string
		word    string
		key     string
		keysize int
	)

	flag.StringVar(&mode, "mode", "encryption", "use 'encryption' to encrypt or 'decryption' to decrypt")
	flag.StringVar(&word, "word", "", "word to encrypt or decrypt")
	flag.StringVar(&key, "key", "", "key to use to decrypt")
	flag.IntVar(&keysize, "size", 0, "key size to generate when encrypting")

	flag.Parse()

	if mode == "" {
		flag.PrintDefaults()
		return
	}

	switch mode {
	case encryption.String():
		if keysize == 0 || word == "" {
			flag.PrintDefaults()
			return
		}
		fmt.Println(NewCipher(keysize, word).Encrypt())
	case decryption.String():
		if key == "" || word == "" || keysize == 0 {
			flag.PrintDefaults()
			return
		}

		mtxK := KeyAsString(key).ToMatrix(keysize)

		cipher := &Cipher{
			keysize: keysize,
			word:    word,
			genKey:  *mtxK,
		}

		fmt.Println(cipher.Decrypt())
	}

}
