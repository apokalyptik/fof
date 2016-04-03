package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/apokalyptik/cfg"
	"github.com/nsqio/go-nsq"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dateTimeFormat = "2006-01-02 15:04:05"
)

var (
	databaseUser     = "my_db_user"
	databasePassword = "my_db_pass"
	databaseHost     = "the_db_host"
	databasePort     = "3306"
	databaseName     = "my_db_name"
	nsqTopic         = "fof-stats"
	nsqChannel       = "stats-to-mysql"
	nsqAddress       = "127.0.0.1:4150"
	db               *sql.DB
)

type statMessage struct {
	Platform string    `json:"platform"`
	Member   string    `json:"member"`
	Product  string    `json:"product"`
	Stat     string    `json:"stat"`
	Sub1     string    `json:"sub1"`
	Sub2     string    `json:"sub2"`
	Sub3     string    `json:"sub3"`
	Info     string    `json:"info"`
	When     time.Time `json:"When"`
	Value    int       `json:"value"`
	Method   string    `json:"method"`
}

func init() {
	dbc := cfg.New("db")
	dbc.StringVar(&databaseUser, "user", databaseUser, "MySQL Username (env: DB_USER)")
	dbc.StringVar(&databasePassword, "pass", databasePassword, "MySQL Password (env: DB_PASS)")
	dbc.StringVar(&databaseHost, "host", databaseHost, "MySQL TCP Hostname (env: DB_HOST)")
	dbc.StringVar(&databasePort, "port", databasePort, "MySQL TCP Port (env: DB_PORT)")
	dbc.StringVar(&databaseName, "name", databaseName, "MySQL Database Name (env: DB_NAME)")

	nsqc := cfg.New("nsq")
	nsqc.StringVar(&nsqTopic, "topic", nsqTopic, "NSQD Topic (env: NSQ_TOPIC)")
	nsqc.StringVar(&nsqChannel, "chan", nsqChannel, "NSQD Channel (env: NSQ_CHAN)")
	nsqc.StringVar(&nsqAddress, "addr", nsqAddress, "NSQD Address (env: NSQ_ADDR)")

	cfg.Parse()
}

func mustPrepare(name string, query string) *sql.Stmt {
	s, err := db.Prepare(query)
	if err != nil {
		log.Fatal("Error creating %s: %s", name, err.Error())
	}
	return s
}

func mustConnect() *sql.DB {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", databaseUser, databasePassword, databaseHost, databasePort, databaseName))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db = mustConnect()
	defer db.Close()
	stats := mustPrepare(
		"stats insert stmt",
		"INSERT IGNORE INTO `stats` (platform,product,stat,sub1,sub2,sub3,info) VALUES(?,?,?,?,?,?,?)")
	latest := mustPrepare(
		"stats_latest update statement",
		"INSERT INTO `stats_latest` (`member`,`stat_id`,`daily`,`hourly`)"+
			"  SELECT ?,ID,DATE_FORMAT(?,'%Y-%m-%d %H:00:00'),DATE_FORMAT(?,'%Y-%m-%d %H:00:00')"+
			"    FROM `stats`"+
			"    WHERE platform=?"+
			"      AND product=?"+
			"      AND stat=?"+
			"      AND sub1=?"+
			"      AND sub2=?"+
			"      AND sub3=?"+
			"	  ON DUPLICATE KEY UPDATE `daily`=VALUES(`daily`),`hourly`=VALUES(`hourly`)")
	setDaily := mustPrepare(
		"daily insert stmt",
		"INSERT INTO `stats_daily` (`when`,`stat_id`,`member`,`value`)"+
			"  SELECT ?,ID,?,?"+
			"    FROM `stats`"+
			"    WHERE platform=?"+
			"      AND product=?"+
			"      AND stat=?"+
			"      AND sub1=?"+
			"      AND sub2=?"+
			"      AND sub3=?"+
			"  ON DUPLICATE KEY UPDATE `value`=?")
	incrDaily := mustPrepare(
		"daily insert stmt",
		"INSERT INTO `stats_daily` (`when`,`stat_id`,`member`,`value`)"+
			"  SELECT ?,ID,?,?"+
			"    FROM `stats`"+
			"    WHERE platform=?"+
			"      AND product=?"+
			"      AND stat=?"+
			"      AND sub1=?"+
			"      AND sub2=?"+
			"      AND sub3=?"+
			"  ON DUPLICATE KEY UPDATE `value`=`value`+?")
	setHourly := mustPrepare(
		"hourly insert stmt",
		"INSERT INTO `stats_hourly` (`when`,`stat_id`,`member`,`value`)"+
			"  SELECT DATE_FORMAT(?,'%Y-%m-%d %H:00:00'),ID,?,?"+
			"    FROM `stats`"+
			"    WHERE platform=?"+
			"      AND product=?"+
			"      AND stat=?"+
			"      AND sub1=?"+
			"      AND sub2=?"+
			"      AND sub3=?"+
			"  ON DUPLICATE KEY UPDATE `value`=?")
	incrHourly := mustPrepare(
		"hourly insert stmt",
		"INSERT INTO `stats_hourly` (`when`,`stat_id`,`member`,`value`)"+
			"  SELECT DATE_FORMAT(?,'%Y-%m-%d %H:00:00'),ID,?,?"+
			"    FROM `stats`"+
			"    WHERE platform=?"+
			"      AND product=?"+
			"      AND stat=?"+
			"      AND sub1=?"+
			"      AND sub2=?"+
			"      AND sub3=?"+
			"  ON DUPLICATE KEY UPDATE `value`=`value`+?")

	consumer, err := nsq.NewConsumer(nsqTopic, nsqChannel, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		var stat statMessage
		if err := json.Unmarshal(m.Body, &stat); err != nil {
			log.Printf("Error unmarshalling stat message: %s -- %s", err.Error(), string(m.Body))
		}
		tx, err := db.Begin()
		if err != nil {
			log.Fatal("Error beginning transaction: %s", err.Error())
		}
		_, err = tx.Stmt(stats).Exec(stat.Platform, stat.Product, stat.Stat, stat.Sub1, stat.Sub2, stat.Sub3, stat.Info)
		switch stat.Method {
		case "", "inc":
			if err == nil {
				_, err = tx.Stmt(incrDaily).Exec(
					stat.When.Format(dateTimeFormat),
					stat.Member,
					stat.Value,
					stat.Platform,
					stat.Product,
					stat.Stat,
					stat.Sub1,
					stat.Sub2,
					stat.Sub3,
					stat.Value)
			}
			if err == nil {
				_, err = tx.Stmt(incrHourly).Exec(
					stat.When.Format(dateTimeFormat),
					stat.Member,
					stat.Value,
					stat.Platform,
					stat.Product,
					stat.Stat,
					stat.Sub1,
					stat.Sub2,
					stat.Sub3,
					stat.Value)
			}
			break
		case "set":
			if err == nil {
				_, err = tx.Stmt(setDaily).Exec(
					stat.When.Format(dateTimeFormat),
					stat.Member,
					stat.Value,
					stat.Platform,
					stat.Product,
					stat.Stat,
					stat.Sub1,
					stat.Sub2,
					stat.Sub3,
					stat.Value)
			}
			if err == nil {
				_, err = tx.Stmt(setHourly).Exec(
					stat.When.Format(dateTimeFormat),
					stat.Member,
					stat.Value,
					stat.Platform,
					stat.Product,
					stat.Stat,
					stat.Sub1,
					stat.Sub2,
					stat.Sub3,
					stat.Value)
			}
			break
		}
		if err == nil {
			_, err = tx.Stmt(latest).Exec(
				stat.Member,
				stat.When.Format(dateTimeFormat),
				stat.When.Format(dateTimeFormat),
				stat.Platform,
				stat.Product,
				stat.Stat,
				stat.Sub1,
				stat.Sub2,
				stat.Sub3)
		}
		if err != nil {
			log.Printf("Error Inserting or updating stats: %s", err.Error())
			tx.Rollback()
			return nil
		}
		tx.Commit()
		log.Println(stat)
		return nil
	}))
	if err := consumer.ConnectToNSQD(nsqAddress); err != nil {
		log.Fatal(err)
	}
	mindHTTP()
}
