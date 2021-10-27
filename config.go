package main

type Config struct {
	Server    string `json:"server"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	From      string `json:"from"`
	To        string `json:"to"`
	Chunksize int64  `json:"chunksize"`
}
