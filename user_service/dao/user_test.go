package dao

import (
	"context"
	"gorm.io/gorm"
	Potesting "server/shared/postgres/testing"
	userpb "server/user_service/api/gen/v1/user"
	"server/user_service/model"
	"testing"
)

func TestUser_CreateUser(t *testing.T) {

}

// 插入
//BaseModel 1~4
//Mobile phone1~4
//PassWord pwd1~4
//NickName name1~4
func insertHelper(c context.Context, t *testing.T, db *gorm.DB) {
	rows := []model.User{
		{
			BaseModel: model.BaseModel{ID: 1},
			Mobile:    "phone1",
			PassWord:  "pwd1",
			NickName:  "name1",
		},
		{
			BaseModel: model.BaseModel{ID: 2},
			Mobile:    "phone2",
			PassWord:  "pwd2",
			NickName:  "name2",
		},
		{
			BaseModel: model.BaseModel{ID: 3},
			Mobile:    "phone3",
			PassWord:  "pwd3.",
			NickName:  "name3",
		},
		{
			BaseModel: model.BaseModel{ID: 4},
			Mobile:    "phone4",
			PassWord:  "pwd4.",
			NickName:  "name4",
		},
	}
	db.AutoMigrate(&model.User{})
	for _, row := range rows {
		db.Create(&row)
	}
	var users []model.User
	find := db.Find(&users)
	if find.RowsAffected != 4 {
		t.Fatal("not found enough")
	}
}
func TestCreate(t *testing.T) {
	cases := []struct {
		name    string
		user    *userpb.UserInfo
		wantErr bool
	}{
		{
			name: "normal_create",
			user: &userpb.UserInfo{
				Id:       10,
				Mobile:   "1234",
				PassWord: "pwd123",
				NickName: "123",
			},
		},
		{
			name: "mobile_repeat",
			user: &userpb.UserInfo{
				Id:       11,
				Mobile:   "phone1",
				PassWord: "pwd123",
				NickName: "123",
			},
			wantErr: true,
		},
		{
			name: "null_pwd",
			user: &userpb.UserInfo{
				Id:       12,
				Mobile:   "phone1",
				PassWord: "",
			},
			wantErr: true,
		},
		{
			name: "big role",
			user: &userpb.UserInfo{
				Id:       13,
				Mobile:   "phone1",
				PassWord: "123",
				Role:     4,
			},
			wantErr: true,
		},
		{
			name: "big gender",
			user: &userpb.UserInfo{
				Id:       13,
				Mobile:   "phone1",
				PassWord: "123",
				Gender:   2,
			},
			wantErr: true,
		},
	}
	ctx := context.Background()
	db, err := Potesting.NewClient(ctx)

	insertHelper(ctx, t, db)
	if err != nil {
		t.Fatal(err)
	}

	dao := Dao{db}
	db.AutoMigrate(&model.User{})

	if err != nil {
		panic(err)
	}

	for _, c := range cases {
		err := dao.CreateUser(ctx, c.user)
		if c.wantErr {
			if err == nil {
				t.Errorf("%s: want err but not", c.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
		}
		var u model.User
		dao.DB.Find(&u, c.user.Id)

		if c.user.Mobile != u.Mobile ||
			c.user.PassWord != u.PassWord ||
			c.user.NickName != u.NickName {
			t.Fatalf("%s:\n want:%v;\n but got: %v;\n", c.name, c.user, u)
		}
	}
}

func TestMain(m *testing.M) {
	Potesting.RunWithMongoInDocker(m)

}
