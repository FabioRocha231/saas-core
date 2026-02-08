package pkg

import (
	"errors"
	"time"

	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
	ttl    time.Duration
	issuer string
	uuid   ports.UUIDInterface
}

func NewJwtService(secret string, ttl time.Duration, issuer string, uuid ports.UUIDInterface) ports.JwtInterface {
	return &Service{
		secret: []byte(secret),
		ttl:    ttl,
		issuer: issuer,
		uuid:   uuid,
	}
}

func (s *Service) Sign(userID string, role string) (string, error) {
	now := time.Now()
	claims := ports.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        s.uuid.Generate(),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}

func (s *Service) Parse(tokenStr string) (*ports.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &ports.Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ports.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *Service) GetJTI(tokenStr string) string {
	claims, err := s.Parse(tokenStr)
	if err != nil {
		return ""
	}
	return claims.ID
}

func (s *Service) GetExpiresAt(tokenStr string) time.Time {
	claims, err := s.Parse(tokenStr)
	if err != nil {
		return time.Time{}
	}
	return claims.ExpiresAt.Time
}
