// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	users syncmap.Map

	clients syncmap.Map

	// Inbound messages from the clients.
	EventEmitter chan types.WSEvent

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		EventEmitter: make(chan types.WSEvent, 4096),
		register:     make(chan *Client, 256),
		unregister:   make(chan *Client, 256),
		users:        syncmap.Map{},
		clients:      syncmap.Map{},
	}
}

func (h *Hub) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.LogMessage("hub", "recovered", "info", logrus.Fields{})
		}
	}()
	for {
		select {
		case client := <-h.register:
			h.clients.Store(client.conn, client)
			if client.userID != nil {
				h.users.Store(*client.userID, client)
				controllers.Chat.ActivateUser(*client.userID)
			}
		case client := <-h.unregister:
			if _, ok := h.clients.Load(client.conn); ok {
				if client.userID != nil {
					h.users.Delete(*client.userID)
					controllers.Chat.DeactivateUser(*client.userID)
				}
				h.clients.Delete(client.conn)
				close(client.send)
				client.hub = nil
			}
		case wsEvent := <-h.EventEmitter:
			count := len(wsEvent.Users) + len(wsEvent.Conns)
			if count == 0 {
				h.clients.Range(func(key, value any) bool {
					if value.(*Client).room != wsEvent.Room && wsEvent.Room != types.Chat {
						return true
					}
					select {
					case value.(*Client).send <- wsEvent.Message:
					default:
						close(value.(*Client).send)
						if value.(*Client).userID != nil {
							h.users.Delete(value.(*Client).userID)
						}
						h.clients.Delete(key)
					}
					return true
				})
			} else {
				for i := 0; i < len(wsEvent.Users); i++ {
					userID := wsEvent.Users[i]
					if user, prs := h.users.Load(userID); prs {
						select {
						case user.(*Client).send <- wsEvent.Message:
						default:
							close(user.(*Client).send)
							h.users.Delete(userID)
							h.clients.Delete(user.(*Client).conn)
						}
					}
				}
				for i := 0; i < len(wsEvent.Conns); i++ {
					conn := wsEvent.Conns[i]
					if client, prs := h.clients.Load(conn); prs {
						select {
						case client.(*Client).send <- wsEvent.Message:
						default:
							close(client.(*Client).send)
							h.clients.Delete(conn)
							if client.(*Client).userID != nil {
								h.users.Delete(client.(*Client).userID)
							}
						}
					}
				}
			}
		}
	}
}
