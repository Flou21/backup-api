package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Backup struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	Target    *Target            `json:"target" bson:"target"`
	Path      string             `json:"path" bson:"path"`
	Size      int64              `json:"size" bson:"size"`
}

func ListBackups(targetName string) ([]*Backup, error) {

	filter := bson.D{{"target.name", targetName}}

	result, err := backupsCollection.Find(context.TODO(), filter)

	if err != nil {
		log.Error().Err(err).Msgf("there was an error when listing backups in mongo")
		return nil, err
	}

	var backups []*Backup
	err = result.All(context.TODO(), &backups)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error decoding backups")
		return nil, err
	}

	return backups, nil
}

func (b *Backup) Insert() error {
	b.ID = primitive.NewObjectID()

	res, err := backupsCollection.InsertOne(context.TODO(), b)
	if err != nil {
		return err
	}

	log.Info().Msgf("inserted backup, insertId: %s", res.InsertedID)
	return nil
}

func (b *Backup) Remove() error {
	return fmt.Errorf("not implemented yet")
}

func LatestBackupForTarget(target *Target) (*Backup, error) {

	matchStage := bson.D{{"$match", bson.D{{"target.name", "test-db"}}}}
	sortStage := bson.D{{"$sort", bson.D{{"timestamp", -1}}}}
	limitStage := bson.D{{"$limit", 1}}

	cursor, err := backupsCollection.Aggregate(context.TODO(), mongo.Pipeline{
		matchStage,
		sortStage,
		limitStage,
	})

	if err != nil {
		return nil, err
	}

	backups := []Backup{}
	err = cursor.All(context.TODO(), &backups)
	if err != nil {
		return nil, err
	}

	if len(backups) == 0 {
		return nil, nil
	}

	latestBackup := &backups[0]

	log.Info().Msgf("found backup from %s as latest backup for target %s", latestBackup.Timestamp, target)

	return latestBackup, nil
}

func ListAllBackups() ([]*Backup, error) {

	targets, err := ListTargets()
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when listing targets")
	}

	backups := make(chan *Backup)
	var result []*Backup

	go func() {
		var wg sync.WaitGroup
		for _, target := range targets {
			wg.Add(1)
			go func(target *Target) {
				defer wg.Done()
				tempResut, err := ListBackups(target.Name)
				if err != nil {
					log.Error().Err(err).Msgf("there was an error when listing backups for target %s", target.Name)
				}
				for _, b := range tempResut {
					backups <- b
				}
			}(target)
		}
		wg.Wait()
		close(backups)
	}()

	for b := range backups {
		result = append(result, b)
	}

	return result, nil
}
