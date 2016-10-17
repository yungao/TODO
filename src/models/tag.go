package models

import (
	"config"
	"github.com/coopernurse/gorp"
	"log"
	"time"
	"fmt"
	"errors"
)

type Tag struct {
	ID        int    `db:"id"        json:"id"`
	CreatorID int    `db:"uid"       json:"uid"`
	Name      string `db:"name"      json:"name"        form:"name"     binding:"required"`
	Remark    string `db:"remark"    json:"remark"      form:"remark"`
	CreateAt  int64  `db:"create"    json:"create"`
	/* Active:
	 *       -1: deleted
	 *       1:  normal
	 */
	Active int8 `db:"active"    json:"active"`
}

type Tags struct {
	Collection []Tag `json:"tags"`
}

// Create tag table if not exist
func CreateTagTable(db *gorp.DbMap) {
	tb := db.AddTableWithName(Tag{}, config.TABLE_NAME_TAGS)
	tb.SetKeys(true, "id")
	tb.ColMap("uid").SetNotNull(true)
	tb.ColMap("name").SetMaxSize(20).SetNotNull(true).SetUnique(true)
	// tb.SetUniqueTogether("uid", "name")
	tb.ColMap("remark").SetMaxSize(255)
	tb.ColMap("create").SetNotNull(true)
	tb.ColMap("active").SetNotNull(true)

	err := db.CreateTablesIfNotExists()
	// db.DropTables()
	if err != nil {
		panic(err)
	}

	log.Printf(">>> Table[%s] created", config.TABLE_NAME_TAGS)
}

func (tag *Tag) PreInsert(s gorp.SqlExecutor) error {
	tag.CreateAt = time.Now().Unix()
	tag.Active = 1
	return nil
}

func (tag *Tag) String() string {
    return fmt.Sprintf("{ID:%d, CreatorID:%d, Name:%s, Remark:%s, CreateAt:%d, Active:%d}", tag.ID, tag.CreatorID, tag.Name, tag.Remark, tag.CreateAt, tag.Active)
}

/*
* Query Tag from database by ID
 */
func GetTagByID(db *gorp.DbMap, id int) (*Tag, error) {
	ret, err := db.Get(Tag{}, id)
	if err != nil {
		return nil, err
	}

	if tag, ok := ret.(*Tag); ok {
		return tag, nil
	}

	return nil, errors.New("Tag does not exist!")
}

// func (tag *Tag) PreUpdate(s gorp.SqlExecutor) error {
// 	tag.UpdateAt = time.Now().Unix()
// 	return nil
// }

type TodoTag struct {
	ID     int `db:"id"        json:"id"`
	TodoID int `db:"todoid"    json:"todoid"     form:"todoid"  binding:"required"`
	TagID  int `db:"tagid"     json:"tagid"      form:"tagid"   binding:"required"`
	Tag    *Tag `db:"-"     json:"tag"`
	CreatorID   int    `db:"uid"       json:"-"`
	/* Active:
	 *       -1: deleted
	 *       1:  normal
	 */
	Active int8 `db:"active"    json:"active"`
}

type TodoTags struct {
	Collection []Tag `json:"tags"`
}

// Create todo tag table if not exist
func CreateTodoTagTable(db *gorp.DbMap) {
	tb := db.AddTableWithName(TodoTag{}, config.TABLE_NAME_TODO_TAGS)
	tb.SetKeys(true, "id")
	tb.ColMap("todoid").SetNotNull(true)
	tb.ColMap("tagid").SetNotNull(true)
	tb.SetUniqueTogether("todoid", "tagid")
	tb.ColMap("active").SetNotNull(true)

	err := db.CreateTablesIfNotExists()
	// db.DropTables()
	if err != nil {
		panic(err)
	}

	log.Printf(">>> Table[%s] created", config.TABLE_NAME_TODO_TAGS)
}

func (ttag *TodoTag) PreInsert(s gorp.SqlExecutor) error {
	ttag.Active = 1
	return nil
}

func (tag *TodoTag) String() string {
    return fmt.Sprintf("{ID:%d, TodoID:%d, TagID:%d, Active:%d}", tag.ID, tag.TodoID, tag.TagID, tag.Active)
}
// func (ttag *TodoTag) PreUpdate(s gorp.SqlExecutor) error {
// 	return nil
// }

func GetTagsByTodoID(db *gorp.DbMap, tid int) ([]*TodoTag, error) {
    var tags []*TodoTag
	_, err := db.Select(tags, fmt.Sprintf("SELECT * FROM %s WHERE %s.todoid=%d", config.TABLE_NAME_TODO_TAGS, config.TABLE_NAME_TODO_TAGS, tid))
    if err == nil {
        for _, tag := range tags {
            tag.Tag, err = GetTagByID(db, tag.ID)
        }
    }

    return tags, err
}
