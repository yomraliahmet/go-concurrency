package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func main() {
	var users []User
	usersCh := make(chan User)

	var wg1 sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg1.Add(1)
		go func(ch chan User) {
			getRandomUser(ch)
			wg1.Done()
		}(usersCh)

	}

	go func() {
		var a int = 1
		for {
			select {
			case user := <-usersCh:
				user.Sequence = a
				users = append(users, user)
				fmt.Println(user)
			}
			a++
		}
	}()
	wg1.Wait()
	fmt.Println("-------------------------------------------")
	fmt.Println(users)
}

func getRandomUser(ch chan User) {

	response, _ := http.Get("https://random-data-api.com/api/users/random_user")

	body, _ := ioutil.ReadAll(response.Body)

	var user User

	json.Unmarshal(body, &user)

	ch <- user
}

type User struct {
	Sequence  int    `json:"sequence"`
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
