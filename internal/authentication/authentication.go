package authentication

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sandbox/internal/database"
	"sandbox/internal/hasher"
	"sandbox/internal/helpers"
	"sandbox/internal/sessionmanager"
	"sandbox/internal/structs"
	"strings"
	"time"
)

func signUsingAuth(img, email, id, name, token string, w http.ResponseWriter) {
	structs.UserToken = token
	structs.IsAuth = true
	var userName, firstName, lastName, DatabaseEmail string
	var err error
	exist := false

	if email != "" {

		exist, err = database.CheckExistance("User", "email", email)
		if err != nil {
			http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
			return
		}
		fmt.Print(exist)

		if exist {

			finalUser, err := database.GetUserByEmail(email)
			if err != nil {
				return
			}

			err = sessionmanager.CreateSessionAndSetCookie(token, w, finalUser)
			if err != nil {
				http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
				return
			}
			return
		}

	} else {
		exist, err = database.CheckExistance("User", "email", strings.ReplaceAll(name, " ", "")+"@sandbox.com")

		if err != nil {
			http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
			return
		}

		if exist {

			finalUser, err := database.GetUserByUsername(strings.ReplaceAll(name, " ", ""))
			if err != nil {
				return
			}

			err = sessionmanager.CreateSessionAndSetCookie(token, w, finalUser)
			if err != nil {
				http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	if userName != "" {

		// check if the username exists
		exist, err := database.CheckExistance("User", "username", strings.ReplaceAll(name, " ", ""))
		if err != nil {
			http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
			return
		}
		if exist {

			finalUser, err := database.GetUserByUsername(strings.ReplaceAll(name, " ", ""))
			if err != nil {
				return
			}

			err = sessionmanager.CreateSessionAndSetCookie(token, w, finalUser)
			if err != nil {
				http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
				return
			} 
			return

		} else {
			if len(id) >= 3 {
				userName = userName + id[0:3]

			} else {
				userName = userName + id

			}
			exist, err = database.CheckExistance("User", "username", userName)
			if err != nil {
				http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
				return
			}

			if exist {

				finalUser, err := database.GetUserByUsername(userName)
				if err != nil {
					return
				}

				err = sessionmanager.CreateSessionAndSetCookie(token, w, finalUser)
				if err != nil {
					http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
					return
				}
				return

			}
		}
	}

	if !exist {
		imageID := 0
		var err error
		if img != "" {
			IsImg, _ := helpers.IsDataImage([]byte(img))
			if IsImg {
				imageID, err = database.UploadImage([]byte(img))
				if err != nil {
					log.Printf("SignupHandler: %s\n", err.Error())
				}
			}
		}
		// structure
		hashedPassword, hashErr := hasher.GetHash("123123123")
		if hashErr != hasher.HasherErrorNone {
			http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
			return
		}

		if name != "" {
			names := strings.Split(name, " ")
			if len(names) > 1 {
				firstName = names[0]
				lastName = names[len(names)-1]
			}

			userName = strings.ReplaceAll(name, " ", "")
			if email == "" {
				DatabaseEmail = userName + "@sandbox.com"

				exist, err := database.CheckExistance("User", "email", DatabaseEmail)
				if err != nil {
					http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
					return
				}

				if exist {
					if len(id) >= 3 {
						DatabaseEmail = userName + id[0:3] + "@sandbox.com"

					} else {
						DatabaseEmail = userName + id + "@sandbox.com"

					}
				}
			} else {
				DatabaseEmail = email

				exist, err := database.CheckExistance("User", "email", DatabaseEmail)
				if err != nil {
					http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
					return
				}

				if exist {
					if len(id) >= 3 {
						DatabaseEmail = DatabaseEmail + id[0:3]

					} else {
						DatabaseEmail = DatabaseEmail + id

					}
				}
			}

			// check if the username exists
			exist, err := database.CheckExistance("User", "username", userName)
			if err != nil {
				http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
				return
			}

			if exist {
				if len(id) >= 3 {
					userName = userName + id[0:3]

				} else {
					userName = userName + id

				}
			}

		}

		if email != "" {
			fmt.Print(email, "hi")
			if name == "" {
				parts := strings.Split(email, "@")
				firstName = parts[0]
				lastName = parts[0]
				userName = parts[0]
				// check if the username exists
				exist, err := database.CheckExistance("User", "username", userName)
				if err != nil {
					http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
					return
				}

				if exist {
					if len(id) >= 3 {
						userName = userName + id[0:3]

					} else {
						userName = userName + id

					}
				}
				DatabaseEmail = email

			}
		}
		fmt.Print(DatabaseEmail)

		cleanedUserData := structs.User{
			Username:       userName,
			Email:          DatabaseEmail,
			FirstName:      firstName,
			LastName:       lastName,
			DateOfBirth:    time.Time{},
			HashedPassword: hashedPassword,
			ImageId:        imageID,
			GithubName:     "",
			LinkedinName:   "",
			TwitterName:    "",
		}
		err = database.CreateUser(cleanedUserData)
		if err != nil {
			http.Error(w, "could not create a user, please try again later", http.StatusBadRequest)
			return
		}

		finalUser, err := database.GetUserByUsername(userName)
		if err != nil {
			return
		}

		err = sessionmanager.CreateSessionAndSetCookie(token, w, finalUser)
		if err != nil {
			http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
			return
		}
	}
}

func loadEnvVariables(APiKEY string) (string, error) {
	envVariables := make(map[string]string)

	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			envVariables[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return envVariables[APiKEY], nil
}
