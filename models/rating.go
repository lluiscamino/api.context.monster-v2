package models

type Rating struct {
	Tweet       *Tweet   `json:"tweet,omitempty" gorm:"foreignkey:TweetID"`
	TweetID     uint     `json:"tweet_id,omitempty" gorm:"primary_key;auto_increment:false"`
	Keyword     *Keyword `json:"keyword,omitempty" gorm:"association_autoupdate:false;association_autocreate:false;foreignkey:KeywordText;"`
	KeywordText string   `json:"keyword_text,omitempty" gorm:"primary_key"`
	Rate        float64  `json:"rate" gorm:"type:decimal"`
}

/*func (r *Rating) Create() {
	GetDB().Create(r)
}*/
