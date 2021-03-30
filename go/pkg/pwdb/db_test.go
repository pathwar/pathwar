package pwdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDB_SoftDelete(t *testing.T) {
	t.Helper()
	db := TestingSqliteDB(t, zap.NewNop())

	user := User{
		Email: "test@test.com",
	}
	err := db.Create(&user).Error
	require.NoError(t, err)

	var usersCount int64
	err = db.Model(&User{}).Count(&usersCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), usersCount)

	err = db.Delete(&user).Error
	require.NoError(t, err)

	err = db.Model(&User{}).Count(&usersCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), usersCount)

	// getting soft deleted instances too
	err = db.Unscoped().Model(&User{}).Count(&usersCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), usersCount)

	// hard delete
	err = db.Unscoped().Delete(&user).Error
	require.NoError(t, err)

	err = db.Unscoped().Model(&User{}).Count(&usersCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), usersCount)
}
