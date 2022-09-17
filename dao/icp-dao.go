package dao

import (
	"database/sql"
	"icp-search/entity"
	init_ "icp-search/init"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)
var (
	db *sql.DB
)

func Init()  {
	var err error
	db, err = sql.Open("mysql", init_.Cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func UnInit()  {
	if db != nil {
		db.Close()
	}
}

func Search(domain string) (icp *entity.Icp, err error) {
	stmt, err := db.Prepare("select * from icps where domain = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	icp = &entity.Icp{
		Domain:   "",
		Unit:     "",
		Type:     "",
		IcpCode:  "",
		Name:     "",
		PassTime: "",
		CacheTime: "",
	}

	err = stmt.QueryRow(domain).Scan(&icp.Domain, &icp.Unit, &icp.Type, &icp.IcpCode, &icp.Name, &icp.PassTime, &icp.CacheTime)
	if err != nil {
		return nil, err
	}
	return icp, nil
}

func Insert(icp *entity.Icp) error  {
	stmt, err := db.Prepare(`insert into icps(domain, unit, type, icpCode, name, passTime, cacheTime)
			    value(?,?,?,?,?,?,?) on duplicate key update unit=?, type=?, icpCode=?,name=?, passTime=?, cacheTime=?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(icp.Domain, icp.Unit, icp.Type, icp.IcpCode, icp.Name, icp.PassTime, icp.CacheTime,
		icp.Unit, icp.Type, icp.IcpCode, icp.Name, icp.PassTime, icp.CacheTime)
	return err
}