package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Profile represents extended user profile information for Socialink (Facebook-like)
type Profile struct {
	ID                 uuid.UUID           `json:"id" db:"id"`
	UserID             uuid.UUID           `json:"user_id" db:"user_id"`
	Hometown           *string             `json:"hometown,omitempty" db:"hometown"`
	CurrentCity        *string             `json:"current_city,omitempty" db:"current_city"`
	RelationshipStatus *string             `json:"relationship_status,omitempty" db:"relationship_status"`
	Languages          StringArray         `json:"languages,omitempty" db:"languages"`
	InterestedIn       StringArray         `json:"interested_in,omitempty" db:"interested_in"`
	Work               []WorkExperience    `json:"work,omitempty" db:"work"`
	Education          []EducationEntry    `json:"education,omitempty" db:"education"`
	ContactInfo        *ContactInfo        `json:"contact_info,omitempty" db:"contact_info"`
	About              *string             `json:"about,omitempty" db:"about"`
	FavoriteQuotes     *string             `json:"favorite_quotes,omitempty" db:"favorite_quotes"`
	Hobbies            StringArray         `json:"hobbies,omitempty" db:"hobbies"`
	Website            *string             `json:"website,omitempty" db:"website"`
	SocialLinks        *SocialLinks        `json:"social_links,omitempty" db:"social_links"`
	FeaturedPhotos     StringArray         `json:"featured_photos,omitempty" db:"featured_photos"`
	Visibility         *ProfileVisibility  `json:"visibility,omitempty" db:"visibility"`
	CreatedAt          time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at" db:"updated_at"`
}

// WorkExperience represents employment history
type WorkExperience struct {
	ID          string     `json:"id"`
	Company     string     `json:"company"`
	Position    string     `json:"position"`
	City        *string    `json:"city,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	IsCurrent   bool       `json:"is_current"`
}

// EducationEntry represents educational background
type EducationEntry struct {
	ID              string     `json:"id"`
	School          string     `json:"school"`
	Degree          *string    `json:"degree,omitempty"`
	FieldOfStudy    *string    `json:"field_of_study,omitempty"`
	StartYear       int        `json:"start_year"`
	EndYear         *int       `json:"end_year,omitempty"`
	Description     *string    `json:"description,omitempty"`
	IsCurrentlyHere bool       `json:"is_currently_here"`
}

// ContactInfo represents contact information
type ContactInfo struct {
	Email          *string `json:"email,omitempty"`
	PhoneNumber    *string `json:"phone_number,omitempty"`
	Address        *string `json:"address,omitempty"`
	City           *string `json:"city,omitempty"`
	State          *string `json:"state,omitempty"`
	ZipCode        *string `json:"zip_code,omitempty"`
	Country        *string `json:"country,omitempty"`
}

// SocialLinks represents social media links
type SocialLinks struct {
	Instagram *string `json:"instagram,omitempty"`
	Twitter   *string `json:"twitter,omitempty"`
	LinkedIn  *string `json:"linkedin,omitempty"`
	YouTube   *string `json:"youtube,omitempty"`
	GitHub    *string `json:"github,omitempty"`
	Website   *string `json:"website,omitempty"`
}

// ProfileVisibility controls who can see profile information
type ProfileVisibility struct {
	Bio              string `json:"bio"`               // public, friends, only_me
	Work             string `json:"work"`              // public, friends, only_me
	Education        string `json:"education"`         // public, friends, only_me
	ContactInfo      string `json:"contact_info"`      // friends, only_me
	RelationshipInfo string `json:"relationship_info"` // public, friends, only_me
	Hometown         string `json:"hometown"`          // public, friends, only_me
	Birthday         string `json:"birthday"`          // public, friends, only_me
}

// StringArray custom type for PostgreSQL array handling
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

// JSON type handlers for complex types
func (w WorkExperience) Value() (driver.Value, error) {
	return json.Marshal(w)
}

func (w *WorkExperience) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, w)
}

func (e EducationEntry) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *EducationEntry) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, e)
}

func (c ContactInfo) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *ContactInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}

func (s SocialLinks) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SocialLinks) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

func (p ProfileVisibility) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *ProfileVisibility) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, p)
}

// DTOs for profile management

// UpdateProfileInfoRequest for basic profile updates
type UpdateProfileInfoRequest struct {
	Hometown           *string    `json:"hometown,omitempty" binding:"omitempty,max=100"`
	CurrentCity        *string    `json:"current_city,omitempty" binding:"omitempty,max=100"`
	RelationshipStatus *string    `json:"relationship_status,omitempty" binding:"omitempty,oneof=single in_relationship engaged married divorced widowed"`
	About              *string    `json:"about,omitempty" binding:"omitempty,max=1000"`
	FavoriteQuotes     *string    `json:"favorite_quotes,omitempty" binding:"omitempty,max=500"`
	Website            *string    `json:"website,omitempty" binding:"omitempty,url"`
}

// AddWorkExperienceRequest for adding work history
type AddWorkExperienceRequest struct {
	Company     string  `json:"company" binding:"required,max=100"`
	Position    string  `json:"position" binding:"required,max=100"`
	City        *string `json:"city,omitempty" binding:"omitempty,max=100"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=500"`
	StartDate   string  `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate     *string `json:"end_date,omitempty"`            // YYYY-MM-DD
	IsCurrent   bool    `json:"is_current"`
}

// AddEducationRequest for adding education
type AddEducationRequest struct {
	School          string  `json:"school" binding:"required,max=200"`
	Degree          *string `json:"degree,omitempty" binding:"omitempty,max=100"`
	FieldOfStudy    *string `json:"field_of_study,omitempty" binding:"omitempty,max=100"`
	StartYear       int     `json:"start_year" binding:"required,min=1900,max=2100"`
	EndYear         *int    `json:"end_year,omitempty" binding:"omitempty,min=1900,max=2100"`
	Description     *string `json:"description,omitempty" binding:"omitempty,max=500"`
	IsCurrentlyHere bool    `json:"is_currently_here"`
}

// UpdateContactInfoRequest for updating contact information
type UpdateContactInfoRequest struct {
	Email       *string `json:"email,omitempty" binding:"omitempty,email"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Address     *string `json:"address,omitempty" binding:"omitempty,max=200"`
	City        *string `json:"city,omitempty" binding:"omitempty,max=100"`
	State       *string `json:"state,omitempty" binding:"omitempty,max=100"`
	ZipCode     *string `json:"zip_code,omitempty" binding:"omitempty,max=20"`
	Country     *string `json:"country,omitempty" binding:"omitempty,max=100"`
}

// UpdateSocialLinksRequest for updating social media links
type UpdateSocialLinksRequest struct {
	Instagram *string `json:"instagram,omitempty" binding:"omitempty,url"`
	Twitter   *string `json:"twitter,omitempty" binding:"omitempty,url"`
	LinkedIn  *string `json:"linkedin,omitempty" binding:"omitempty,url"`
	YouTube   *string `json:"youtube,omitempty" binding:"omitempty,url"`
	GitHub    *string `json:"github,omitempty" binding:"omitempty,url"`
	Website   *string `json:"website,omitempty" binding:"omitempty,url"`
}

// UpdateVisibilityRequest for privacy settings
type UpdateVisibilityRequest struct {
	Bio              *string `json:"bio,omitempty" binding:"omitempty,oneof=public friends only_me"`
	Work             *string `json:"work,omitempty" binding:"omitempty,oneof=public friends only_me"`
	Education        *string `json:"education,omitempty" binding:"omitempty,oneof=public friends only_me"`
	ContactInfo      *string `json:"contact_info,omitempty" binding:"omitempty,oneof=friends only_me"`
	RelationshipInfo *string `json:"relationship_info,omitempty" binding:"omitempty,oneof=public friends only_me"`
	Hometown         *string `json:"hometown,omitempty" binding:"omitempty,oneof=public friends only_me"`
	Birthday         *string `json:"birthday,omitempty" binding:"omitempty,oneof=public friends only_me"`
}

// ProfileResponse for API responses
type ProfileResponse struct {
	ID                 uuid.UUID          `json:"id"`
	UserID             uuid.UUID          `json:"user_id"`
	User               *UserResponse      `json:"user,omitempty"`
	Hometown           *string            `json:"hometown,omitempty"`
	CurrentCity        *string            `json:"current_city,omitempty"`
	RelationshipStatus *string            `json:"relationship_status,omitempty"`
	Languages          []string           `json:"languages,omitempty"`
	InterestedIn       []string           `json:"interested_in,omitempty"`
	Work               []WorkExperience   `json:"work,omitempty"`
	Education          []EducationEntry   `json:"education,omitempty"`
	ContactInfo        *ContactInfo       `json:"contact_info,omitempty"`
	About              *string            `json:"about,omitempty"`
	FavoriteQuotes     *string            `json:"favorite_quotes,omitempty"`
	Hobbies            []string           `json:"hobbies,omitempty"`
	Website            *string            `json:"website,omitempty"`
	SocialLinks        *SocialLinks       `json:"social_links,omitempty"`
	FeaturedPhotos     []string           `json:"featured_photos,omitempty"`
	Visibility         *ProfileVisibility `json:"visibility,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

// ToProfileResponse converts Profile to ProfileResponse
func (p *Profile) ToProfileResponse() *ProfileResponse {
	return &ProfileResponse{
		ID:                 p.ID,
		UserID:             p.UserID,
		Hometown:           p.Hometown,
		CurrentCity:        p.CurrentCity,
		RelationshipStatus: p.RelationshipStatus,
		Languages:          p.Languages,
		InterestedIn:       p.InterestedIn,
		Work:               p.Work,
		Education:          p.Education,
		ContactInfo:        p.ContactInfo,
		About:              p.About,
		FavoriteQuotes:     p.FavoriteQuotes,
		Hobbies:            p.Hobbies,
		Website:            p.Website,
		SocialLinks:        p.SocialLinks,
		FeaturedPhotos:     p.FeaturedPhotos,
		Visibility:         p.Visibility,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}
