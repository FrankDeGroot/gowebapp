# Go Web App

This app is my learning playground for Go, CockroachDB, RedPanda and Web Components.

It is a simple to-do app to maintain a list of task descriptions and a checkbox to sign them off.

It features
- [CockroachDB](https://github.com/cockroachdb/cockroach) to store tasks,
- an HTTP/JSON API to get, post, put and delete tasks,
- [RedPanda](https://github.com/redpanda-data/redpanda/) to track changes to tasks,
- a WebSocket connection to receive changes,
- a static file server for the web site,
- Web Components to show and edit the tasks and show the saving status.