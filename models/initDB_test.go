package models

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestInitMysql(t *testing.T) {
	InitMysql("root@tcp(localhost)/xx1")

	sql := "select rightName from `right` where rightID=1"
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*8)
	defer cancel()
	query(ctx, sql)
}

func query(ctx context.Context, sql string) {
	c := make(chan string)
	go func() {
		rows, err := db.Query(sql)
		if err != nil {
			fmt.Println("query error: ", err)
			close(c)
		}
		defer rows.Close()
		var name string
		if rows.Next() {
			rows.Scan(&name)
		}
		c <- name
	}()
	select {
	case <-ctx.Done():
		fmt.Println("Done: ", ctx.Err())
	case v := <-c:
		fmt.Println("v:", v)
	}
}
