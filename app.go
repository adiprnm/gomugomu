package main

import (
    "os"
    "os/exec"
    "io/ioutil"
    "encoding/json"
    "fmt"
    "strings"
)

func main() {
    command := os.Args[1]

    // The logic
    if command == "init" {
        var APP_NAME, APP_PORT string
        fmt.Printf("Your app name: ");
        fmt.Scanf("%s", &APP_NAME)
        fmt.Printf("Your app port: ");
        fmt.Scanf("%s", &APP_PORT)

        file, err := os.Create("config.json")
        if err != nil {
            fmt.Println("[ERROR] Error creating config.json")
        } else {
            file.WriteString("{\n")
            file.WriteString("    \"APP_NAME\": \"" + APP_NAME + "\",\n")
            file.WriteString("    \"APP_PORT\": \"" + APP_PORT + "\",\n")
            file.WriteString("}")
            fmt.Println("[INFO] File config.json has been successfully created!")
        }
    } else {
        // Open config.json
        jsonFile, err := os.Open("config.json")
        defer jsonFile.Close()
        if err != nil {
            panic(err)
        }

        // Change config.json to struct
        byteVal, _ := ioutil.ReadAll(jsonFile)
        var config map[string]interface{}
        json.Unmarshal([]byte(byteVal), &config)

        if command == "serve" {
            fmt.Printf("[INFO] %s is starting on port %s\n\n", config["APP_NAME"], config["APP_PORT"])
            exec.Command(fmt.Sprintf("%s", config["APP_NAME"])).Run()
        } else if command == "install" {
            ctrlFiles, err := ioutil.ReadDir("controllers")
            if err != nil {
                fmt.Println("[ERROR] Directory not found!")
            } else {
                for _, v := range(ctrlFiles) {
                    var command = fmt.Sprintf("go install %s/controllers/%s", config["APP_NAME"], v.Name())
                    fmt.Println("[INFO] Running command " + command)
                    exec.Command(command).Run()
                }
            }
        } else if strings.Split(command, ":")[0] == "create" {
            module := strings.Split(command, ":")[1]

            if module == "controller" {
                if len(os.Args) < 2 {
                    fmt.Println("[ERROR] Controller name not provided!")
                    os.Exit(1)
                } else {
                    controllerName := os.Args[2]

                    if _, err := os.Stat("controllers/" + controllerName); !os.IsNotExist(err) {
                        fmt.Println("[ERROR] Controller " + controllerName + " already exist!")
                        os.Exit(1)
                    } else {
                        os.MkdirAll("controllers/" + controllerName, os.ModePerm)
                    }

                    f, err := os.Create("controllers/" + controllerName + "/" + controllerName + ".go")
                    if err != nil {
                        fmt.Println("[ERROR] Error while opening controllers/" + controllerName + "/" + controllerName + ".go")
                    } else {
                            f.WriteString("package " + controllerName)
                            f.WriteString("\n\n")
                            f.WriteString("import (\n")
                            f.WriteString("    \"net/http\"")
                            f.WriteString("\n)")
                            f.WriteString("\n\n")
                            f.WriteString("// Write your code here...")
                            fmt.Printf("[INFO] Creating %s/controllers/%s/%s.go ...\n",
                                config["APP_NAME"],
                                controllerName,
                                controllerName)

                    }

                }
            }
        }
    }
}
