package main

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  string
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://ac88-46-2-221-175.ngrok-free.app/auth/google/callback",
		ClientID:     "476930073014-18kmcf0bv2ltthokm1tt4gapg64jrkbr.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	oauthStateString = generateStateString()
	log.Println("Google OAuth configuration initialized")
}

func generateStateString() string {
	b := make([]byte, 32)
	cryptorand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	log.Printf("Redirecting to Google OAuth URL: %s", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Google OAuth callback")

	state := r.FormValue("state")
	code := r.FormValue("code")

	if state != oauthStateString {
		log.Printf("Invalid OAuth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	content, err := getUserInfo(state, code)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		http.Error(w, "Error getting user info", http.StatusInternalServerError)
		return
	}

	var userInfo struct {
		Email string `json:"email"`
	}
	err = json.Unmarshal(content, &userInfo)
	if err != nil {
		log.Printf("Error unmarshaling user info: %v", err)
		http.Error(w, "Error processing user information", http.StatusInternalServerError)
		return
	}

	if userInfo.Email == "" {
		log.Println("Error: Email is empty")
		http.Error(w, "Invalid email received from Google", http.StatusBadRequest)
		return
	}

	log.Printf("Received user info for email: %s", userInfo.Email)

	membershipID := generateMembershipID()

	// Check if the user already exists
	existingUser, err := getUser(userInfo.Email)
	if err == nil {
		log.Printf("User already exists with email: %s", userInfo.Email)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message":       "User already exists",
			"membership_id": existingUser.MembershipID,
			"email":         userInfo.Email,
		})
		return
	}

	// Generate a 6-digit numerical password
	password := generateSixDigitPassword()
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// User doesn't exist, create a new one
	err = createUser(membershipID, userInfo.Email, hashedPassword)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	log.Printf("User created successfully via Google OAuth: %s", userInfo.Email)

	log.Printf("User created successfully via Google OAuth: %s", userInfo.Email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":       "User created successfully via Google OAuth",
		"membership_id": membershipID,
		"email":         userInfo.Email,
	})
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func generateSixDigitPassword() string {
	return fmt.Sprintf("%06d", randGen.Intn(1000000))
}

func generateRandomPassword() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
