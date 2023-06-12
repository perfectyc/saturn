package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "time"
)

type Earnings struct {
    FilAmount  float64 `json:"filAmount"`
    Timestamp string  `json:"timestamp"`
}

type Data struct {
    Earnings []Earnings `json:"earnings"`
}

type Node struct {
    Count int    `json:"count"`
    State string `json:"state"`
}

type NodeResponse struct {
    Nodes []Node `json:"nodes"`
}

var filAddressLocationMap = map[string]string{
        "f1i65c367tx................": "china-01", // soaaid01a@126.com
       
    
        

}

func fetchNodeData(filAddress string) (*NodeResponse, error) {
    currentTime := time.Now()
    currentTimestampInt := int(currentTime.Unix())
    currentTimestampStr := strconv.Itoa(currentTimestampInt)
    url := fmt.Sprintf("https://uc2x7t32m6qmbscsljxoauwoae0yeipw.lambda-url.us-west-2.on.aws/?filAddress=%s&startDate=%d000&endDate=%d000&step=hour&currentTimestamp=%s000", filAddress, currentTimestampInt-3600, currentTimestampInt, currentTimestampStr)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var data NodeResponse
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }

    return &data, nil
}

func main() {
    for filAddress := range filAddressLocationMap {
        data, err := fetchNodeData(filAddress)
        if err != nil {
            fmt.Printf("Error fetching data for filAddress %s: %s\n", filAddress, err)
            continue
        }

        nodes := data.Nodes
        if len(nodes) > 0 {
            node := nodes[0]
            fmt.Printf("count: %d, state: %s, location: %s\n", node.Count, node.State, filAddressLocationMap[filAddress])
        } else {
            fmt.Println("No nodes found", filAddressLocationMap[filAddress])
        }
    }
}
