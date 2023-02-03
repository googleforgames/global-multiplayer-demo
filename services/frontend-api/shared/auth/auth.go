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

package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(id string, days int) (string, error) {
	expirationTime := time.Now().Add(24 * 31 * time.Hour)

	// Create the JWT claims, which includes google's profile id and expiry time
	claims := &Claims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(endpointHandler func(id string, rw http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		prefix := "Bearer "
		authHeader := req.Header.Get("Authorization")
		reqToken := strings.TrimPrefix(authHeader, prefix)

		if len(reqToken) != 0 {
			claims := &Claims{}

			// Parse the JWT string and store the result in `claims`.
			// Note that we are passing the key in this method as well. This method will return an error
			// if the token is invalid (if it has expired according to the expiry time we set on sign in),
			// or if the signature does not match
			tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_KEY")), nil
			})
			if err != nil {
				fmt.Println(err.Error())
				if err == jwt.ErrSignatureInvalid {
					rw.WriteHeader(http.StatusUnauthorized)
					return
				}
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			if !tkn.Valid {
				fmt.Println("Invalid token")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			endpointHandler(claims.Id, rw, req)
		} else {
			fmt.Println("Authorization token is not present")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}
