package users

import (
	"github.com/bubulearn/bubucore"
	"github.com/bubulearn/bubucore/i18n"
	"net/mail"
	"time"
)

// User represents a common user data, students and teachers both
type User struct {
	ID string `json:"id" example:"4452dda6-4fde-453f-a41d-4c043e0ea6d1"`

	// Role ID, 1=student, 500=teacher, 1000=admin
	Role  int    `json:"role" example:"500"`
	Email string `json:"email" example:"john@example.com"`

	Name  string `json:"name" example:"John Doe"`
	Phone string `json:"phone" example:"+7 (900) 123-45-67"`

	// Gender ID, 0=unset, 1=male, 2=female
	Gender byte `json:"gender" example:"1"`

	// Code for teacher to assign to the student
	Code string `json:"code" example:"123456"`

	Language        i18n.Language `json:"lang" swaggertype:"string" enums:"en,ru" example:"ru"`
	LessonsLanguage i18n.Language `json:"lessons_lang" swaggertype:"string" enums:"en,ru" example:"ru"`

	// TeacherID for student, ID of teacher assigned
	TeacherID string `json:"teacher_id,omitempty" example:"b21b949e-8495-4f56-ab9e-502199af48cf"`

	IsBlocked     bool   `json:"is_blocked" example:"false"`
	BlockedReason string `json:"blocked_reason,omitempty" example:"Really bad ass"`

	TimeCreated   *time.Time `json:"time_created" example:"2021-07-12T20:29:40+03:00"`
	TimeUpdated   *time.Time `json:"time_updated" example:"2021-07-12T20:29:40+03:00"`
	TimeLastLogin *time.Time `json:"time_last_login" example:"2021-07-12T20:29:40+03:00"`

	// List of assigned students to the teacher
	StudentsAssigned []*User `json:"students_assigned,omitempty"`

	// DeviceTokens is a list of user's registered devices tokens
	DeviceTokens []string `json:"device_tokens"`
}

// IsStudent checks if user has a student role
func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

// IsTeacher checks if user has a teacher role
func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

// IsAdmin checks if user has an admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// GetEmailAddress return user's name and email as email address instance
func (u *User) GetEmailAddress() *mail.Address {
	return &mail.Address{
		Address: u.Email,
		Name:    u.Name,
	}
}

// ValidateUnblocked check if user is not blocked.
// Returns common.ErrUserBlocked in case of block.
func (u *User) ValidateUnblocked() error {
	if u.IsBlocked {
		return bubucore.ErrUserBlocked
	}
	return nil
}
