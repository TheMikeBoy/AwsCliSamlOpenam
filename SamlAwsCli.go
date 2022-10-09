package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/term"
)

//structs

type Awscred struct {
	Credentials struct {
		AccessKeyId     string `json:"AccessKeyId"`
		SecretAccessKey string `json:"SecretAccessKey"`
		SessionToken    string `json:"SessionToken"`
	} `json:"Credentials"`
}

// Acessa o OpenAm e recupera o token
func ReturnOpenAmToken(url string, user string, pass string) string {

	body, _ := json.Marshal(map[string]string{
		"uri": "ldapService",
	})
	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		log.Fatalf("Error ao gerar a requesicao do openam %s", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("X-OpenAM-Username", user)
	req.Header.Add("X-OpenAM-Password", pass)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	bodyf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	s := strings.Split(string(bodyf), "\"")
	Iplanetpro := s[3]

	return Iplanetpro
}

// Funcao para usuario e senha
func credentials() (string, string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Digite Seu usuario ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}

	fmt.Print("Digite Seu Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}

	fmt.Print("Digite Seu url ")
	url, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}
	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), strings.TrimSpace(url), nil
}

func main() {
	username, password, url, _ := credentials()
	Iplanetpro := ReturnOpenAmToken(url, username, password)
	println(Iplanetpro)
	os := runtime.GOOS

	println(os)
}
