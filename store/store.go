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