package main

import (
        "database/sql"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "time"

        _ "github.com/go-sql-driver/mysql"
)

type Node struct {
        Count int    `json:"count"`
        State string `json:"state"`
}

func main() {
        // 创建数据库连接
        db, err := sql.Open("mysql", "user:passwd@tcp(xx.xx.xx.xx:3306)/globaldate")
        if err != nil {
                log.Fatal(err)
        }
        defer db.Close()

        // 获取当前时间戳
        currentTime := time.Now().Unix() * 1000

        // 构建 URL
        url := fmt.Sprintf("https://uc2x7t32m6qmbscsljxoauwoae0yeipw.lambda-url.us-west-2.on.aws/?filAddress=all&startDate=%d&endDate=%d&step=hour", currentTime, currentTime)

        // 发起 HTTP 请求获取 JSON 数据
        resp, err := http.Get(url)
        if err != nil {
                log.Fatal(err)
        }
        defer resp.Body.Close()

        // 读取响应体
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatal(err)
        }

        // 解析 JSON 数据
        var data struct {
                Nodes []Node `json:"nodes"`
        }
        err = json.Unmarshal(body, &data)
        if err != nil {
                log.Fatal(err)
        }

        // 创建新的表
        tableName := "globaldate"
        createTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INT AUTO_INCREMENT PRIMARY KEY, state VARCHAR(50) NOT NULL, count INT NOT NULL)", tableName)
        _, err = db.Exec(createTableQuery)
        if err != nil {
                log.Fatal(err)
        }

        // 打印和保存数据到 MySQL 数据库
        for _, node := range data.Nodes {
                fmt.Printf("State: %s, Count: %d\n", node.State, node.Count)

                // 将数据插入到新的表中
                insertQuery := fmt.Sprintf("INSERT INTO %s (state, count) VALUES (?, ?)", tableName)
                _, err := db.Exec(insertQuery, node.State, node.Count)
                if err != nil {
                        log.Fatal(err)
                }
        }
}
