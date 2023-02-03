package interfaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tmkshy1908/Portfolio/domain"
	"github.com/tmkshy1908/Portfolio/pkg/infrastructure/db"
	"github.com/tmkshy1908/Portfolio/pkg/infrastructure/line"
)

type CommonRepository struct {
	DB  db.SqlHandler
	Bot line.Client
}

const (
	// SELECT_SCHEDULE string = "select * from schedule;"
	SELECT_CONTENTS string = "select * from contents;"
	INSERT_CONTENTS string = "insert into contents (contents_day, location, event_title, act, other_info) values(TO_DATE('%s', 'YY-MM-DD'),'%s','%s','%s','%s')"
	UPDATE_CONTENTS string = "update contents set (contents_day, location, event_title, act, other_info) values(TO_DATE('%s', 'YY-MM-DD'),'%s','%s','%s','%s') "
	DELETE_CONTENTS string = "delete from contents where contents_day = TO_DATE('%s', 'YY-MM-DD')"
	DAY_CHECK       string = "select * from test where day = TO_DATE($1, 'YY-MM-DD HH24:MI:SS')"
	USER_CHECK      string = "select count(user_id) from users where user_id = '%s'"
	INSERT_USERS    string = "insert into users (user_id, condition) values('%s',%b)"
	TEST_CHECK      string = "select * from %s where %s = '%s'"
)

func (r *CommonRepository) Find(ctx context.Context) (contents []*domain.Contents, err error) {
	rows, err := r.DB.Query(ctx, SELECT_CONTENTS)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	contents = make([]*domain.Contents, 0)

	for rows.Next() {
		contentsTable := domain.Contents{}
		if err = rows.Scan(
			&contentsTable.ID,
			&contentsTable.Contents_Day,
			&contentsTable.Location,
			&contentsTable.EventTitle,
			&contentsTable.Act,
			&contentsTable.OtherInfo,
		); err != nil {
			fmt.Println(err)
			return
		}
		contents = append(contents, &contentsTable)
	}
	return
}

func (r *CommonRepository) Add(ctx context.Context, contents *domain.Contents) (err error) {
	fmt.Println(&contents)
	fmt.Println(contents.Contents_Day, contents.EventTitle, contents.Location, contents.Act, contents.OtherInfo)

	// contentsTable := make([]*domain.Contents,0)saa
	value := fmt.Sprintf("insert into schedule (day) values (TO_DATE('%s', 'YY-MM-DD HH24:MI:SS'))", contents.Contents_Day)
	_, err = r.DB.Exec(ctx, value)
	if err != nil {
		fmt.Println(err, "Execエラー")
		return err
	}
	values := fmt.Sprintf(INSERT_CONTENTS, contents.Contents_Day, contents.Location, contents.EventTitle, contents.Act, contents.OtherInfo)
	_, err = r.DB.Exec(ctx, values)
	if err != nil {
		fmt.Println(err, "Execエラー")
		return err
	}
	return
}

func (r *CommonRepository) Update(ctx context.Context, contents *domain.Contents) (err error) {
	values := fmt.Sprintf(UPDATE_CONTENTS, contents.Contents_Day, contents.EventTitle, contents.Location, contents.Act, contents.OtherInfo)
	_, err = r.DB.Exec(ctx, values)
	if err != nil {
		fmt.Println(err, "Updateエラー")
		return err
	}
	return
}

func (r *CommonRepository) Delete(ctx context.Context, contents *domain.Contents) (err error) {
	values := fmt.Sprintf(DELETE_CONTENTS, contents.Contents_Day)
	_, err = r.DB.Exec(ctx, values)
	if err != nil {
		fmt.Println(err, "Deleteエラー")
		return err
	}
	return
}

func (r *CommonRepository) DivideEvent(ctx context.Context, req *http.Request) (msg string, userId string) {
	msg, userId = r.Bot.CathEvents(ctx, req)
	fmt.Println(userId)

	return
}

func (r *CommonRepository) CallReply(msg string, userId string) {
	r.Bot.MsgReply(msg, userId)
}

func (r *CommonRepository) WaitMsg(ctx context.Context) (contents *domain.Contents, err error) {
	day, location, title, act, info := r.Bot.WaitEvents(ctx)
	fmt.Println(day)
	contents = &domain.Contents{Contents_Day: day, Location: location, EventTitle: title, Act: act, OtherInfo: info}
	return
}

func (r *CommonRepository) UserCheck(ctx context.Context, userId string) {
	values := fmt.Sprintf(USER_CHECK, userId)
	var i int
	err := r.DB.QueryRow(ctx, values).Scan(&i)
	if err != nil {
		fmt.Println(err, "クエリ")
	}
	if i == 0 {
		values = fmt.Sprintf(INSERT_USERS, userId, 0)
		fmt.Println(values)
		_, err = r.DB.Exec(ctx, values)
		if err != nil {
			fmt.Println(err, "eku")
		}
	} else {
		fmt.Println("登録されています")
	}
}

// func (r *CommonRepository) dayCheck(ctx context.Context, day string) {
// 	values := fmt.Sprintf(DAY_CHECK, day)
// 	_, err := r.DB.Exec(ctx,values)
// 	if err != nil {
// 		fmt.Println("Create Execエラー:", err)
// 		t = false
// 		return
// 	}
// 	t = true
// }

func (r *CommonRepository) TestTest(ctx context.Context, req *http.Request) {

	// for i := 0; i < 5; i++ {
	// 	r.Bot.TestFunc(ctx, req)
	// }
	// a := "Contents"
	// b := "Location"
	// c := "新宿"
	// values := fmt.Sprintf(TEST_CHECK, a, b, c)
	// // fmt.Println((values))
	// // aa := domain.Contents{}
	// rows, err := r.DB.Query(ctx, values)
	// if err != nil {
	// 	fmt.Println(err, "TestCheck ERROR")
	// 	return
	// }
	// for rows.Next(){
	// 	contents := domain.Contents{}
	// 	err = rows.Scan(&contents.ID, &contents.Contents_Day, &contents.Location,&contents.EventTitle, &contents.Act, &contents.OtherInfo)
	// 	if err != nil{
	// 		fmt.Println(err,"ScanError")
	// 	}

	// }
	// if aa == nil {
	// 	fmt.Println("値なし")
	// }
	// fmt.Println(aa)
	// fmt.Println("ループ処理終わり")
}
