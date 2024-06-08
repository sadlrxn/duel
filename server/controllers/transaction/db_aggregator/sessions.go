package db_aggregator

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/sync/syncmap"
	"gorm.io/gorm"
)

// UUID is used for identifying session.
// All the functions in the db-aggregator package, have one more
// optional parameter, session.
// The session is for calling that function as a part of
// transaction, following all-or-nothing strategy.
type UUID uuid.UUID

// UUID_NIL is a reserved constant for nil value of UUID type.
var UUID_NIL = UUID(uuid.Nil)

// sessions is a mapping from UUID to session db.
// The main session pointer is indexed by UUID_NIL key.
var sessions = syncmap.Map{}

// Returns main session pointer id.
func MainSessionId() UUID {
	return UUID_NIL
}

// This is a helper function to get proper session id from session id
// array. It returns main session id for empty parameter.
func getFirstSessionId(sessionId ...UUID) UUID {
	var sessionIdChecked = MainSessionId()
	if len(sessionId) > 0 {
		sessionIdChecked = sessionId[0]
	}
	return sessionIdChecked
}

// Checks whether the db pointer is valid or not.
func isValidSession(sessionId ...UUID) bool {
	if _, prs := sessions.Load(getFirstSessionId(sessionId...)); !prs {
		return false
	}
	return true
}

// Gets the db pointer for the given session id.
// If the sessionId is not provided, it returns the main db pointer.
func getSession(sessionId ...UUID) (*gorm.DB, error) {
	if !isValidSession(sessionId...) {
		return nil, errors.New("db point is not initialized")
	}

	session, ok := sessions.Load(getFirstSessionId(sessionId...))
	if !ok {
		return nil, nil
	}
	return session.(*gorm.DB), nil
}

// Generates a new uuid.
func generateUUID() UUID {
	return UUID(uuid.New())
}

// @External
// Creates a session db object and returns corresponding UUID.
func startSession() (UUID, error) {
	mainSession, err := getSession()
	if err != nil {
		return UUID_NIL, nil
	}

	new_uuid := generateUUID()
	sessions.Store(new_uuid, mainSession.Begin())

	return new_uuid, nil
}

// @External
// Commit a session to main db and remove the saved one from
// map.
func commitSession(sessionId UUID) error {
	session, err := getSession(sessionId)
	if err != nil {
		return err
	}

	sessions.Delete(sessionId)
	return session.Commit().Error
}

// @External
// Remove a session from the memory.
func removeSession(sessionId UUID) error {
	if sessionId == MainSessionId() {
		return errors.New("cannot remove main session")
	}
	session, prs := sessions.Load(sessionId)
	if !prs {
		return nil
	}

	result := (session.(*gorm.DB)).Rollback()
	sessions.Delete(sessionId)
	return result.Error
}
