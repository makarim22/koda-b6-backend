package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func main() {
	// order_id=2
	// status_code=200
	// gross_amount=50000.00
	h := sha512.New()
	h.Write([]byte("2" + "200" + "50000.00" + "test-server-key"))
	fmt.Println(hex.EncodeToString(h.Sum(nil)))
}
