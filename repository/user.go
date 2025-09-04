package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepo struct {
	Database *mongo.Database
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.Database.Collection("users").InsertOne(ctx, user)
	return err
}

var ErrUserNotFound = errors.New("user not found")

func (r *UserRepo) GetUser(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.Database.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, err
}

func (r *UserRepo) DoesExist(ctx context.Context, username string) bool {
	res := r.Database.Collection("users").FindOne(ctx, bson.M{"username": username})
	return !errors.Is(res.Err(), mongo.ErrNoDocuments)
}

func (r *UserRepo) AddLoginSession(ctx context.Context, username string, session models.UserSession) error {
	update := bson.M{
		"$push": bson.M{
			"sessions": session,
		},
	}

	_, err := r.Database.Collection("users").UpdateOne(ctx, bson.M{"username": username}, update)
	return err
}

func (r *UserRepo) FinalizeSession(ctx context.Context, username string, sessionToken string) error {
	update := bson.M{
		"$pull": bson.M{
			"sessions": bson.M{
				"session_token": sessionToken,
			},
		},
	}

	_, err := r.Database.Collection("users").UpdateOne(ctx, bson.M{"username": username}, update)
	return err
}

func (r *UserRepo) AuthCheck(ctx context.Context, username string, sessionToken string, csrfToken string) bool {
	fmt.Println(username)
	fmt.Println(sessionToken)
	fmt.Println(csrfToken)
	res := r.Database.Collection("users").FindOne(ctx, bson.M{
		"username":               username,
		"sessions.session_token": sessionToken,
		"sessions.csrf_token":    csrfToken,
	})

	return !errors.Is(res.Err(), mongo.ErrNoDocuments)
}
