package dao

type Lyric struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Type  int // 1 qqmusic
	Lyric []byte
}

func GetLyric(id string, _type int) *Lyric {
	var entity Lyric
	Db.Where("id = ? and type = ?", id, _type).First(&entity)
	return &entity
}

func GetLyricById(id string) *Lyric {
	var entity Lyric
	Db.First(&entity, id)
	return &entity
}

func SaveLyric(id string, lyric []byte, name string, _type int) {
	Db.Save(&Lyric{ID: id, Lyric: lyric, Name: name, Type: _type})
}

func SaveLyricEntity(entity *Lyric) {
	Db.Save(&entity)
}

// 查询单个id
// func (entity *Lyric) GetStr(id string) *Lyric {
// 	if id != "" {
// 		_id, _ := strconv.Atoi(id)
// 		db.First(&entity, _id)
// 	}
// 	return entity
// }
