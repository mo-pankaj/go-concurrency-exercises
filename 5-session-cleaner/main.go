//////////////////////////////////////////////////////////////////////
//
// Given is a SessionManager that stores session information in
// memory. The SessionManager itself is working, however, since we
// keep on adding new sessions to the manager our program will
// eventually run out of memory.
//
// Your task is to implement a session cleaner routine that runs
// concurrently in the background and cleans every session that
// hasn't been updated for more than 5 seconds (of course usually
// session times are much longer).
//
// Note that we expect the session to be removed anytime between 5 and
// 7 seconds after the last update. Also, note that you have to be
// very careful in order to prevent race conditions.
//

package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// SessionManager keeps track of all sessions from creation, updating
// to destroying.
type SessionManager struct {
	mutex    *sync.Mutex
	sessions map[string]Session
}

// Session stores the session's data
type Session struct {
	Data map[string]interface{}
}

// NewSessionManager creates a new sessionManager
func NewSessionManager() *SessionManager {
	m := &SessionManager{
		mutex:    &sync.Mutex{},
		sessions: make(map[string]Session),
	}
	go m.cleanSession()

	return m
}

func (m *SessionManager) cleanSession() {
	t := time.NewTicker(5 * time.Second)

	f := func() {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		for k, v := range m.sessions {
			updateTime, ok := v.Data["updated"]
			if ok {
				upd := updateTime.(time.Time)
				fmt.Println("Session", k, "last updated at", time.Since(upd))
				if time.Since(upd) > 5*time.Second {
					fmt.Println("Session", k, "expired")
					delete(m.sessions, k)
				}
			}
		}
	}

	for {
		select {
		case <-t.C:
			f()
		}
	}
}

// CreateSession creates a new session and returns the sessionID
func (m *SessionManager) CreateSession() (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	sessionID, err := MakeSessionID()
	if err != nil {
		return "", err
	}
	data := make(map[string]interface{})
	data["updated"] = time.Now()

	m.sessions[sessionID] = Session{
		Data: data,
	}

	return sessionID, nil
}

// ErrSessionNotFound returned when sessionID not listed in
// SessionManager
var ErrSessionNotFound = errors.New("SessionID does not exists")

// GetSessionData returns data related to session if sessionID is
// found, errors otherwise
func (m *SessionManager) GetSessionData(sessionID string) (map[string]interface{}, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session.Data, nil
}

// UpdateSessionData overwrites the old session data with the new one
func (m *SessionManager) UpdateSessionData(sessionID string, data map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	data["updated"] = time.Now()
	_, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	// Hint: you should renew expiry of the session here
	m.sessions[sessionID] = Session{
		Data: data,
	}

	return nil
}

func main() {
	// Create new sessionManager and new session
	m := NewSessionManager()
	sID, err := m.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created new session with ID", sID)

	// Update session data
	data := make(map[string]interface{})
	data["website"] = "longhoang.de"

	err = m.UpdateSessionData(sID, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Update session data, set website to longhoang.de")

	// Retrieve data from manager again
	updatedData, err := m.GetSessionData(sID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get session data:", updatedData)
}
