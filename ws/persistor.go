package ws

import (
	"time"

	"github.com/cheersappio/matchserver/models"
	"github.com/cheersappio/matchserver/utils"
	"gopkg.in/mgo.v2/bson"
)

const (
	// max distance for the query in meters
	maxDistance = 5
	// seconds range to do match
	secondsRange = 3
)

type persistor struct {
	persist  chan *postStroke
	response chan []string
}

func createPersistor(response chan []string) *persistor {
	return &persistor{
		persist:  make(chan *postStroke),
		response: response,
	}
}

func createStrokeFrom(postStrokeVar *postStroke) models.Stroke {
	return models.Stroke{
		Location:  postStrokeVar.Loc,
		UserID:    postStrokeVar.userID,
		CreatedAt: time.Now(),
	}
}

func (p *persistor) run() {
	defer close(p.persist)
	postStrokeVar := <-p.persist
	stroke := createStrokeFrom(postStrokeVar)
	if err := p.save(&stroke); err != nil {
		return
	}
	nearUsers, err := p.findNear(&stroke)
	if err != nil {
		panic(err)
	}
	p.response <- nearUsers
}

func (p *persistor) save(stroke *models.Stroke) error {
	utils.Log.Infof("Persisting %s stroke: %v", stroke.UserID, stroke)
	return models.StrokesCollection.Insert(stroke)
}

func (p *persistor) findNear(stroke *models.Stroke) ([]string, error) {
	results := []models.Stroke{}
	query := bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": stroke.Location,
				},
				"$maxDistance": maxDistance,
			},
		},
		"user_id": bson.M{
			"$ne": stroke.UserID,
		},
		"created_at": bson.M{
			"$gte": time.Now().Add(-1 * secondsRange * time.Second),
			"$lte": time.Now().Add(secondsRange * time.Second),
		},
	}
	err := models.StrokesCollection.Find(query).All(&results)
	utils.Log.Infof("Query executed by %s: %v", stroke.UserID, query)
	utils.Log.Infof("Actor %s found matches %d", stroke.UserID, len(results))
	names := make([]string, len(results))
	for i, r := range results {
		names[i] = r.UserID
	}
	return names, err
}
