package postgres

import (
	"context"
	"fmt"
	"log"
	"social/internal/entity"
	"testing"
)

func TestQuery(t *testing.T) {
	pg, err := New(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		"postgres", "123456", "127.0.0.1", 5433, "social"),
		MaxPoolSize(2))
	if err != nil {
		panic(err)
	}
	sql, _, err := pg.Builder.
		Select("id", "username", "password", "nickname", "sex", "avatar").
		Where("username=?").
		From("users").ToSql()
	if err != nil {
		log.Fatal(err)
	}
	//s := `SELECT id, username, password, nickname, sex, avatar FROM users WHERE username=$1`
	e := &entity.User{}
	err = pg.Pool.QueryRow(context.Background(), sql, "zs").Scan(&e.ID, &e.Username, &e.Password, &e.Nickname, &e.Sex, &e.Avatar)
	t.Logf("%+v, err: %v \n", e, err)
}

func TestList(t *testing.T) {
	pg, err := New(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		"postgres", "123456", "127.0.0.1", 5433, "social"),
		MaxPoolSize(2))
	if err != nil {
		panic(err)
	}

	sql, _, err := pg.Builder.Select("id", "username", "password", "nickname", "sex", "avatar").From("users").ToSql()
	if err != nil {
		t.Fatal(err)
	}
	rows, err := pg.Pool.Query(context.Background(), sql)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	type User struct {
		ID          int64  `json:"id"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		Nickname    string `json:"nickname"`
		Sex         int32  `json:"sex"`
		Salt        string `json:"salt"`
		Avatar      string `json:"avatar"`
		Summary     string `json:"summary"`
		AccessToken string `json:"access_token"`
		ExpireTime  int64  `json:"expire_time"`
		CreatedTime int64  `json:"created_time"`
		UpdatedTime int64  `json:"updated_time"`
	}

	entities := make([]User, 0, 64)
	for rows.Next() {
		e := User{}
		err = rows.Scan(&e.ID, &e.Username, &e.Password, &e.Nickname, &e.Sex, &e.Avatar)
		if err != nil {
			t.Fatal(err)
		}
		entities = append(entities, e)
	}
	t.Logf("%+v", entities)
}
