package users

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = mockDB.Close() })

	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: mockDB}),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	require.NoError(t, err)

	return gormDB, mock
}

func testEditReq() *EditProfileRequest {
	return &EditProfileRequest{
		Username:       "newuser",
		Email:          "new@example.com",
		Address:        "456 New St",
		ProfilePicture: "https://example.com/avatar.jpg",
	}
}

var userColumns = []string{
	"id", "username", "email", "password", "address",
	"province_id", "city_id", "district_id", "profile_picture", "role",
	"created_at", "updated_at", "deleted_at",
}

func newUserRow() *sqlmock.Rows {
	return sqlmock.NewRows(userColumns).
		AddRow(
			int64(1), "newuser", "new@example.com", "hashed", "456 New St",
			int64(1), int64(2), int64(3), "https://example.com/avatar.jpg", "user",
			nil, nil, nil,
		)
}

func TestGetProfile_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := get{repository: gormDB}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(newUserRow())

	user, err := g.GetProfile(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, uint(1), user.ID)
	require.Equal(t, "newuser", user.Username)
	require.Equal(t, "new@example.com", user.Email)
	require.Equal(t, "456 New St", user.Address)
	require.Equal(t, "https://example.com/avatar.jpg", user.ProfilePicture)
	require.Equal(t, "user", user.Role)
}

func TestGetProfile_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	g := get{repository: gormDB}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id =`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := g.GetProfile(context.Background(), 999)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUpdate_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testEditReq()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WithArgs("newuser", "new@example.com", "456 New St", "https://example.com/avatar.jpg", sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnRows(newUserRow())

	user, err := u.Update(context.Background(), req, 1)
	require.NoError(t, err)
	require.Equal(t, uint(1), user.ID)
	require.Equal(t, "newuser", user.Username)
	require.Equal(t, "new@example.com", user.Email)
	require.Equal(t, "456 New St", user.Address)
	require.Equal(t, "https://example.com/avatar.jpg", user.ProfilePicture)
}

func TestUpdate_ErrorOnUpdate(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testEditReq()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WithArgs("newuser", "new@example.com", "456 New St", "https://example.com/avatar.jpg", sqlmock.AnyArg(), int64(1)).
		WillReturnError(errors.New("update failed"))

	_, err := u.Update(context.Background(), req, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "update failed")
}

func TestUpdate_ErrorOnGetProfile(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	getOp := get{repository: gormDB}
	u := update{repository: gormDB, get: getOp}
	req := testEditReq()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WithArgs("newuser", "new@example.com", "456 New St", "https://example.com/avatar.jpg", sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id =`)).
		WithArgs(int64(1), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := u.Update(context.Background(), req, 1)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
