package auth

import (
    "testing"
    "time"
    "net/http"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
    userID := uuid.New()
    tokenSecret := "test_secret"
    expiresIn := time.Hour

    tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
    if err != nil {
        t.Fatalf("MakeJWT returned an error: %v", err)
    }

    // Парсим токен без проверки подписи для извлечения claims
    token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &jwt.RegisteredClaims{})
    if err != nil {
        t.Fatalf("Failed to parse token: %v", err)
    }

    claims, ok := token.Claims.(*jwt.RegisteredClaims)
    if !ok {
        t.Fatalf("Expected jwt.RegisteredClaims, got %T", token.Claims)
    }

    if claims.Subject != userID.String() {
        t.Errorf("Expected Subject to be %v, got %v", userID.String(), claims.Subject)
    }

    if claims.Issuer != "chirpy" {
        t.Errorf("Expected Issuer to be 'chirpy', got %v", claims.Issuer)
    }

    if !claims.ExpiresAt.Time.After(time.Now()) {
        t.Errorf("Expected ExpiresAt to be in the future, got %v", claims.ExpiresAt.Time)
    }
}

func TestValidateJWT(t *testing.T) {
    userID := uuid.New()
    tokenSecret := "test_secret"
    expiresIn := time.Hour

    tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
    if err != nil {
        t.Fatalf("MakeJWT returned an error: %v", err)
    }

    validatedUserID, err := ValidateJWT(tokenString, tokenSecret)
    if err != nil {
        t.Fatalf("ValidateJWT returned an error: %v", err)
    }

    if validatedUserID != userID {
        t.Errorf("Expected userID to be %v, got %v", userID, validatedUserID)
    }
}

func TestValidateJWT_InvalidToken(t *testing.T) {
    tokenString := "invalid_token"
    tokenSecret := "test_secret"

    _, err := ValidateJWT(tokenString, tokenSecret)
    if err == nil {
        t.Fatal("Expected an error for invalid token, but got nil")
    }
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
    userID := uuid.New()
    tokenSecret := "test_secret"
    expiresIn := -time.Hour // Токен уже истек

    tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
    if err != nil {
        t.Fatalf("MakeJWT returned an error: %v", err)
    }

    _, err = ValidateJWT(tokenString, tokenSecret)
    if err == nil {
        t.Fatal("Expected an error for expired token, but got nil")
    }
}

func TestGetBearerToken(t *testing.T) {
    headers := http.Header{}
    tokenString := "Bearer TOKEN_STRING"
    headers.Add("Authorization", tokenString)
    token, _ := GetBearerToken(headers)
    if token != "TOKEN_STRING" {
        t.Fatal("Error JWT token")
    }
}
