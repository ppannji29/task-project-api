package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"task-project/models"
	"task-project/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var jwtIssuer = os.Getenv("JWT_ISSUER")

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email string `json:"email"`
}

type VerifyOtpRequest struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

type UserAuth struct {
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type TokenResponse struct {
	Message  string   `json:"message"`
	Token    string   `json:"token"`
	UserAuth UserAuth `json:"user_auth"`
}

func RequestOtp(userCol *mongo.Collection, otpCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		ctx := context.Background()
		var user models.User
		err := userCol.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		code := utils.GenerateOtp(6)

		otp := models.Otp{
			UserEmail: req.Email,
			OtpCode:   code,
			OtpExpiry: time.Now().Add(2 * time.Minute),
			IsClaimed: false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, _ = otpCol.InsertOne(ctx, otp)

		errEmail := utils.SendEmail(req.Email, "Your OTP Code", code)
		if errEmail != nil {
			http.Error(w, "Failed to send OTP email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"otp": code})
	}
}

func VerifyOtp(userCol *mongo.Collection, otpCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req VerifyOtpRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		ctx := context.Background()
		var otp models.Otp
		err := otpCol.FindOne(ctx, bson.M{
			"user_email": req.Email,
			"otp_code":   req.Otp,
			"is_claimed": false,
		}).Decode(&otp)
		if err != nil || time.Now().After(otp.OtpExpiry) {
			http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
			return
		}

		var user models.User
		err = userCol.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		accessToken, _ := generateToken(user.UserID, user.Email, "access", 1*time.Hour)
		refreshToken, _ := generateToken(user.UserID, user.Email, "refresh", 24*time.Hour)

		setTokenCookie(w, "access_token", accessToken, 1*time.Hour)
		setTokenCookie(w, "refresh_token", refreshToken, 24*time.Hour)

		_, _ = otpCol.UpdateOne(ctx, bson.M{"_id": otp.ID}, bson.M{"$set": bson.M{"is_claimed": true}})

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TokenResponse{
			Message:  "Login successful",
			Token:    accessToken,
			UserAuth: UserAuth{UserID: user.UserID, Email: user.Email, Name: user.Name},
		})
	}
}

func RefreshToken(userCol *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "Refresh token not found", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid || claims.Type != "refresh" {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		objID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		ctx := context.Background()
		var user models.User
		err = userCol.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		newAccessToken, err := generateToken(user.UserID, user.Email, "access", 1*time.Hour)
		if err != nil {
			http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
			return
		}

		setTokenCookie(w, "access_token", newAccessToken, 1*time.Hour)

		json.NewEncoder(w).Encode(TokenResponse{
			Message:  "Token refreshed successfully",
			Token:    newAccessToken,
			UserAuth: UserAuth{UserID: user.UserID, Email: user.Email, Name: user.Name},
		})
	}
}

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clearTokenCookie(w, "access_token")
		clearTokenCookie(w, "refresh_token")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TokenResponse{Message: "Logout successful"})
	}
}

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				http.Error(w, "Access token not found", http.StatusUnauthorized)
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil || !token.Valid || claims.Type != "access" {
				http.Error(w, "Invalid access token", http.StatusUnauthorized)
				return
			}

			if _, err := primitive.ObjectIDFromHex(claims.UserID); err != nil {
				http.Error(w, "Invalid user ID format", http.StatusUnauthorized)
				return
			}

			r.Header.Set("X-User-ID", claims.UserID)
			r.Header.Set("X-User-Email", claims.Email)

			next.ServeHTTP(w, r)
		})
	}
}

func generateToken(userID, email, tokenType string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtIssuer,
			Subject:   userID,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func setTokenCookie(w http.ResponseWriter, name, value string, duration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(duration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func clearTokenCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
