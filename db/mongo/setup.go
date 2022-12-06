package mongo

import (
	"context"
	"sync"
	"time"

	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const initTimeout time.Duration = 30 * time.Second

// singleton instance of mongo client
var mongoClient *mongodb.Client
var isInitialized bool = false
var lock sync.Mutex

// getMongoClient returns mongoClient singleton
func getMongoClient() (*mongodb.Client, error) {
	// check lock check
	if !isInitialized {
		lock.Lock()
		defer lock.Unlock()

		if !isInitialized {
			if err := initClient(); err != nil {
				return nil, err
			}
			
			isInitialized = true
			return mongoClient, nil
		}
	}

	return mongoClient, nil
}

func initClient() error {
	// read config	
	uri := "mongodb+srv://dioji.64lmttz.mongodb.net/?authSource=%24external&authMechanism=MONGODB-X509&retryWrites=true&w=majority&tlsCertificateKeyFile=X509-cert-6416201640233700126.pem"

	// mongodb.SessionContext
	initContext, cancel := context.WithTimeout(context.Background(), initTimeout)
	defer cancel()

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
  	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)

	var connErr error
	mongoClient, connErr = mongodb.Connect(initContext, clientOptions)
	if connErr != nil {
		return nil
	}

	return mongoClient.Ping(initContext, nil);
}