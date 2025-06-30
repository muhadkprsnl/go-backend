// // scripts/hash.go
// package main

// import (
// 	"fmt"

// 	"golang.org/x/crypto/bcrypt"
// )

// func main() {
// 	password := "admin@Ma" // Change this to your desired password

// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		panic(err)
// 	}

//		fmt.Println("Hashed password:")
//		fmt.Println(string(hash))
//	}
package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin@Mail" // exact casing
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Println(string(hash))
}
