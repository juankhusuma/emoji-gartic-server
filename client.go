package main

import (
	ws "github.com/gofiber/websocket/v2"
)

type Client struct {
	hub  *Hub
	conn *ws.Conn
	send chan []byte
}
