package store
import(
	"sync"
	"github.com/Arjun8242/ticket-system.git/models"
)

type Store struct {
	mu sync.RWMutex
	Users map[string]*models.User
	Tickets map[string]*models.Ticket
}

func NewStore() *Store {
	return &Store{
		Users: make(map[string]*models.User),
		Tickets: make(map[string]*models.Ticket),
	}
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range s.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (s *Store) CreateUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Users[user.ID] = user
	return nil
}