# saturn
### 该程序拉取了所有节点部分数据.并保存到数据库,以便于分析,文件目录为:
````
go.sum
go.mode
grafana_data_global_data.go
grafana_data_global_data_1_30.go
````
#### init\加载mysql模块
````
go mod init go.mod
go get github.com/go-sql-driver/mysql
````
#### 修改代码中得mysql权限,
````
示例:func main() {
        // 创建数据库连接
        db, err := createDBConnection("alvin:123456@tcp(192.168.1.1:3306)/global_data")
        if err != nil {
                log.Fatal(err)
        }
        defer db.Close()
````
#### 需要先在数据库中创建表
为grafana_data_global_data.go创建数据表
````

````
为grafana_data_global_data_1_30.go创建数据表
````
CREATE DATABASE global_data;

CREATE TABLE IF NOT EXISTS earning1 (
  node_id VARCHAR(255),
  earning FLOAT,
  status VARCHAR(255),
  isp VARCHAR(255),
  country VARCHAR(255),
  city VARCHAR(255),
  region VARCHAR(255),
  created TIMESTAMP,
  PRIMARY KEY (node_id)
);

CREATE TABLE IF NOT EXISTS earning30 (
  node_id VARCHAR(255),
  earning FLOAT,
  status VARCHAR(255),
  isp VARCHAR(255),
  country VARCHAR(255),
  city VARCHAR(255),
  region VARCHAR(255),
  created TIMESTAMP,
  PRIMARY KEY (node_id)
);

````
##### 之后使用命令启动程序,成功运行后,数据将保存在数据库中
````
go run *.go
````
