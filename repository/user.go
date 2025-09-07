package repository

import (
	"context"
	"errors"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepo struct {
	Database *mongo.Database
}

type UserAuth struct {
	Username string
	UserID   bson.ObjectID
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

func (r *UserRepo) AuthCheck(ctx context.Context, sessionToken string, csrfToken string) (*UserAuth, error) {
	var user models.User

	err := r.Database.Collection("users").FindOne(ctx, bson.M{
		"sessions": bson.M{
			"$elemMatch": bson.M{
				"session_token": sessionToken,
				"csrf_token":    csrfToken,
			},
		},
	}).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		logrus.Error(err.Error())
		return nil, err
	}

	userAuth := &UserAuth{
		Username: user.Username,
		UserID:   user.ID,
	}

	return userAuth, nil
}
