package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"

	"github.dena.jp/swet/go-sampleapi/internal/config"
	"github.dena.jp/swet/go-sampleapi/internal/handler"
	"github.dena.jp/swet/go-sampleapi/internal/logging"
)

// testで作成するuserのdata
const (
	email    = "test@example.com"
	password = "passw0rd!123"
)

// POST /users に対するtest
// 正常なパラメータでリクエストを行う
func Test_E2E_PostUser(t *testing.T) {
	db := sqlx.MustConnect("mysql", config.Config().DBSrc())
	defer func() {
		// DBのcleanを行う
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table users")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&handler.ReqPostUserJSON{
		FirstName: "テスト姓",
		LastName:  "テスト名",
		Email:     email,
		Password:  password,
	}); err != nil {
		t.Fatal(err)
	}

	// requestをシュミレートする
	req := httptest.NewRequest(http.MethodPost, "/", &body)
	rec := httptest.NewRecorder()
	handler.PostUser(db, logging.Logger()).ServeHTTP(rec, req)

	// responseのStatus Codeをチェックする
	if rec.Code != http.StatusCreated {
		t.Errorf("status code must be 201 but: %d", rec.Code)
		t.Fatalf("body: %s", rec.Body.String())
	}

	var result handler.ResPostUserJSON
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	// responseで返ってきたIDでuserが作られているかどうかをチェックする
	var actual string
	if err := db.Get(&actual, "select email from users where id = ?", result.ID); err != nil {
		t.Fatal(err)
	}
	if actual != email {
		t.Errorf("email must be %s but %s", email, actual)
	}
}

// POST /users に対するtest
// 重複するuserでリクエストを行う
func Test_E2E_PostUser_DuplicateEmail(t *testing.T) {
	// TODO: testを記述していく
	db := sqlx.MustConnect("mysql", config.Config().DBSrc())
	defer func() {
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table users")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()

	email := "test@example.com"

	if _, err := db.Exec("insert into users(first_name, last_name, email, password_hash) values (?, ?, ?, ?)", "dummy_first_name", "dummy_last_name", email, "dummy_password"); err != nil {
		t.Fatal(err)
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&handler.ReqPostUserJSON{
		FirstName: "test family",
		LastName: "test last",
		Email: email,
		Password: "passw0rd1234",
	}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/",  &body)
	rec := httptest.NewRecorder()
	handler.PostUser(db, logging.Logger()).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status code must be 400 but: %d", rec.Code)
	}

	var result handler.ResError
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if result.Message != string(handler.ErrEmailAlreadyExists) {
		t.Errorf("error Message must be %s but %s", handler.ErrEmailAlreadyExists, result.Message)
	}

}

func Test_E2E_PostUser_InvalidEmail(t *testing.T){
	db := sqlx.MustConnect("mysql", config.Config().DBSrc())
	defer func() {
		db.MustExec("set foreign_key_checks = 0")
		db.MustExec("truncate table users")
		db.MustExec("set foreign_key_checks = 1")
		db.Close()
	}()

	invalid_email :="test.com"
}