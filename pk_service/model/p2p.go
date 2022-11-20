package model

import "time"

type MatchRecords struct {
	p1        int32     `gorm:"type:int;index"` //player1
	score1    int32     `gorm:"type:int"`
	p2        int32     `gorm:"type:int;index"` // player2
	score2    int32     `gorm:"type:int"`
	CreateAt  time.Time // 创建时间
	EndTime   time.Time // 结束时间
	winner_id int32     // 其中的一个
}

// TODO: 历届排行榜
