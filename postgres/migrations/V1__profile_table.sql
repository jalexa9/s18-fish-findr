-- User table to store information about new users profiles.
-- Does not have links to other information yet but will add.
CREATE TABLE profile (
	user_id 				serial PRIMARY KEY,
	user_name 				TEXT UNIQUE NOT NULL,
	first_name 				TEXT NOT NULL,
	last_name 				TEXT NOT NULL,
	password 				TEXT NOT NULL,
	phone_number 			TEXT NOT NULL,
	email_address 			TEXT NOT NULL,
	facebook_profile_link 	TEXT,
	bio 					TEXT,
	updated_at				timestamptz DEFAULT NOW() NOT NULL
);

-- Table to contain types of interest
CREATE TABLE interest (
	interest_id 	serial PRIMARY KEY,
	interest_type 	TEXT UNIQUE NOT NULL
);

-- Interest to user relation table
CREATE TABLE user_interest (
	interest_id 			INT REFERENCES interest(interest_id) NOT NULL,
	user_id 			INT REFERENCES profile(user_id) NOT NULL,
	PRIMARY KEY (interest_id, user_id)
);

-- autoupdate is a generic trigger to update the updated_at column to current time.
CREATE FUNCTION autoupdate() RETURNS trigger AS $$
BEGIN
    new.updated_at = now();
    RETURN new;
END;
$$ LANGUAGE plpgsql;

-- profile_trigger updates the profile_trigger table's updated_at column on update.
CREATE TRIGGER profile_trigger
    BEFORE INSERT OR UPDATE
    ON profile FOR EACH ROW EXECUTE PROCEDURE autoupdate();

-- Some info for fake tester profiles to impot
INSERT INTO PROFILE (user_name, first_name, last_name, password, phone_number, email_address, facebook_profile_link, bio) VALUES ('mm1', 'Mitchell', 'McKenzie', 'password', '1234325134', 'email', 'facebook.com', 'my bio');
INSERT INTO PROFILE (user_name, first_name, last_name, password, phone_number, email_address, facebook_profile_link, bio) VALUES ('mm2', 'Mitchell', 'McKenzie', 'password', '1234325134', 'email', 'facebook.com', 'my bio');
INSERT INTO PROFILE (user_name, first_name, last_name, password, phone_number, email_address, facebook_profile_link, bio) VALUES ('mm3', 'Mitchell', 'McKenzie', 'password', '1234325134', 'email', 'facebook.com', 'my bio');
INSERT INTO PROFILE (user_name, first_name, last_name, password, phone_number, email_address, facebook_profile_link, bio) VALUES ('mm4', 'Mitchell', 'McKenzie', 'password', '1234325134', 'email', 'facebook.com', 'my bio');