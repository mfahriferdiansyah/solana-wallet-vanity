package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mr-tron/base58"
)

func worker(prefix string, maxAttempts int, result chan<- string, privateKeyChan chan<- string, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	for attempts := 0; maxAttempts == 0 || attempts < maxAttempts; attempts++ {
		publicKey, privateKey := generateKeypair()

		if strings.HasPrefix(publicKey, prefix) {
			result <- publicKey
			privateKeyChan <- privateKey
			return
		}

		if attempts%100000 == 0 && id == 0 { 
			fmt.Printf("[Worker %d] Attempts: %d\n", id, attempts)
		}
	}
}

func generateKeypair() (string, string) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println("Error generating keypair:", err)
		os.Exit(1)
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)
	return base58.Encode(publicKey), base58.Encode(privateKey)
}

func saveResultToFile(publicKey, privateKey string) {
	file, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening result file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Public Key: %s\nPrivate Key: %s\n\n", publicKey, privateKey))
	if err != nil {
		fmt.Println("Error writing to result file:", err)
	}
}

func main() {
	prefix := flag.String("prefix", "", "The prefix to match (required)")
	maxAttempts := flag.Int("attempts", 0, "Maximum number of attempts per worker (0 for unlimited)")
	workers := flag.Int("workers", 4, "Number of parallel workers")
	flag.Parse()

	if *prefix == "" {
		fmt.Println("Error: --prefix is required")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Searching for wallets starting with: %s\nUsing %d workers...\n", *prefix, *workers)

	result := make(chan string, 1)
	privateKeyChan := make(chan string, 1)

	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(*prefix, *maxAttempts, result, privateKeyChan, &wg, i)
	}

	go func() {
		wg.Wait()
		close(result)
		close(privateKeyChan)
	}()

	select {
	case publicKey := <-result:
		privateKey := <-privateKeyChan
		fmt.Printf("Match found!\n")
		fmt.Printf("Public Key: %s\n", publicKey)
		fmt.Printf("Private Key: %s\n", privateKey)
		fmt.Printf("Elapsed Time: %s\n", time.Since(start))

		saveResultToFile(publicKey, privateKey)
	case <-time.After(time.Hour * 24): 
		fmt.Println("Timeout reached. No match found.")
	}
}

