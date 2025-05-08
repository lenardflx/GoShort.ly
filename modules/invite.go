package modules

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"urlshort-backend/db"
	"urlshort-backend/models"
)

// CreateInvite deletes existing invites for email and returns a new one
func CreateInvite(email string) (*models.Invite, error) {
	// Cleanup old invites for the same email
	_ = db.DB.Where("email = ?", email).Unscoped().Delete(&models.Invite{}).Error

	// Generate secure token
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	invite := &models.Invite{
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(90 * 24 * time.Hour),
	}

	if err := db.DB.Create(invite).Error; err != nil {
		return nil, err
	}

	return invite, nil
}

// GetActiveInvite returns a valid invite for email + token
func GetActiveInvite(email, token string) (*models.Invite, error) {
	var invite models.Invite
	err := db.DB.
		Where("email = ? AND token = ? AND expires_at > ?", email, token, time.Now()).
		First(&invite).Error

	if err != nil {
		return nil, err
	}

	return &invite, nil
}

// DeleteInvite removes the invite (e.g. after signup)
func DeleteInvite(id uint) error {
	return db.DB.Delete(&models.Invite{}, id).Error
}

// generateSecureToken generates a secure random token of a specified length
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	return hex.EncodeToString(bytes), err
}
