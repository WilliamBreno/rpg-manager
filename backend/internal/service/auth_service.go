package service

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "rpg-manager/internal/domain"
    "rpg-manager/internal/repository"
)

type AuthService struct {
    UserRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
    return &AuthService{UserRepo: userRepo}
}

type Claims struct {
    UserID uint            `json:"user_id"`
    Role   domain.UserRole `json:"role"`
    jwt.RegisteredClaims
}

func (s *AuthService) Register(name, email, password string, role domain.UserRole) (domain.User, error) {
    // Verifica se o email já existe
    _, err := s.UserRepo.FindByEmail(email)
    if err == nil {
        return domain.User{}, errors.New("email já cadastrado")
    }

    // Gera o hash da senha
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return domain.User{}, errors.New("erro ao processar senha")
    }

    user := domain.User{
        Name:     name,
        Email:    email,
        Password: string(hashedPassword),
        Role:     role,
    }

    if err := s.UserRepo.Create(&user); err != nil {
        return domain.User{}, errors.New("erro ao criar usuário")
    }

    return user, nil
}

func (s *AuthService) Login(email, password string) (string, domain.User, error) {
    // Busca o usuário
    user, err := s.UserRepo.FindByEmail(email)
    if err != nil {
        return "", domain.User{}, errors.New("email ou senha inválidos")
    }

    // Valida a senha
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", domain.User{}, errors.New("email ou senha inválidos")
    }

    // Gera o token JWT
    token, err := s.generateToken(user)
    if err != nil {
        return "", domain.User{}, errors.New("erro ao gerar token")
    }

    return token, user, nil
}

func (s *AuthService) generateToken(user domain.User) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "rpg-manager-secret"
    }

    claims := Claims{
        UserID: user.ID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "rpg-manager-secret"
    }

    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })

    if err != nil || !token.Valid {
        return nil, errors.New("token inválido")
    }

    claims, ok := token.Claims.(*Claims)
    if !ok {
        return nil, errors.New("token inválido")
    }

    return claims, nil
}