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

func (s *Store) CreateTicket(ticket *models.Ticket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Tickets[ticket.ID] = ticket
	return nil
}

func (s *Store) GetTicketsByOwner(ownerID string) []*models.Ticket {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]*models.Ticket, 0)
	for _, ticket := range s.Tickets {
		if ticket.OwnerID == ownerID {
			result = append(result, ticket)
		}
	}
	return result
}

func (s *Store) GetTicketByID(id string) (*models.Ticket, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ticket, exists := s.Tickets[id]
	return ticket, exists
}