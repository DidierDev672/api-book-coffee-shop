package usecase

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type AuthUseCase interface {
	Register(token, nameFull, phone, idNumber, dateOfBirth, email, password string) (*domain.User, error)
	Login(token, email, password string) (*domain.User, string, error)
	GetAll() ([]*domain.User, error)
	GetProfile(id string) (*domain.User, error)
}

type authUseCase struct {
	repo   repository.UserRepository
	hasher repository.PasswordHasher
	tokens repository.TokenService
}

func NewAuthUseCase(
	repo repository.UserRepository,
	hasher repository.PasswordHasher,
	tokens repository.TokenService,
) AuthUseCase {
	return &authUseCase{repo: repo, hasher: hasher, tokens: tokens}
}

func (uc *authUseCase) Register(token, nameFull, phone, idNumber, dateOfBirth, email, password string) (*domain.User, error) {
	if err := validateRegisterFields(nameFull, phone, idNumber, dateOfBirth, email, password); err != nil {
		return nil, err
	}

	email = strings.ToLower(strings.TrimSpace(email))
	if _, err := uc.repo.GetByEmail(email); err == nil {
		return nil, errors.New("email already registered")
	}

	hash, err := uc.hasher.Hash(password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	now := time.Now()
	u := &domain.User{
		ID:           generateID(),
		NameFull:     strings.TrimSpace(nameFull),
		Phone:        strings.TrimSpace(phone),
		IDNumber:     strings.TrimSpace(idNumber),
		DateOfBirth:  dateOfBirth,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *authUseCase) Login(_ string, email, password string) (*domain.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, "", errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, "", errors.New("password cannot be empty")
	}

	u, err := uc.repo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if err := uc.hasher.Compare(u.PasswordHash, password); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	newToken, err := uc.tokens.Generate(u.ID)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}
	if err := uc.repo.UpdateAuthToken(u.ID, newToken); err != nil {
		return nil, "", err
	}
	return u, newToken, nil
}

func (uc *authUseCase) GetAll() ([]*domain.User, error) {
	return uc.repo.GetAll()
}

func (uc *authUseCase) GetProfile(id string) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *authUseCase) verifyToken(token string) error {
	if token == "" {
		return errors.New("authorization token is required")
	}
	if _, err := uc.tokens.Validate(token); err != nil {
		return err
	}
	if _, err := uc.repo.GetByAuthToken(token); err != nil {
		return errors.New("invalid or expired token")
	}
	return nil
}

func validateRegisterFields(nameFull, phone, idNumber, dateOfBirth, email, password string) error {
	if strings.TrimSpace(nameFull) == "" {
		return errors.New("name_full cannot be empty")
	}
	if strings.TrimSpace(phone) == "" {
		return errors.New("phone cannot be empty")
	}
	if strings.TrimSpace(idNumber) == "" {
		return errors.New("id_number cannot be empty")
	}
	if !isNumeric(idNumber) {
		return errors.New("id_number must be numeric")
	}
	if dateOfBirth == "" {
		return errors.New("date_of_birth cannot be empty")
	}
	if _, err := time.Parse("2006-01-02", dateOfBirth); err != nil {
		return errors.New("date_of_birth must be in YYYY-MM-DD format")
	}
	if strings.TrimSpace(email) == "" {
		return errors.New("email cannot be empty")
	}
	if !emailPattern.MatchString(strings.TrimSpace(email)) {
		return errors.New("email format is invalid")
	}
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func isNumeric(s string) bool {
	for _, r := range strings.TrimSpace(s) {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(strings.TrimSpace(s)) > 0
}
