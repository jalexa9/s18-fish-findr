package postgres

import (
	"strconv"

	"github.com/Clemson-CPSC-4910/s18-fish-findr/go/fisher"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// GetProfiles will get the current profiles from the DB.
func (db *DB) GetProfiles() ([]fisher.Profile, error) {
	var fishers []fisher.Profile
	var interest []idForType
	var finalFishers []fisher.Profile

	err := db.Transact(func(tx *sqlx.Tx) error {
		err := tx.Select(&fishers, `
			SELECT * FROM profile`)

		if err != nil {
			return errors.Wrapf(err, "Error getting profile.")
		}
		err = tx.Select(&interest,
			`SELECT user_id, u.interest_id, interest_type 
				FROM user_interest u 
				INNER JOIN interest i ON i.interest_id = u.interest_id;`)
		if err != nil {
			return errors.Wrapf(err, "Error getting interests.")
		}
		for _, p := range fishers {
			for _, i := range interest {
				if p.ID == i.UserID {
					p.Interest = append(p.Interest, fisher.Interest{
						ID:   i.InterestID,
						Type: i.Type,
					})
				}
			}
			finalFishers = append(finalFishers, p)
		}
		return errors.Wrapf(err, "Error adding interest to profile")
	})
	return finalFishers, errors.Wrapf(err, "Could not get profiles from the DB.")
}

// GetProfileByUserName will take a user name and return a profile of that type or error.
func (db *DB) GetProfileByUserName(s string) (fisher.Profile, error) {
	var fishers []fisher.Profile
	var interest []idForType

	err := db.Transact(func(tx *sqlx.Tx) error {
		err := tx.Select(&fishers, `
			SELECT * FROM profile WHERE user_name = '`+s+`'`)
		if err != nil {
			return errors.Wrapf(err, "Error getting profile.")
		}
		err = tx.Select(&interest,
			`SELECT user_id, u.interest_id, interest_type 
				FROM user_interest u
				INNER JOIN interest i ON i.interest_id = u.interest_id;`)
		if err != nil {
			return errors.Wrapf(err, "Error getting interests.")
		}
		for _, i := range interest {
			if fishers[0].ID == i.UserID {
				fishers[0].Interest = append(fishers[0].Interest, fisher.Interest{
					ID:   i.InterestID,
					Type: i.Type,
				})
			}
		}
		return errors.Wrapf(err, "Error adding interest to profile")
	})
	if len(fishers) != 1 {
		return fisher.Profile{}, errors.Wrapf(err, "Did not get a single profile.")
	}
	return fishers[0], errors.Wrapf(err, "Error getting profile.")
}

// GetProfileInterest gets a list of interest for the profile user_name.
func (db *DB) GetProfileInterest(s string) ([]fisher.Interest, error) {
	// MAYBE, TODO
	return nil, nil
}

// GetProfileMatches gets the ordered profile matches from the DB.
func (db *DB) GetProfileMatches(s string) ([]fisher.Profile, error) {
	profile, err := db.GetProfileByUserName(s)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting info on profile.")
	}
	var interest []idForType
	var matchProfiles []fisher.Profile
	var finalFishers []fisher.Profile
	var params []interface{}
	if len(profile.Interest) == 0 {
		return nil, nil
	}
	for _, i := range profile.Interest {
		params = append(params, i.ID)
	}
	err = db.Transact(func(tx *sqlx.Tx) error {
		err = tx.Select(&matchProfiles, `
			SELECT 
				p.user_id, 
				p.first_name, 
				p.last_name, 
				p.user_name, 
				p.password, 
				p.phone_number,
				p.email_address, 
				p.facebook_profile_link, 
				p.bio 
			FROM profile p INNER JOIN user_interest u 
			ON u.user_id = p.user_id 
			WHERE u.user_id != `+strconv.Itoa(profile.ID)+` 
			AND u.interest_id IN `+buildValues(len(params), 1)+` 
			GROUP BY (p.user_id) 
			ORDER BY COUNT(u.user_id) DESC;
		`, params...)
		if err != nil {
			return errors.Wrapf(err, "Error getting matching interest from DB.")
		}
		err = tx.Select(&interest,
			`SELECT user_id, u.interest_id, interest_type 
				FROM user_interest u
				INNER JOIN interest i ON i.interest_id = u.interest_id;`)
		if err != nil {
			return errors.Wrapf(err, "Error getting interests in matching.")
		}
		for _, p := range matchProfiles {
			for _, i := range interest {
				if p.ID == i.UserID {
					p.Interest = append(p.Interest, fisher.Interest{
						ID:   i.InterestID,
						Type: i.Type,
					})
				}
			}
			finalFishers = append(finalFishers, p)
		}
		return nil
	})
	return finalFishers, errors.Wrapf(err, "Error getting profile matches.")
}

// UpdateProfile accepts a profile and updates it or creates it in the DB.
func (db *DB) UpdateProfile(profile fisher.Profile) error {
	err := db.Transact(func(tx *sqlx.Tx) error {
		var params []interface{}
		params = append(params, profile.UserName, profile.FirstName, profile.LastName,
			profile.Password, profile.PhoneNumber, profile.Email, profile.Facebook, profile.Bio)
		query := `
			INSERT INTO profile (user_name, first_name, last_name, password,
				phone_number, email_address, facebook_profile_link, bio)
			VALUES ` + buildValues(8, 1) + `
			ON CONFLICT (user_name) DO UPDATE SET
				user_name = EXCLUDED.user_name,
				first_name = EXCLUDED.first_name,
				last_name = EXCLUDED.last_name,
				password = EXCLUDED.password,
				phone_number = EXCLUDED.phone_number,
				email_address = EXCLUDED.email_address,
				facebook_profile_link = EXCLUDED.facebook_profile_link,
				bio = EXCLUDED.bio
		`
		_, err := tx.Exec(query, params...)
		if err != nil {
			return errors.Wrapf(err, "Error inserting the profile into the database.")
		}
		err = db.UpdateInterest(profile, tx)
		return errors.Wrapf(err, "Error updating interest for profile.")
	})
	return err
}

// UpdateInterest will take a profile and update that profiles interest.
func (db *DB) UpdateInterest(profile fisher.Profile, tx *sqlx.Tx) error {
	var params []interface{}
	var count int
	var userID int
	for _, i := range profile.Interest {
		ids, err := db.getIDsForType(i.Type, profile.UserName, tx)
		if err != nil {
			return errors.Wrapf(err, "Error getting interest id from DB")
		}
		params = append(params, ids.InterestID, ids.UserID)
		userID = ids.UserID
		count++
	}
	err := db.pruneUserInterest(userID, tx)
	if err != nil {
		return errors.Wrapf(err, "Error pruning interest.")
	}
	if count != 0 {
		query := `
		INSERT INTO user_interest (interest_id, user_id)
		VALUES ` + buildValues(2, count) + `
		ON CONFLICT (interest_id, user_id) DO UPDATE SET
			interest_id = EXCLUDED.interest_id,
			user_id = EXCLUDED.user_id
		`
		_, err = tx.Exec(query, params...)
	}

	return errors.Wrapf(err, "Error inserting the user interest into the database.")
}

func (db *DB) pruneUserInterest(i int, tx *sqlx.Tx) error {
	_, err := tx.Exec(`DELETE FROM user_interest WHERE user_id = ` + strconv.Itoa(i))
	return errors.Wrapf(err, "Error deleting from user interest for this user.")
}

func (db *DB) getIDsForType(t string, un string, tx *sqlx.Tx) (idForType, error) {
	var ids []idForType
	err := tx.Select(&ids, `
		SELECT interest_id, user_id FROM interest, profile 
		WHERE interest_type = '`+t+`' AND user_name = '`+un+`'`)
	return ids[0], err
}

func (db *DB) getIDforInterest(t string, tx *sqlx.Tx) (int, error) {
	var ids []int
	err := tx.Select(&ids, `
		SELECT interest_id FROM interest 
		WHERE interest_type = '`+t+`'`)
	if err != nil {
		return -1, errors.Wrapf(err, "Error getting interest id from the DB.")
	}
	return ids[0], err
}

type idForType struct {
	InterestID int    `db:"interest_id"`
	UserID     int    `db:"user_id"`
	Type       string `db:"interest_type"`
}

// GetIfLogin will see if the user name and password are valid.
func (db *DB) GetIfLogin(un string, p string) (fisher.Profile, error) {
	var fishers []fisher.Profile

	err := db.Transact(func(tx *sqlx.Tx) error {
		err := tx.Select(&fishers, `
			SELECT * FROM profile WHERE user_name = '`+un+`' AND password = '`+p+`'`)

		return err
	})
	if len(fishers) != 1 {
		return fisher.Profile{}, errors.Wrapf(errors.New("error"), "Could not get profiles from the DB.")
	}
	return fishers[0], errors.Wrapf(err, "Could not get profiles from the DB.")
}

// GetProfilesWithInterest returns an array of users with the interest sorted by most matching to least.
// These users will no include the user who is sending the request.
func (db *DB) GetProfilesWithInterest(s string, subInt []fisher.Interest) ([]fisher.Profile, error) {
	profile, err := db.GetProfileByUserName(s)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting info on profile.")
	}
	var interest []idForType
	var matchProfiles []fisher.Profile
	var finalFishers []fisher.Profile
	var params []interface{}
	if len(profile.Interest) == 0 {
		return nil, nil
	}
	err = db.Transact(func(tx *sqlx.Tx) error {
		for _, i := range subInt {
			intID, err := db.getIDforInterest(i.Type, tx)
			if err != nil {
				return err
			}
			params = append(params, intID)
		}
		err = tx.Select(&matchProfiles, `
			SELECT 
				p.user_id, 
				p.first_name, 
				p.last_name, 
				p.user_name, 
				p.password, 
				p.phone_number,
				p.email_address, 
				p.facebook_profile_link, 
				p.bio 
			FROM profile p INNER JOIN user_interest u 
			ON u.user_id = p.user_id 
			WHERE u.user_id != `+strconv.Itoa(profile.ID)+` 
			AND u.interest_id IN `+buildValues(len(params), 1)+` 
			GROUP BY (p.user_id) 
			ORDER BY COUNT(u.user_id) DESC;
		`, params...)
		if err != nil {
			return errors.Wrapf(err, "Error getting matching interest from DB.")
		}
		err = tx.Select(&interest,
			`SELECT user_id, u.interest_id, interest_type 
				FROM user_interest u
				INNER JOIN interest i ON i.interest_id = u.interest_id;`)
		if err != nil {
			return errors.Wrapf(err, "Error getting interests in matching.")
		}
		for _, p := range matchProfiles {
			for _, i := range interest {
				if p.ID == i.UserID {
					p.Interest = append(p.Interest, fisher.Interest{
						ID:   i.InterestID,
						Type: i.Type,
					})
				}
			}
			finalFishers = append(finalFishers, p)
		}
		return nil
	})
	return finalFishers, errors.Wrapf(err, "Error getting profile matches.")
}
