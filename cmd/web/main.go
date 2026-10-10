package main

import (
	"aphelios_website/internal/models"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	pages          *models.PageModel
	tils           *models.TilModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder //表单解析
	users          *models.UserModel
	sessionManager *scs.SessionManager
	blogs          *models.BlogModel
	logger         *slog.Logger
}

// func OpenDB(dsn string) (*sql.DB, error) {
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		db.Close()
// 		return nil, err
// 	}
// 	return db, nil
// }

func loadPostGres() (models.PostGresConfig, error) {
	var pgcfg models.PostGresConfig
	err := godotenv.Load()
	if err != nil {
		return pgcfg, err
	}
	pgcfg = models.PostGresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSL_MODE"),
	}
	if pgcfg.Host == "" && pgcfg.Port == "" {
		return pgcfg, fmt.Errorf("No PSQL Config provided.")
	}

	return pgcfg, nil
}

func main() {

	//数据库连接
	// dsn := "root:19980131@/aphelios_website?parseTime=true"
	// db, err := OpenDB(dsn)
	// if err != nil {
	// 	fmt.Printf("DB conn err:%v", err)
	// 	return
	// }

	// defer db.Close()
	pgcfg, err := loadPostGres()
	if err != nil {
		panic(err)
	}

	db, err := models.Open(pgcfg)
	if err != nil {
		fmt.Printf("DB conn err:%v", err)
		return
	}
	fmt.Println("db connection success")
	defer db.Close()

	templateHtml, err := newTemplates()
	if err != nil {
		fmt.Printf("templatecache err:%v", err)
		return
	}

	//表单解析器
	formDecoder := form.NewDecoder()

	//配置它以使用我们的 MySQL 数据库作为会话存储，并设置一个有效期为 12 小时（因此会话将在 12 小时后自动过期,在首次创建之后）。
	//用户的登录状态 / 会话数据会失效，需要重新登录；网站本身依旧可以正常访问。
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	//日志记录器
	// 初始化一个新的结构化日志记录器，将日志条目写入标准输出
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		pages:          &models.PageModel{DB: db},
		templateCache:  templateHtml,
		formDecoder:    formDecoder,
		users:          &models.UserModel{DB: db},
		sessionManager: sessionManager,
		tils:           &models.TilModel{DB: db},
		blogs:          &models.BlogModel{DB: db},
		logger:         logger,
	}

	err = app.server()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
