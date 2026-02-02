package auth

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ra341/glacier/internal/user"
	"github.com/ra341/glacier/pkg/syncmap"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestService_RoleChecks(t *testing.T) {
	uts := &TestUserStore{}
	us := user.NewService(uts)
	ts := &TestSessionStore{}
	srv := New(ts, us, false)

	u, err := srv.userSrv.GetByUsername(user.DefaultUser)
	require.NoError(t, err)

	use := "new"
	p := "new"

	err = srv.Register(use, p, user.Magos, &u)
	require.NoError(t, err)

	adminUser, err := uts.GetByUsername(use)
	require.NoError(t, err)

	use = "new22"
	p = "new333"
	err = srv.Register(use, p, user.Omnissiah, &u)
	require.NoError(t, err)

	err = srv.Register(use, p, user.Omnissiah, &adminUser)
	require.Error(t, err)

	err = srv.Register(use, p, user.TechPriest, &adminUser)
	require.NoError(t, err)
}

func TestService_Register(t *testing.T) {
	uts := &TestUserStore{}
	us := user.NewService(uts)
	ts := &TestSessionStore{}
	srv := New(ts, us, false)

	u := "test"
	p := "test"

	err := srv.Register(u, p, user.Magos, nil)
	require.Error(t, err)

	srv.openRegistration = true
	err = srv.Register(u, p, user.Magos, nil)
	require.NoError(t, err)

	usd, err := uts.GetByUsername(u)
	require.NoError(t, err)
	require.Equal(t, usd.Role, user.TechPriest)

	expectedSess, session, refresh, err := srv.Login(u, p, Web)
	require.NoError(t, err)

	t.Log(session, refresh)

	verifySession, err := srv.VerifySession(session)
	require.NoError(t, err)
	require.Equal(t, expectedSess, verifySession, "session changed between login and verify")

	// invalidate session
	expectedSess.SessionTokenExpiry = time.Now().Add(-time.Second * 5)
	err = ts.Edit(&expectedSess)
	require.NoError(t, err)

	verifySession, err = srv.VerifySession(session)
	require.ErrorIs(t, err, ErrTokenExpired, "token should be expired")

	_, session, _, err = srv.RefreshSession(&expectedSess)
	require.NoError(t, err)

	_, err = srv.VerifySession(session)
	require.NoError(t, err)

	expectedSess.RefreshTokenExpiry = time.Now().Add(-time.Second * 5)
	err = ts.Edit(&expectedSess)
	require.NoError(t, err)

	_, session, _, err = srv.RefreshSession(&expectedSess)
	require.ErrorIs(t, err, ErrTokenExpired, "token should be expired")
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// test user
type TestUserStore struct {
	store   syncmap.Map[uint, user.User]
	idCount atomic.Int64
}

func (t *TestUserStore) GetByUsername(username string) (user.User, error) {
	var found user.User
	var ok bool
	t.store.Range(func(key uint, value user.User) bool {
		if value.Username == username {
			found = value
			ok = true
			return false
		}
		return true
	})
	if !ok {
		return user.User{}, errors.New("user not found")
	}
	return found, nil
}

func (t *TestUserStore) GetByID(id uint) (user.User, error) {
	val, ok := t.store.Load(id)
	if !ok {
		return user.User{}, gorm.ErrRecordNotFound
	}
	return val, nil
}

func (t *TestUserStore) New(u *user.User) error {
	newID := uint(t.idCount.Add(1))
	u.ID = newID
	t.store.Store(newID, *u)
	return nil
}

func (t *TestUserStore) Edit(u *user.User) error {
	if _, ok := t.store.Load(u.ID); !ok {
		return errors.New("user does not exist")
	}
	t.store.Store(u.ID, *u)
	return nil
}

func (t *TestUserStore) Delete(id uint) error {
	t.store.Delete(id)
	return nil
}

func (t *TestUserStore) List() ([]user.User, error) {
	var users []user.User
	t.store.Range(func(key uint, value user.User) bool {
		users = append(users, value)
		return true
	})
	return users, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// session

type TestSessionStore struct {
	store   syncmap.Map[uint, Session]
	idCount atomic.Int64
}

func (t *TestSessionStore) New(session *Session) error {
	newID := uint(t.idCount.Add(1))
	session.ID = newID
	t.store.Store(newID, *session)
	return nil
}

func (t *TestSessionStore) Edit(session *Session) error {
	if _, ok := t.store.Load(session.ID); !ok {
		return errors.New("session does not exist")
	}
	t.store.Store(session.ID, *session)
	return nil
}

func (t *TestSessionStore) Delete(session *Session) error {
	t.store.Delete(session.ID)
	return nil
}

func (t *TestSessionStore) List() ([]Session, error) {
	var sessions []Session
	t.store.Range(func(key uint, value Session) bool {
		sessions = append(sessions, value)
		return true
	})
	return sessions, nil
}

func (t *TestSessionStore) ListByUser(userId uint) ([]Session, error) {
	var sessions []Session
	t.store.Range(func(key uint, value Session) bool {
		if value.UserId == userId {
			sessions = append(sessions, value)
		}
		return true
	})
	return sessions, nil
}

func (t *TestSessionStore) GetBySessionId(sessionId uint) (Session, error) {
	val, ok := t.store.Load(sessionId)
	if !ok {
		return Session{}, errors.New("session not found")
	}
	return val, nil
}

func (t *TestSessionStore) GetBySessionToken(token string) (Session, error) {
	var found Session
	var ok bool
	t.store.Range(func(key uint, value Session) bool {
		if value.HashedSessionToken == token {
			found = value
			ok = true
			return false
		}
		return true
	})
	if !ok {
		return Session{}, errors.New("session not found")
	}
	return found, nil
}
