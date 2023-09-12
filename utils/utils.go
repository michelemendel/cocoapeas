package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"
)

func HandleInterrupts(postOps ...func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("\ninterrupted, running postOps...")
	for _, op := range postOps {
		op()
	}
}

var once sync.Once

func InitEnv() {
	once.Do(InitEnvExec)
}

const projectDirName = "cocopeas"

func InitEnvExec() {
	// fmt.Println("Initializing environment...")
	rootDir := GetRootDir()
	envFile := filepath.Join(string(rootDir), ".env")
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("error loading env file:", envFile)
	}
}

// We have to dynamically find the project root directory, because
// it's different for tests and the main server.
func GetRootDir() string {
	projectName := regexp.MustCompile(`^.*` + projectDirName + ``)
	cwd, _ := os.Getwd()
	rootDir := projectName.Find([]byte(cwd))
	return string(rootDir)
}

type UUID string

func GenerateUUID() UUID {
	return UUID(ksuid.New().String())
}

func PP(s any) {
	res, err := PrettyStruct(s)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(res)
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
