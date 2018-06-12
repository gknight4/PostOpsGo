package main 

import (
	"golang.org/x/crypto/bcrypt"
//	"crypto/sha256"
	"fmt"
//	"encoding/base64"
	"net/http"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
//	"log"
)

const (
	SIGNKEY = "thumbsup"
)

type ClaimsJWT struct {
	id int
	expire int
	flags int
}

/*
2018/05/06 13:49:23 map[expire:2018-05-06T13:36:36.750487137-07:00 flags:1 id:5aef4ae301914f5d3985f57f]
*/

/*func demoToken(){
// sample token string taken from the New example
tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

// Parse takes the token string and a function for looking up the key. The latter is especially
// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
// head of the token to identify which key to use, but the parsed token (head and claims) is provided
// to the callback, providing flexibility.
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    // Don't forget to validate the alg is what you expect:
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
    }

    // hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
    return []byte("my_secret_key"), nil
})

if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    fmt.Println(claims["foo"], claims["nbf"])
} else {
    fmt.Println(err)
}	
}*/

func ValidateToken(tokenString string) (interface{}, bool) {
	token, _ := jwt.Parse(tokenString, 
		func(token *jwt.Token) (interface{}, error) {
	    // Don't forget to validate the alg is what you expect:
//	    lo("val2")
	    if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
		    return []byte(SIGNKEY), nil
	    } else {
	        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	    }
	})
	if token.Valid {
		return token.Claims.(jwt.MapClaims), token.Valid
	} else {
		return nil, false
	}
}

func GetJwtToken(user User, flags int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	    "id": user.ID,
	    "expire": time.Now(),
	    "flags": flags,
	})
//	h := sha256.New()
//	h.Write([]byte(SIGNKEY))
	tokenString, _ := token.SignedString([]byte(SIGNKEY)) // h.Sum(nil)
//	log.Println("token: " + tokenString)
//	if claims, ok := ValidateToken(tokenString); ok {
//		log.Println("claims:")
//		log.Println(claims["id"])
//	}
	return tokenString
}


func ReturnJson (w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Msr(r string) map[string]string   {
// make success / failure map return
/*	var resp bool
	if r {
		resp = true} else {
		resp = false
	}*/
		return map[string]string{"result": r}
//	returnJson(w, http.StatusOK, )
	
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func bToOk (r bool) string {
    var res string
	if r {
            res = "ok"
        } else {
            res = "err"
        }
        return res
}


func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

