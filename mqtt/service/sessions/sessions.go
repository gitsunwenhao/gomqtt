package sessions

/* Sessions is used to provide session for user in mqtt server
   Author- Sunface*/
import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
)

var (
	// ErrSessionsProviderNotFound means no providers are available
	ErrSessionsProviderNotFound = errors.New("Session: Session provider not found")

	// ErrKeyNotAvailable means the key is invalid
	ErrKeyNotAvailable = errors.New("Session: not item found for key.")

	providers = make(map[string]Provider)
)

// Provider is a session provider, you can write your own provider by implementing this interface
type Provider interface {
	New(id string) (*Session, error)
	Get(id string) (*Session, error)
	Del(id string)
	Count() int
	Close() error
}

// Register makes a session provider available by the provided name.
// If a Register is called twice with the same name or if the driver is nil,
// it panics.
func Register(name string, provider Provider) {
	if provider == nil {
		log.Panicln("[FATAL] session: Register provide is nil")
	}

	if _, dup := providers[name]; dup {
		log.Panicln("[FATAL] session: Register called twice for provider " + name)
	}

	providers[name] = provider
}

// Unregister will delete a provider
func Unregister(name string) {
	delete(providers, name)
}

// Manager is used for session providers
type Manager struct {
	p Provider
}

// NewManager create a new session provider manager from the given providername
func NewManager(pn string) (*Manager, error) {
	p, ok := providers[pn]
	if !ok {
		return nil, fmt.Errorf("session: unknown provider %q", pn)
	}

	return &Manager{p: p}, nil
}

// New creates a session
func (m *Manager) New(id string) (*Session, error) {
	if id == "" {
		id = m.sessionID()
	}
	return m.p.New(id)
}

// Get reuturn the specify session
func (m *Manager) Get(id string) (*Session, error) {
	return m.p.Get(id)
}

// Del delete the specify session
func (m *Manager) Del(id string) {
	m.p.Del(id)
}

// Count returns the session conunt of the provider
func (m *Manager) Count() int {
	return m.p.Count()
}

// Close close a provider
func (m *Manager) Close() error {
	return m.p.Close()
}

func (m *Manager) sessionID() string {
	b := make([]byte, 15)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
