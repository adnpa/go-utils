package db

//
//import (
//	"database/sql"
//	"errors"
//	"fmt"
//	"github.com/go-sql-driver/mysql"
//	"log"
//	"sync"
//	"sync/atomic"
//)
//
//type MysqlConn struct {
//	db *sql.DB
//}
//
//var (
//	defaultConn     atomic.Pointer[MysqlConn]
//	defaultConnOnce sync.Once
//	defaultConfig   mysql.Config = mysql.Config{}
//)
//
//func DefaultConn() *MysqlConn {
//	conn := defaultConn.Load()
//	if conn == nil {
//		defaultConnOnce.Do(func() {
//			defaultConn.CompareAndSwap(
//				nil, NewConnWithConfig(defaultConfig))
//		})
//	}
//	return conn
//}
//
//func NewConnWithConfig(config mysql.Config) *MysqlConn {
//	//todo 替换config
//	db, err := sql.Open("mysql", "mydb")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	pingErr := db.Ping()
//	if pingErr != nil {
//		log.Fatal(pingErr)
//	}
//
//	conn := MysqlConn{
//		db: db,
//	}
//	return &conn
//}
//
//func (conn *MysqlConn) Close() {
//	conn.Close()
//}
//
////example
//
//type Album struct {
//	ID     int64
//	Title  string
//	Artist string
//	Price  float32
//}
//
//func albumByID(id int64) (Album, error) {
//	// An album to hold data from the returned row.
//	var alb Album
//
//	row := DefaultConn().db.QueryRow("SELECT * FROM album WHERE id = ?", id)
//	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
//		if err == sql.ErrNoRows {
//			return alb, fmt.Errorf("albumsById %d: no such album", id)
//		}
//		return alb, fmt.Errorf("albumsById %d: %v", id, err)
//	}
//	return alb, nil
//}
//
//func addAlbum(alb Album) (int64, error) {
//	result, err := DefaultConn().db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
//	if err != nil {
//		return 0, fmt.Errorf("addAlbum: %v", err)
//	}
//	id, err := result.LastInsertId()
//	if err != nil {
//		return 0, fmt.Errorf("addAlbum: %v", err)
//	}
//	return id, nil
//}
//
//func PrepareAlbumByID(id int64) (Album, error) {
//	db := DefaultConn().db
//	sqlStr := "SELECT * FROM album WHERE id = ?"
//	stmt, err := db.Prepare(sqlStr)
//	if err != nil {
//		return Album{}, err
//	}
//	defer stmt.Close()
//
//	var alb Album
//	row := stmt.QueryRow(id)
//	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return alb, fmt.Errorf("albumsById %d: no such album", id)
//		}
//		return alb, fmt.Errorf("albumsById %d: %v", id, err)
//	}
//	return alb, nil
//}
//
//func PrepareAddAlbum(alb Album) (int64, error) {
//	db := DefaultConn().db
//	sqlStr := "INSERT INTO album (title, artist, price) VALUES (?, ?, ?)"
//	stmt, err := db.Prepare(sqlStr)
//	if err != nil {
//		return 0, err
//	}
//	defer stmt.Close()
//
//	result, err := stmt.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
//	if err != nil {
//		return 0, fmt.Errorf("addAlbum: %v", err)
//	}
//	id, err := result.LastInsertId()
//	if err != nil {
//		return 0, fmt.Errorf("addAlbum: %v", err)
//	}
//	return id, nil
//}
