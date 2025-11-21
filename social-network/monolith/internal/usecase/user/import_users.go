package user

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

func (uc *UseCases) ImportUsers(ctx context.Context, filePath string) error {
	//nolint:gosec // путь к файлу определяется выше по стеку, не из пользовательского ввода.
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close() //nolint:errcheck // не нужно обрабатывать.

	csvReader := csv.NewReader(f)

	const batchSize = 1000

	users := make([]domain.User, 0, batchSize)

Loop:
	for {
		row, err := csvReader.Read()

		switch {
		case errors.Is(err, io.EOF):
			break Loop
		case err != nil:
			return fmt.Errorf("failed to read records: %w", err)
		}

		if len(row) != 3 {
			return fmt.Errorf("expected 3 records, got %d", len(row))
		}

		user, err := parseUserFromRow(row)
		if err != nil {
			return fmt.Errorf("failed to parse user: %w", err)
		}

		users = append(users, user)

		if len(users) >= batchSize {
			err = uc.saveUsers(ctx, users)
			if err != nil {
				return fmt.Errorf("failed to save users: %w", err)
			}

			users = make([]domain.User, 0, batchSize)
		}
	}

	if len(users) > 0 {
		err = uc.saveUsers(ctx, users)
		if err != nil {
			return fmt.Errorf("failed to save users: %w", err)
		}
	}

	return nil
}

func (uc *UseCases) saveUsers(ctx context.Context, users []domain.User) error {
	f := func(ctx context.Context, tx transaction.Tx) error {
		err := uc.repo.WithTx(tx).MassCreate(ctx, users)
		if err != nil {
			return fmt.Errorf("failed to mass create users: %w", err)
		}

		return nil
	}

	err := uc.tx.Exec(ctx, f, nil)
	if err != nil {
		return fmt.Errorf("failed to save users in tx: %w", err)
	}

	return nil
}

func parseUserFromRow(row []string) (domain.User, error) {
	var out domain.User

	const (
		lastNameFirstNameIdx = iota
		birthDateIdx
		cityIdx
	)

	//nolint:gosec // хэш пароля не может быть пустым,
	// поэтому присваиваем всем пользователям из файла один и тот же хэш.
	const passwordHash = "$2a$10$WJDqvagiuI8eMSkrxjBn.OGYvKjjbKFB/fSLS2I56ebnSZWujFz.e"

	fullName := row[lastNameFirstNameIdx]
	birthDate := row[birthDateIdx]
	city := row[cityIdx]

	nameSplit := strings.Split(fullName, " ")
	if len(nameSplit) != 2 {
		return out, fmt.Errorf("failed to parse full name: %s", fullName)
	}

	lastName, firstName := nameSplit[0], nameSplit[1]

	birthDateTime, err := time.Parse(time.DateOnly, birthDate)
	if err != nil {
		return out, fmt.Errorf("failed to parse birth date '%s': %w", birthDate, err)
	}

	return domain.User{
		ID:           uuid.New(),
		PasswordHash: passwordHash,
		FirstName:    firstName,
		SecondName:   lastName,
		BirthDate:    birthDateTime,
		Gender:       defineGender(lastName),
		Biography:    "",
		City:         city,
	}, nil
}

func defineGender(lastName string) domain.Gender {
	if strings.HasSuffix(lastName, "ва") ||
		strings.HasSuffix(lastName, "на") {
		return domain.GenderFemale
	}

	return domain.GenderMale
}
