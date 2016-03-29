package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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
	if val := os.Getenv("DB_USER"); val != "" {
		databaseUser = val
	}
	if val := os.Getenv("DB_PASS"); val != "" {
		databasePassword = val
	}
	if val := os.Getenv("DB_HOST"); val != "" {
		databaseHost = val
	}
	if val := os.Getenv("DB_PORT"); val != "" {
		databasePort = val
	}
	if val := os.Getenv("DB_NAME"); val != "" {
		databaseName = val
	}
	if val := os.Getenv("NSQ_TOPIC"); val != "" {
		nsqTopic = val
	}
	if val := os.Getenv("NSQ_CHAN"); val != "" {
		nsqChannel = val
	}
	flag.StringVar(&databaseUser, "dbuser", databaseUser, "MySQL Username (env: DB_USER)")
	flag.StringVar(&databasePassword, "dbpass", databasePassword, "MySQL Password (env: DB_PASS)")
	flag.StringVar(&databaseHost, "dbhost", databaseHost, "MySQL TCP Hostname (env: DB_HOST)")
	flag.StringVar(&databasePort, "dbport", databasePort, "MySQL TCP Port (env: DB_PORT)")
	flag.StringVar(&databaseName, "dbname", databaseName, "MySQL Database Name (env: DB_NAME)")
	flag.StringVar(&nsqTopic, "nsqtopic", nsqTopic, "NSQD Topic (env: DB_TOPIC)")
	flag.StringVar(&nsqChannel, "nsqchan", nsqChannel, "NSQD Channel (env: NSQ_CHAN)")
	flag.StringVar(&nsqAddress, "nsqaddr", nsqAddress, "NSQD Address (env: NSQ_ADDR)")
	flag.Parse()
}

func mustPrepare(db *sql.DB, name string, query string) *sql.Stmt {
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
	db := mustConnect()
	defer db.Close()

	stats := mustPrepare(
		db,
		"stats insert stmt",
		"INSERT IGNORE INTO `stats` (platform,product,stat,sub1,sub2,sub3,info) VALUES(?,?,?,?,?,?,?)")
	setDaily := mustPrepare(
		db,
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
		db,
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
		db,
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
		db,
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
				_, err = tx.Stmt(incrHourly).Exec(stat.When.Format(dateTimeFormat),
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
				_, err = tx.Stmt(setHourly).Exec(stat.When.Format(dateTimeFormat),
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
	var wait chan struct{}
	<-wait
}
