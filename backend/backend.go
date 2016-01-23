package backend

import (
	"fmt"
	"github.com/gophergala2016/gobench/backend/model"
	"labix.org/v2/mgo"
	"log"
	"time"
)

// DatabaseConfig holds MongoDB connection params
type databaseConfig struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Name string `json:"name"`
}

// Config holds backend configuration params
type Config struct {
	Mongo databaseConfig `json:"mongo"`
	Debug bool           `json:"debug"`
}

// Backend provides single point of access to business layer
type Backend struct {
	session  *mgo.Session
	log      *log.Logger
	dbConfig databaseConfig

	Model *model.Model
}

// New creates Backend instance, connects to database and initialise caches
func New(cfg *Config, l *log.Logger) (*Backend, error) {

	var err error
	b := &Backend{log: l, dbConfig: cfg.Mongo}

	a := 1
	for {
		err = b.connectDB()
		if err == nil {
			break
		}
		if cfg.Debug {
			return nil, err
		}
		log.Printf("Database connection error! Attempt: %d", a)
		a++
		time.Sleep(3 * time.Second)
	}

	b.Model, err = model.New(b.session.DB(b.dbConfig.Name), l)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Start launches background processes
func (b *Backend) Start() error {

	return nil
}

func (b *Backend) connectDB() error {

	var err error

	b.session, err = mgo.DialWithTimeout(fmt.Sprintf("%s:%d/%s", b.dbConfig.IP, b.dbConfig.Port, b.dbConfig.Name), time.Second*3)
	if err != nil {
		return err
	}
	b.session.SetMode(mgo.Monotonic, true)

	return b.session.Ping()
}
