package auth

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func CheckForAuth() (string, string, string) {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("No .env file found, generating a new one to store your user data")
	}

	//Check for the App ID
	appId, exists := os.LookupEnv("APP_ID")
	if !exists {
		fmt.Println("APP_ID not found, Please neter a new APP_ID : ")

		reader := bufio.NewReader(os.Stdin)
		auth, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		appId = strings.TrimSpace(auth)

		fmt.Printf("New App ID set %s\n", appId)

		//save the new AUTH_KEY to .env file
		saveAuthKeyToFile("APP_ID", appId)

	} else {
		fmt.Printf("APP_ID found\n")
	}

	//Check for the OSID
	osID, exists := os.LookupEnv("OS_ID")
	if !exists {
		fmt.Println("OS_ID not found, Please neter a new OneSignal ID : ")

		reader := bufio.NewReader(os.Stdin)
		auth, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		osID = strings.TrimSpace(auth)

		fmt.Printf("New OS_ID set %s\n", osID)

		//save the new AUTH_KEY to .env file
		saveAuthKeyToFile("OS_ID", osID)

	} else {
		fmt.Printf("OS_ID found\n")
	}

	//Check for the User Auth Key
	authKey, exists := os.LookupEnv("AUTH_KEY")
	if !exists {
		fmt.Println("AUTH_KEY not found, Please neter a new AUTH_KEY : ")

		reader := bufio.NewReader(os.Stdin)
		auth, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		authKey = strings.TrimSpace(auth)

		fmt.Printf("New AUTH_KEY set %s\n", authKey)

		//save the new AUTH_KEY to .env file
		saveAuthKeyToFile("AUTH_KEY", authKey)

	} else {
		fmt.Printf("AUTH_KEY found\n")
	}

	return appId, authKey, osID
}

func saveAuthKeyToFile(key, value string) {

	file, err := os.OpenFile(".env", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	fmt.Println("Saved to .env file successfully")

}
