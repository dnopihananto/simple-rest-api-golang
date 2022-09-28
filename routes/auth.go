package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danilopolani/gocialite/structs"
	"github.com/dgrijalva/jwt-go"
	"github.com/dnopihananto/gin-full-api/config"
	"github.com/dnopihananto/gin-full-api/models"
	"github.com/gin-gonic/gin"
)

// Redirect to correct oAuth URL
func RedirectHandler(c *gin.Context) {
	// Retrieve provider from route
	provider := c.Param("provider")

	// In this case we use a map to store our secrets, but you can use dotenv or your framework configuration
	// for example, in revel you could use revel.Config.StringDefault(provider + "_clientID", "") etc.
	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GH"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GH"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
		"google": {
			"clientID":     os.Getenv("CLIENT_ID_G"),
			"clientSecret": os.Getenv("CLIENT_SECRET_G"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/google/callback",
		},
	}

	providerScopes := map[string][]string{
		"github": []string{"public_repo"},
		"google": []string{},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

// Handle callback of provider
func CallbackHandler(c *gin.Context) {
	// Retrieve query params for state and code
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	// Handle callback and check for errors
	user, _, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	var newUser = getOrRegisterUser(provider, user)
	var jwtToken = createToken(&newUser)

	c.JSON(200, gin.H{
		"data":    newUser,
		"token":   jwtToken,
		"message": "berhasil login",
	})
}

func getOrRegisterUser(provider string, user *structs.User) models.User {
	var userData models.User
	config.DB.Where("provider = ? AND social_id = ?", provider, user.ID).First(&userData)

	if userData.ID == 0 {
		newUser := models.User{
			FullName: user.FullName,
			Email:    user.Email,
			SocialId: user.ID,
			Provider: provider,
			Avatar:   user.Avatar,
		}
		config.DB.Create(&newUser)
		return newUser
	} else {
		return userData
	}
}

func createToken(user *models.User) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"user_role": user.Role,
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
		"iat":       time.Now().Unix(),
	})
	fmt.Println(jwtToken)

	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	fmt.Println(tokenString, err)

	return tokenString
}

// temporary check token
func CheckToken(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "berhasil login",
	})
}

func GetProfile(c *gin.Context) {
	var user models.User
	user_id := uint(c.MustGet("jwt_user_id").(float64))

	config.DB.Where("id = ?", user_id).Preload("Article", "user_id = ?", user_id).Find(&user)

	c.JSON(200, gin.H{
		"message": "berhasil ke halaman profile",
		"data":    user,
	})
}
