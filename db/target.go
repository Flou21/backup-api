package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Target struct {
	ID                     primitive.ObjectID `bson:"_id" json:"id"`
	Type                   string             `bson:"type" json:"type"`
	Ip                     string             `bson:"ip" json:"ip"`
	Port                   int64              `bson:"port" json:"port"`
	Name                   string             `bson:"name" json:"name"`
	Username               string             `bson:"username" json:"username"`
	Password               string             `bson:"password" json:"password"`
	Database               string             `bson:"database" json:"database"`
	Interval               int64              `json:"interval" bson:"interval"`
	AuthenticationDatabase string             `bson:"authentication_database,omitempty" json:"authenticationDatabase"`
}

func ListTargets() ([]*Target, error) {
	filter := bson.D{}

	cursor, err := targetCollection.Find(context.TODO(), filter)

	if err != nil {
		log.Error().Err(err).Msgf("could not retrieve backup targets from database")
		return nil, err
	}

	targets := []*Target{}
	for cursor.Next(context.Background()) {
		el := &Target{}
		err = cursor.Decode(el)
		if err != nil {
			log.Error().Err(err).Msgf("could not decode backup target")
			return nil, err
		}
		targets = append(targets, el)
	}

	return targets, nil
}

func ListTargetsWithPendingBackups(pendingTargets chan<- *Target) error {
	targets, err := ListTargets()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, target := range targets {
		wg.Add(1)

		go func(target *Target, wg *sync.WaitGroup) {
			defer wg.Done()

			b, err := LatestBackupForTarget(target)
			if err != nil {
				log.Error().Err(err).Msgf("failed to load latest backup for target %s", target.Name)
				return
			}

			if b == nil {
				log.Info().Msgf("target %s has no backup so far")
				pendingTargets <- target
				return
			}

			if doesTargetNeedNewBackup(target, b) {
				pendingTargets <- target
			}

		}(target, &wg)
	}

	wg.Wait()
	close(pendingTargets)
	return nil
}

func TargetByName(name string) (*Target, error) {
	if name == "" {
		return nil, fmt.Errorf("target name must not be empty")
	}

	filter := bson.D{{Key: "name", Value: name}}
	result := targetCollection.FindOne(context.TODO(), filter)
	err := result.Err()
	if err != nil {
		log.Error().Err(err).Msgf("an error occurred while searching for target with name %s", name)
		return nil, err
	}

	target := new(Target)
	if err = result.Decode(target); err != nil {
		log.Error().Err(err).Msgf("an error occurred while decoding target with name %s", name)
		return nil, err
	}

	return target, nil
}

func (b *Target) Insert() error {
	b.ID = primitive.NewObjectID()

	res, err := targetCollection.InsertOne(context.TODO(), b)
	if err != nil {
		return err
	}

	log.Info().Msgf("inserted target, insertId: %s", res.InsertedID)
	return nil
}

func (t *Target) Update() error {
	beforeUpdate, err := TargetByName(t.Name)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when searching for target for update, name: %s", t.Name)
		return err
	}

	t.ID = beforeUpdate.ID
	filter := bson.D{{Key: "name", Value: t.Name}}
	_, err = targetCollection.ReplaceOne(context.TODO(), filter, t)

	if err != nil {
		log.Error().Err(err).Msgf("there was an error when updating target with name %s", t.Name)
		return err
	}

	return nil
}

func doesTargetNeedNewBackup(target *Target, latestBackup *Backup) bool {
	currentTimestamp := time.Now().Unix()

	diff := currentTimestamp - latestBackup.Timestamp.Unix()

	if diff > target.Interval {
		log.Info().Msgf("target has a interval of %d, latest backup is from %s, therefore new backup is needed", target.Interval, latestBackup.Timestamp)
		return true
	}

	log.Info().Msgf("target has a interval of %d, latest backup is from %s, therefore no new backup is needed", target.Interval, latestBackup.Timestamp)
	return false
}
