package main

type DatabaseConfig struct {
}

// Key: topic
// value: {isActive: bool, group: string}

type GroupEntity struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type Topic struct {
	Name  string
	Group []GroupEntity
}

type Database interface {
	getActiveTopics() []Topic
	listAllTopics() []Topic
	getAllGroupsByTopic(topic string) Topic
	toggleTopic(topic, group string)
	setNewTopicGroup()
}
