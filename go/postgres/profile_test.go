package postgres

import (
	"testing"

	"github.com/Clemson-CPSC-4910/s18-fish-findr/go/fisher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProfiles(t *testing.T) {
	db, err := CreateDB()
	require.Nil(t, err)
	err = db.CleanDB(t)
	require.Nil(t, err)
	defer db.CloseOrFail(t)

	db.assertCountRows(t, 0, "profile")

	// TODO, add test for info stored inside of the rows
	err = db.UpdateProfile(fisher.Profile{
		ID:          1,
		FirstName:   "FTest",
		LastName:    "LName",
		UserName:    "user1",
		Password:    "testPASS",
		PhoneNumber: "123456789",
		Email:       "test.email",
		Facebook:    "www.facebook.com",
		Bio:         "test bio",
		Interest: []fisher.Interest{
			fisher.Interest{
				Type: "fly fishing",
			},
		},
	})
	assert.Nil(t, err)
	db.assertCountRows(t, 1, "profile")
	db.assertCountRows(t, 1, "user_interest")
	err = db.UpdateProfile(fisher.Profile{
		ID:          2,
		FirstName:   "FTest",
		LastName:    "LName",
		UserName:    "user2",
		Password:    "testPASS",
		PhoneNumber: "123456789",
		Email:       "test.email",
		Facebook:    "www.facebook.com",
		Bio:         "test bio",
		Interest: []fisher.Interest{
			fisher.Interest{
				Type: "fly fishing",
			},
			fisher.Interest{
				Type: "redfish",
			},
		},
	})
	assert.Nil(t, err)
	db.assertCountRows(t, 2, "profile")
	db.assertCountRows(t, 3, "user_interest")
	err = db.UpdateProfile(fisher.Profile{
		ID:          3,
		FirstName:   "FTest",
		LastName:    "LName",
		UserName:    "user3",
		Password:    "testPASS",
		PhoneNumber: "123456789",
		Email:       "test.email",
		Facebook:    "www.facebook.com",
		Bio:         "test bio",
		Interest:    nil,
	})
	assert.Nil(t, err)
	db.assertCountRows(t, 3, "profile")
	db.assertCountRows(t, 3, "user_interest")

	db.Close()
}

func TestUpdateProfile(t *testing.T) {
	db, err := CreateDB()
	require.Nil(t, err)
	err = db.CleanDB(t)
	require.Nil(t, err)
	defer db.CloseOrFail(t)

	db.assertCountRows(t, 0, "profile")
	err = db.UpdateProfile(fisher.Profile{
		ID:          1,
		FirstName:   "FTest",
		LastName:    "LName",
		UserName:    "user1",
		Password:    "testPASS",
		PhoneNumber: "123456789",
		Email:       "test.email",
		Facebook:    "www.facebook.com",
		Bio:         "test bio",
		Interest: []fisher.Interest{
			fisher.Interest{
				Type: "fly fishing",
			},
		},
	})
	assert.Nil(t, err)
	db.assertCountRows(t, 1, "profile")
	db.assertCountRows(t, 1, "user_interest")
	db.assertValueOfWhere(t, 1, "profile", "first_name='FTest'")
	db.assertValueOfWhere(t, 1, "user_interest", "interest_id=1")

	err = db.UpdateProfile(fisher.Profile{
		ID:          2,
		FirstName:   "FTest",
		LastName:    "LName",
		UserName:    "user2",
		Password:    "testPASS",
		PhoneNumber: "123456789",
		Email:       "test.email",
		Facebook:    "www.facebook.com",
		Bio:         "test bio",
		Interest: []fisher.Interest{
			fisher.Interest{
				Type: "fly fishing",
			},
			fisher.Interest{
				Type: "redfish",
			},
		},
	})
	assert.Nil(t, err)
	db.assertCountRows(t, 2, "profile")
	db.assertCountRows(t, 3, "user_interest")
	db.assertValueOfWhere(t, 2, "profile", "first_name='FTest'")
	db.assertValueOfWhere(t, 2, "user_interest", "interest_id=1")
	db.assertValueOfWhere(t, 1, "user_interest", "interest_id=3")

	// Test changing data within a profile on update.
	err = db.UpdateProfile(fisher.Profile{
		ID:          2,
		FirstName:   "NewName",
		LastName:    "LName",
		UserName:    "user2",
		Password:    "newPass",
		PhoneNumber: "newPhone",
		Email:       "new.email",
		Facebook:    "www.facebook.com",
		Bio:         "new bio",
		Interest:    nil,
	})
	assert.Nil(t, err)
	db.assertCountRows(t, 2, "profile")
	db.assertCountRows(t, 3, "user_interest")
	db.assertValueOfWhere(t, 1, "profile", "first_name='FTest'")
	db.assertValueOfWhere(t, 1, "profile", "first_name='NewName'")
	db.assertValueOfWhere(t, 1, "profile", "password='newPass'")
	db.assertValueOfWhere(t, 1, "profile", "phone_number='newPhone'")
	db.assertValueOfWhere(t, 1, "profile", "email_address='new.email'")
	db.assertValueOfWhere(t, 1, "profile", "bio='new bio'")

	// TODO, make it where this is no longer true, if nil then remove this users interst.
	db.assertValueOfWhere(t, 2, "user_interest", "interest_id=1")
	db.assertValueOfWhere(t, 1, "user_interest", "interest_id=3")
	db.Close()
}

func TestGetProfileByUserName(t *testing.T) {

}
