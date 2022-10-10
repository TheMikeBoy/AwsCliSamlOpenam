package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"

	"github.com/abdfnx/gosh"
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

type UrlsFiles struct {
	OpenAmUrl string
	IdpUrl    string
}

// Acessa o OpenAm e recupera o token
func ReturnOpenAmToken(url string, user string, pass string) string {

	body, _ := json.Marshal(map[string]string{
		"uri": "ldapService",
	})
	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		log.Fatalf("Error %s", err)
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
func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

func checkFileOS(OSystem string, homeDir string) {
	//	var URL string
	urlurlfile := (homeDir + "/AwsCliSaml/URLS.json")
	if _, err := os.Stat(urlurlfile); err == nil {
		fmt.Println("File Exist")
		return

	} else {
		createConfFileOS(OSystem, homeDir)

		return
	}

}

func createConfFileOS(OSystem string, homeDir string) {
	var OpenAmUrl string
	var IdpUrl string
	if _, err := os.Stat(homeDir + "/AwsCliSaml/"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(homeDir+"/AwsCliSaml/", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Printf("Autentication OPENAM URL: ")
	fmt.Scan(&OpenAmUrl)
	fmt.Printf("Autentication IDP URL: ")
	fmt.Scan(&IdpUrl)

	UrlsStrings := []byte("{ \"OpenAmUrl\":" + "\"" + string(OpenAmUrl) + "\"," + "\"IDPURL\":" + "\"" + string(IdpUrl) + "\"}")
	aux := UrlsFiles{}
	err := json.Unmarshal(UrlsStrings, &aux)
	if err != nil {
		// nozzle.printError("opening config file", err.Error())
	}
	UrlString, _ := json.Marshal(&aux)
	err1 := os.WriteFile(homeDir+"/AwsCliSaml/URLS.json", UrlString, 0755)
	if err1 != nil {
		fmt.Printf("Unable to write file: %v", err1)
	} else {
		println("File Create")
	}
	return
}
func returnURls(homeDir string) (string, string) {

	urls := UrlsFiles{}

	b, err := os.ReadFile(homeDir + "/AwsCliSaml/URLS.json")
	if err != nil {
		fmt.Print(err)
	}

	json.Unmarshal([]byte(b), &urls)

	return urls.OpenAmUrl, urls.IdpUrl
}

func IdpRequest(url string, cookie string) ([]byte, string) {
	body, _ := json.Marshal(map[string]string{
		"": "",
	})
	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		log.Fatalln(err)
	}
	req.AddCookie(&http.Cookie{Name: "iPlanetDirectoryPro", Value: cookie})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	BodyDecode := html.UnescapeString(string(respbody))
	//println(BodyDecode)
	//println(teste[strings.Index(teste):])
	re := returnstring(BodyDecode, "name=\"SAMLResponse\" value=\"", "\" />")
	samlToken := strings.ReplaceAll(re, "\r\n", "")
	data, err := base64.StdEncoding.DecodeString(samlToken)
	if err != nil {
		log.Fatal("error:", err)
	}
	return data, samlToken
}

func returnstring(body string, textone string, texttwo string) string {

	pos1 := strings.Index(body, textone)
	pos2 := strings.Index(body, texttwo)
	partstring := body[pos1+len(textone) : pos2]

	return partstring
}
func returnrole(data []byte) (string, string) {
	Rolesmap := make(map[int]string)
	i := 0
	var rolechoose int
	var validID = regexp.MustCompile(`arn[:A-Za-z\d/,-]*`)
	roles := validID.FindAll(data, -1)
	for _, role := range roles {

		Rolesmap[i] = string(role)
		i++
	}

	Text3 := "::"
	Text4 := ":role"
	Text5 := "role/"
	Text6 := ","

	for k := range Rolesmap {
		pos3 := strings.Index(Rolesmap[k], Text3)
		pos4 := strings.Index(Rolesmap[k], Text4)
		account := Rolesmap[k][pos3+len(Text3) : pos4]

		pos5 := strings.Index(Rolesmap[k], Text5)
		pos6 := strings.Index(Rolesmap[k], Text6)
		Roleaccount := Rolesmap[k][pos5+len(Text5) : pos6]

		fmt.Println(k, " Conta: ", account, "Role: ", Roleaccount)

	}
	println("Set you Role: ")
	fmt.Scan(&rolechoose)
	//rolechoose = 5

	return strings.Split(Rolesmap[rolechoose], ",")[0], strings.Split(Rolesmap[rolechoose], ",")[1]

}

const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}
func main() {
	username, password, _ := credentials()
	println("\n")
	operationSystem := runtime.GOOS
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	checkFileOS(operationSystem, dirname)
	OpenAmUrl, IdpUrl := returnURls(dirname)
	Iplanetpro := ReturnOpenAmToken(OpenAmUrl, username, password)
	data, SAMLResponse := IdpRequest(IdpUrl, Iplanetpro)
	RoleArn, PrincipalArn := returnrole(data)

	AssumeRole := ("aws sts assume-role-with-saml --role-arn " + RoleArn + " --principal-arn " + PrincipalArn + " --saml-assertion " + SAMLResponse)
	if operationSystem == "linux" {
		err, cmdreturn, errout := Shellout(AssumeRole)
		if err != nil {
			log.Printf("error: %v\n %s", err, string(errout))

		}
		//println(AssumeRole)
		//Create Json with  return from aws cli
		var awscred Awscred
		json.Unmarshal([]byte(cmdreturn), &awscred)
		//println(fmt.Println(string(errout)))

		exportCreds := []byte("[default]\naws_access_key_id = " + awscred.Credentials.AccessKeyId + "\naws_secret_access_key = " + awscred.Credentials.SecretAccessKey + "\naws_session_token =" + awscred.Credentials.SessionToken)
		err1 := os.WriteFile(dirname+"/.aws/credentials", exportCreds, 0755)
		if err1 != nil {
			fmt.Printf("Unable to write file: %v", err1)
		}
		fmt.Println("File Created")
	} else if operationSystem == "windows" {

		err, out, errout := gosh.PowershellOutput(`AssumeRole`)

		if err != nil {
			log.Printf("error: %v\n", err)
			fmt.Print(errout)
		}
		var awscred Awscred
		json.Unmarshal([]byte(out), &awscred)
		exportCreds := []byte("[default]\naws_access_key_id = " + awscred.Credentials.AccessKeyId + "\naws_secret_access_key = " + awscred.Credentials.SecretAccessKey + "\naws_session_token =" + awscred.Credentials.SessionToken)
		err1 := os.WriteFile(dirname+"/.aws/credentials", exportCreds, 0755)
		if err1 != nil {
			fmt.Printf("Unable to write file: %v", err1)
		}
		fmt.Println("File Created")

	} else {
		fmt.Println("Operation System Not Supported")
	}

}
