package model

import "time"

type TopicInfo struct {
	ID         string    `gorm:"size:20;primary_key;not null" json:"id"`
	Link       string    `gorm:"size:256;not null" json:"link"`
	Title      string    `gorm:"size:100;not null" json:"title"`
	Createtime time.Time `gorm:"column:create_time;type:time;not null" json:"create_time"`
}
