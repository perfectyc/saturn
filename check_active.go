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
        "f1i65c367txjuw6vcum3hzj6lnvqe3q32scxz5vfi": "lima-186", // hubei01a@126.com
        "f17sh2gebt2ntvmlzj4xyvdxouvmzw43tjakx2gji": "Sao paulo-106",  // perfectycc@gmail.com
        "f154uxz725new5mo76nvsak2f5mwfowgprmrehuxa": "indon-03",  // sichuan02a@126.com
        "f1aiv6nrmig7j3ju7xh6hjvftk5mmbln5cxirz5ja": "santiago-131",  // 1131128793@qq.com
        "f1khf2av45fyi25fa2mbay5b3kyrj7raqxx4g6i2q": "indon-05",
        "f1pctvfzyxbvj5ehvz6dhyfuwcvp3jvmfnh24tkgq": "indon-06",  // wuhan04de@126.com
        "f1hkeyu6xxul76eyweawpu7xkwjg3x7c74liqmgpy": "lima-189",   // shandongba01@126.com
        "f1vk3irnkdgcn3h2rgtkurfe66ot6cry2hrlwgtmq": "sao paulo-108",  // 578624081@qq.com
        "f1nqcoldyemhh7w4p52wx22afh4c7t4r27ua25vlq": "santiago-130",  // 799936791@qq.com
        "f1xcpclgfd5pflh6heicpgsxy25q7o5ricnjuiwaq": "Lima-106",  // baliabcper@126.com
        "f1ylo4acvzjscupr5au74npwkxl4qton5ijuhqz2y": "Indon-30", // ttbachz@126.com
        "f1lv6dmxggzhj4txmxn2mzley7j4swhucoqefd4pa": "indon-59",  // perfectycc@gmail.com
        "f1b4g3cyqnrkqgb76bife4j722stxctlqscm2l4gq": "indon-65",  // 1131128793@qq.com
        "f1d6zmchctssl346sjpxjsxqj7j3azrmoq6qmznia": "indon-67",  // filecoinabs@126.com
        "f13fyplwsgrbbe76tqea2o3o3h73msibm2dwfmtsq": "indon-18",  // zhejiangab11@126.com
        "f1fosgrg5z4e2mta5mrwi6rew56qz4yvsocjuhqzi": "indon-70",
        "f17pn4jewes2xohg7hjjosygieyhf2uhfs3mumtvi": "indon-71",
        "f1e2kvqz5khhopsy65zmmtfeq5esfmmaattbjohii": "indon-72",
        "f1lscdphje5c6youmttpz3eoflphgdr5i7b7hg35a": "indon-73",
        "f16q2qazmbggotwrjh7wep5tjjm26odpvenfmpmwq": "indon-74",
        "f1ane4tns2mkojlh3abr2jzr7dgldxrwpn3zm23xa": "Bogota-114", //  tanwan0921@126.com
        "f126d5j7kyifeq4y6nc453e33j2gi7fucn64riamy": "Bogota-115", // lrktqbgzl@126.com
        "f1f3sponexpo3h3hmfcuvhgp5jza3wmhuvwyi4rzi": "Bogota-116", // vsvjyqjljm@126.com
        "f1m3bnnudvly6bxayz2h4mvmatuul7ujk4mwr7ela": "Bogota-117", // fshcsea@126.com
        "f1wgxrxdrbf5ypl6a63wxry7lwn5npfslejgrag5y": "Bogota-118", // yktdwpuv@126.com
        "f1tmexrdyfkwe7vd57f527ly3taxsrumjzvpw3v7a": "Bogota-106",
        "f1fbwc2bivj447famqggkeiidtdjj2mds4erggacq": "Bogota-110",
        "f1uaraceel4ti22ukwhggrcuzfkguyxoqnkqt36by": "Bogota-107",
        "f1ibbsr32jrpbumfv4moobf3a5xb25oo2wjzrl32q": "Bogota-108",
        "f1y2vd4cqgmnhc4l7gafzbbfx37rxzyxgxjvodlea": "Bogota-109",
        "f1sstlv6r6eabaot4qlt4wec4x2o6kvy5plooabpa": "lima-107",   // pqygue@126.com
        "f1a5pc7kqpgmjpdpjc6ua6lds4i7rwjdnpkvfbnai": "lima-108", // ikdzfgodjsmy@126.com
        "f1e5m6xmpaovgnd3dqcmlawos4pp4bffynyn5gpni": "lima-109", // femmufb@126.com
        "f1omt6mjc7z6ukpkpos2ppcuh65tdmnvcfhvwsgta": "lima-110", // cwbpffeg@126.com
        "f1lr4v3om6nsge4xbusajcjye5yajmxs65ku7mkti": "lima-187",
        "f1izsvptqibbvkmt6ct7d7rjze6xm2xjl22r77gzi": "lima-188", // oloabdyzuwpk@126.com
        "f1gxuuigdjlneu4czitw5ognfpop2a5exdd5fozpy": "lima-190",
        "f1zug3m3pokmrhrgalj2r7g5enhz252hdszbfzboi": "Sao paulo-109", // suxmrokzsqpro@126.com
        "f1x3kcbv3bvyulsltsxjanruap6j4ri5jdf64i2fy": "Sao paulo-107", // aynkensk@126.com
        "f1ck3vwhecdx7gnxdkawgvnv7zowfhqxwsujkl3qa": "Sao paulo-110", // zwmfinpvt@126.com
        "f14py5udegnazgkbspouj3urpjcejhy6bh2soie4a": "Sao paulo-138",
        "f1x4xqn76m5o4xp3s77fjcw5fkp2s3ozqpxdlh4ui": "Sao paulo-139",
        "f142ap442yjcaizaixk2kedic2mqdtf6wyonysv7q": "Sao paulo-140",
        "f15zlkcmhqouxqtfgjscacw2xrhhw4pzlqllfvvry": "Sao paulo-141",
        "f12gbyazhqw4kahgoyq4ov25jzooqbfc45kb3m4xy": "santiago-58", // hangzhou11a@126.com
        "f1isqweuhjuk52bsbkzfazmzmx6ihu26rnwinyzai": "santiago-59",  // hlqryhgq@126.com
        "f1ye3t2agqglpuogv6m6sxihx3zeh5uva3w6bdm3a": "santiago-60", // ryolxdsn@126.com
        "f15ieg2vwwf3d2k6htc34x2yutiwwz5ruks3qjawa": "santiago-61", // wlozrqsk@126.com
        "f1ami23ecxlufob366x643df6ljtslnp6uifskb2i": "santiago-62", // indonb2c@126.com
        "f1byjtiuwi7kjt3wsda3mspc24mvatdp55adbwuqy": "santiago-132",
        "f1bwivkebqxnjhhxcsu3sdbljvo2nctnm2d7ueawa": "santiago-133",
        "f1ne2bbossqapb36ct23dqauqx73636usmlmwenby": "santiago-134",
        "f14whyjseme3lea23llrkb7atpbsl5ouhiqyjnaza": "HL-02",
        "f1nq3tpy3wconl4zi66wbzbtgl5rdhbz5ehl32qei": "HL-03",
        "f13mva246baiicxdkepgs4jlg2ir74qqyxujycuca": "HL-04",
        "f1aan76kpfdq4vn3vx6ccu53clqyrcgjl4rn4uwya": "HL-05",
        "f1yo4zwwxifud474hqjtqs2tmjajfq5x7ofh4nida": "HL-06",
        "f1xfpwurvmvt4lbhcm6mlprm7tz5jkk2z22dcyuzi": "HL-10",
        "f12obaoelrh3szgsnd3cndhhqu3cfx25gu4lioa2a": "HL-11",
        "f1fv6mlcermjveuknqlbticukgyiiu6prwuwclcvy": "HL-12",
        "f1hjlinnku6jazowda2y2frcipyf7gtnvqcj6paia": "HL-13",
        "f1pbgjumlm7qud5v654fqxzlz5k73qvzjg2p32gki": "HL-14",
   
    
        

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
