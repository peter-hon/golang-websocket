package changes

import (
	rethinkcore "bitbucket.org/chaatzltd/webcontent/core/rethink"

	"github.com/Sirupsen/logrus"
	r "gopkg.in/dancannon/gorethink.v2"
)

var session *r.Session
var database = "chaatz"

func GetSession() *r.Session {
	session = rethinkcore.GetSession()
	return session
}

func ActiveChanges(ch chan interface{}) {
	log := logrus.WithFields(logrus.Fields{"module": "changes", "method": "ActiveChanges"})
	GetSession()
	go func() {
		for {
			res, err := r.DB("chaatz").Table("broadcastmessage").Changes().Run(session)
			if err != nil {
				log.WithField("error", err).Error("Failed to get changes from rethinkdb")
			}

			var response interface{}
			for res.Next(&response) {
				ch <- response
			}

			if res.Err() != nil {
				log.WithField("error", err).Error("Failed to get changes from rethinkdb")
			}
		}
	}()
}
