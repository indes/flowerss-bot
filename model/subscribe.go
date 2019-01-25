package model

import "errors"

type Subscribe struct {
	ID       uint `gorm:"primary_key";"AUTO_INCREMENT"`
	UserID   int64
	SourceID uint
	EditTime
}

func RegistFeed(userID int64, sourceID uint) error {
	var subscribe Subscribe
	db := getConnect()
	defer db.Close()

	if err := db.Where("user_id=? and source_id=?", userID, sourceID).Find(&subscribe).Error; err != nil {
		if err.Error() == "record not found" {
			subscribe.UserID = userID
			subscribe.SourceID = sourceID
			err := db.Create(&subscribe).Error
			if err == nil {
				return nil
			}
		}
		return err
	}
	return nil
}

func GetSubscribeByUserID(userID int64) []Source {
	db := getConnect()
	defer db.Close()
	user := FindOrInitUser(userID)
	return user.Source
}

func GetSubscriberBySource(s *Source) []Subscribe {
	db := getConnect()
	defer db.Close()
	var subs []Subscribe

	db.Where("source_id=?", s.ID).Find(&subs)
	return subs
}

func UnsubByUserIDAndSource(userID int, source *Source) error {
	db := getConnect()
	defer db.Close()
	var sub Subscribe
	db.Where("user_id=? and source_id=?", userID, source.ID).First(&sub)
	if sub.UserID != int64(userID) {
		return errors.New("未订阅该RSS源")
	}
	db.Delete(&sub)
	if source.GetSubscribeNum() < 1 {
		source.DeleteDueNoSubscriber()
	}
	return nil
}

func GetSubsByUserID(userID int) []Subscribe {
	db := getConnect()
	defer db.Close()

	var subs []Subscribe

	db.Where("user_id=?", userID).Find(&subs)

	return subs
}

func UnsubByUserIDAndSourceURL(userID int, url string) error {
	source, err := GetSourceByUrl(url)
	if err != nil {
		return err
	}
	err = UnsubByUserIDAndSource(userID, source)
	return err
}
