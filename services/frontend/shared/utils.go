// Copyright 2023 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, code int, context string, err error) bool {
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error(), "context": context})
		log.Printf("error occured @ %s: %s\n", context, err.Error())
		return true
	}
	return false
}

func ValidateEnvVars() {
	_, present := os.LookupEnv("CLIENT_ID")
	if !present {
		log.Fatal("CLIENT_ID is not present")
	}
	_, present = os.LookupEnv("CLIENT_SECRET")
	if !present {
		log.Fatal("CLIENT_SECRET is not present")
	}
	_, present = os.LookupEnv("LISTEN_PORT")
	if !present {
		log.Fatal("LISTEN_PORT is not present")
	}
	_, present = os.LookupEnv("CLIENT_LAUNCHER_PORT")
	if !present {
		log.Fatal("CLIENT_LAUNCHER_PORT is not present")
	}
	_, present = os.LookupEnv("PROFILE_SERVICE")
	if !present {
		log.Fatal("PROFILE_SERVICE is not present")
	}
	_, present = os.LookupEnv("PING_SERVICE")
	if !present {
		log.Fatal("PING_SERVICE is not present")
	}
	_, present = os.LookupEnv("JWT_KEY")
	if !present {
		log.Fatal("JWT_KEY is not present")
	}
}
