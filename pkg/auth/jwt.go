package auth

import (
	"errors"
	"log"
	"time"
)

//import (
//	"go.uber.org/zap"
//	"srbbs/src/conf"
//	"srbbs/src/srlogger"
//)

const (
	AccessTokenExpireDuration  = 24 * time.Hour //令牌过期时间
	RefreshTokenExpireDuration = 7 * 24 * time.Hour
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个UserID字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
//type MyClaims struct {
//	UserID   uint64 `json:"user_id"`
//	Username string `json:"username"`
//}

type MyCustomClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var secret = []byte("abab")

// 获取密钥的回调函数
func keyFunc(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//	return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	//}

	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	return []byte(conf.Cfg.JwtKey), nil
}

//func (claims *MyCustomClaims) GetExpirationTime() (*NumericDate, error) {
//
//}
//func (claims *MyCustomClaims) GetIssuedAt() (*NumericDate, error) {
//
//}
//func (claims *MyCustomClaims) GetNotBefore() (*NumericDate, error) {
//
//}
//func (claims *MyCustomClaims) GetIssuer() (string, error) {
//
//}
//func (claims *MyCustomClaims) GetSubject() (string, error) {
//
//}
//func (claims *MyCustomClaims) GetAudience() (ClaimStrings, error) {
//
//}

func GenToken(userID uint64, username string) (aToken, rToken string, err error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpireDuration)),
			Issuer:    "srbbs",
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	//这里注意要用
	aToken, err = token.SignedString([]byte(conf.Cfg.JwtKey))
	if err != nil {
		return
	}

	//rToken不需要自定义字段
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(RefreshTokenExpireDuration).Unix(),
		"iss": "srbbs",
	})
	rToken, err = token.SignedString([]byte(conf.Cfg.JwtKey))
	return
}

func ParseToken(tokenString string) (claims *MyCustomClaims, err error) {
	claims = &MyCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		srlogger.Logger().Info("error parse token", zap.Error(err))
	} else if _, ok := token.Claims.(*MyCustomClaims); ok {
		return
	} else {
		log.Fatal("unknown claims type, cannot proceed")
	}
	return
}

// RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// rToken失效直接return 这里用不带claims的Parse方法
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	//var claims *jwt.MapClaims
	claims, err := ParseToken(aToken)

	// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
	if errors.Is(err, jwt.ErrTokenExpired) {
		return GenToken(claims.UserID, claims.Username)
	}
	return
}
