package controllers

import (
	mdchanges "bitbucket.org/chaatzltd/webcontent/models/changes"
	wshub "bitbucket.org/chaatzltd/webcontent/services/hub"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	Controller
}

func (this *WebSocketController) Join() {
	log := logrus.WithFields(logrus.Fields{"controller": "broadcast", "controller_method": "Get"})
	h := wshub.NewHub()
	go h.Run()
	mdchanges.ActiveChanges(h.Broadcast)

	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		log.WithFields(logrus.Fields{"err": err}).Error("Failed to upgrade to websocket")
		this.OutputJson(500, nil)
		return
	} else if err != nil {
		log.WithFields(logrus.Fields{"err": err}).Error("Failed to upgrade to websocket")
		this.OutputJson(500, nil)
		return
	}

	c := &wshub.Connection{Send: make(chan interface{}, 256), Ws: ws}
	h.Register <- c
	defer func() { h.Unregister <- c }()
	go c.Writer()
	c.Reader()
}
