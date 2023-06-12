package main

import (
        "database/sql"
        "fmt"
        "io"
        "log"
        "net/http"
        "time"

        "github.com/goccy/go-json"
        _ "github.com/go-sql-driver/mysql"
        "github.com/lnzx/strnx/tools"
)

const (
        NodeStatusUrl   = "https://orchestrator.strn.pl/stats?sortColumn=id"
        NodesEarningUrl = "https://uc2x7t32m6qmbscsljxoauwoae0yeipw.lambda-url.us-west-2.on.aws/?filAddress=all&startDate=%d&endDate=%d&step=day"
        UPSERT_EARN_1   = "INSERT INTO earning1(node_id, earning, status, isp, country, city, region, created) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE earning = VALUES(earning)"
        UPSERT_EARN_30  = "INSERT INTO earning30(node_id, earning, status, isp, country, city, region, created) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE earning = VALUES(earning)"
)

var NodeStatusMap = map[string]Status{}
var DB *sql.DB

type Earning struct {
        FilAmount float32 `json:"filAmount"`
        Timestamp string  `json:"timestamp"`
}

type NodeMetrics struct {
        PerNodeMetrics []PerNodeMetrics `json:"perNodeMetrics"`
}

type PerNodeMetrics struct {
        NodeId       string    `json:"nodeId"`
        FilAmount    float32   `json:"filAmount"`
        PayoutStatus string    `json:"payoutStatus"`
        Max          string    `json:"max"`
        Isp          string    `json:"isp"`
        Country      string    `json:"country"`
        City         string    `json:"city"`
        Region       string    `json:"region"`
        Created      time.Time `json:"created"`
}

type NodeStatsResult struct {
        Nodes []Status `json:"nodes"`
}

type Status struct {
        Id     string `json:"id"`
        State  string `json:"state"`
        Geoloc struct {
                Country string `json:"country"`
                City    string `json:"city"`
                Region  string `json:"region"`
        } `json:"geoloc"`
        Speedtest struct {
                Isp string `json:"isp"`
        } `json:"speedtest"`
        Created time.Time `json:"createdAt"`
}

func main() {
        // 创建数据库连接
        db, err := createDBConnection("user:passwd@tcp(xx.xx.xx.xx:3306)/global_data")
        if err != nil {
                log.Fatal(err)
        }
        defer db.Close()

        // 设置数据库连接池的最大连接数
        db.SetMaxOpenConns(10)

        // 检查数据库连接是否正常
        err = db.Ping()
        if err != nil {
                log.Fatal(err)
        }

        // 设置全局DB对象
        DB = db

        // 调用FetchNodesEarningJob()函数获取节点收益数据
        FetchNodesEarningJob()
}

func createDBConnection(connectionString string) (*sql.DB, error) {
        // 创建数据库连接
        db, err := sql.Open("mysql", connectionString)
        if err != nil {
                return nil, err
        }

        return db, nil
}

func FetchNodesEarningJob() {
        now := time.Now().UTC()
        start1, end1 := tools.GetBeforeDayN(now, 1)
        start30, end30 := tools.GetBeforeDayN(now, 30)

        metrics1, err := fetchNodesEarning(start1, end1)
        if err != nil {
                log.Println(err)
                return
        }
        if len(metrics1) == 0 {
                log.Println("cron FetchNodesEarningJob metrics1 0 skip")
                return
        }

        metrics30, err := fetchNodesEarning(start30, end30)
        if err != nil {
                log.Println(err)
                return
        }
        if len(metrics30) == 0 {
                log.Println("cron FetchNodesEarningJob metrics30 0 skip")
                return
        }

        statusMap, err := fetchNodesStatus()
        if err == nil {
                NodeStatusMap = statusMap
        }

        // 开始事务
        batch, err := DB.Begin()
        if err != nil {
                log.Println(err)
                return
        }
        defer batch.Rollback()

        for _, node := range metrics1 {
                if node.FilAmount == 0 {
                        continue
                }
                if status, ok := NodeStatusMap[node.NodeId]; ok {
                        geo := status.Geoloc
                        _, err = batch.Exec(UPSERT_EARN_1, node.NodeId, node.FilAmount, node.PayoutStatus,
                                status.Speedtest.Isp, geo.Country, geo.City, geo.Region, status.Created)
                        if err != nil {
                                log.Println(err)
                        }
                }
        }

        for _, node := range metrics30 {
                if node.FilAmount == 0 {
                        continue
                }
                if status, ok := NodeStatusMap[node.NodeId]; ok {
                        geo := status.Geoloc
                        _, err = batch.Exec(UPSERT_EARN_30, node.NodeId, node.FilAmount, node.PayoutStatus,
                                status.Speedtest.Isp, geo.Country, geo.City, geo.Region, status.Created)
                        if err != nil {
                                log.Println(err)
                        }
                }
        }

        // 提交事务
        err = batch.Commit()
        if err != nil {
                log.Println(err)
                return
        }

        log.Printf("cron FetchNodesEarningJob started %s\n", time.Now().UTC().Sub(now).String())
}

func fetchNodesEarning(start, end time.Time) ([]PerNodeMetrics, error) {
        url := fmt.Sprintf(NodesEarningUrl, start.UnixMilli(), end.UnixMilli())
        rsp, err := tools.Get(url)
        if err != nil {
                if rsp != nil {
                        e := rsp.Body.Close()
                        if e != nil {
                                log.Println(e)
                        }
                }
                return nil, err
        }

        bytes, err := io.ReadAll(rsp.Body)
        if err != nil {
                return nil, err
        }
        var metrics NodeMetrics
        err = json.Unmarshal(bytes, &metrics)
        if err != nil {
                return nil, err
        }
        return metrics.PerNodeMetrics, nil
}

func fetchNodesStatus() (map[string]Status, error) {
        req, err := http.NewRequest("GET", NodeStatusUrl, nil)
        if err != nil {
                return nil, err
        }
        req.Header.Add("Accept", "application/json")

        rsp, err := tools.Do(req)
        if err != nil {
                if rsp != nil {
                        if e := rsp.Body.Close(); e != nil {
                                log.Println(e)
                        }
                }
                return nil, err
        }

        bytes, err := io.ReadAll(rsp.Body)
        if err != nil {
                return nil, err
        }
        var r NodeStatsResult
        err = json.Unmarshal(bytes, &r)
        if err != nil {
                return nil, err
        }

        return ConvertNodesToMap(r.Nodes), nil
}

func ConvertNodesToMap(nodes []Status) map[string]Status {
        statusMap := make(map[string]Status, len(nodes))
        for _, node := range nodes {
                statusMap[node.Id] = node
        }
        return statusMap
}
