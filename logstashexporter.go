
package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
	"net/http"
	"strings"
	"strconv"
	"regexp"
	"io/ioutil"
	"log"
)

//type ExtraStuff map[string] interface{}
type jsonData struct {
	LogMessages []Messages `json:"log_message"`
	ProcessName string `json:"process_name"`
	HostName    string `json:"host_name"`
	ThreadID   string `json:"thread_id"`
	RequestFull   string `json:"request"`
	ProcessID   string `json:"process_id"`
	Level       string `json:"level"`
	LogTime     string `json:"time"`
	//FileName    string
	//LineNum     int
}


type Messages struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type conf struct {
	Regex string `yaml:"regex"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("logex.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func parseGhPost(rw http.ResponseWriter, request *http.Request) {
    	decoder := json.NewDecoder(request.Body)
	var t jsonData
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println("Not recognized as monitoring message!",err.Error())
		return
	}
	  /*newt, err := json.Marshal(&t)
            if err != nil {
                    fmt.Println("Error Converting to Json!")
                }
            fmt.Println(string(newt))*/
	    //Creating Endpoint labels
	    	var RequestEndpoint string
		Regexp := regexp.MustCompile(con.Regex)
		if(string(t.RequestFull)==""){
			//Request field empty
		}else
		{
			RequestEndpoint = Regexp.FindString(string(t.RequestFull))
		}
		RequestEndpoint = strings.Replace(RequestEndpoint, " ", "_", -1)

            for _, message := range t.LogMessages {
                switch string(message.Key) {
                case "latency":
                    setvalue, err := strconv.ParseFloat(message.Value, 64)
                    if err != nil {
                        panic(err)
                    }
                    Latency.WithLabelValues(string(t.ProcessName), string(t.HostName), string(t.ThreadID), RequestEndpoint).Set(setvalue)
		case "request_time":
                    setvalue, err := strconv.ParseFloat(message.Value, 64)
                    if err != nil {
                        panic(err)
                    }
		    RequestTime.WithLabelValues(string(t.ProcessName), string(t.HostName), string(t.ThreadID), RequestEndpoint).Set(setvalue)
		case "upstream_response_time":
                    setvalue, err := strconv.ParseFloat(message.Value, 64)
                    if err != nil {
                        panic(err)
                    }
	            UpstreamResponseTime.WithLabelValues(string(t.ProcessName), string(t.HostName), string(t.ThreadID), RequestEndpoint).Set(setvalue)
		case "status":
                    setvalue, err := strconv.ParseFloat(message.Value, 64)
                    if err != nil {
                        panic(err)
                    }
		    if (setvalue >= 500 && setvalue < 600){
			    Http5XXcode.WithLabelValues(string(t.ProcessName), string(t.HostName), string(t.ThreadID), RequestEndpoint).Inc()} else if
		    (setvalue >= 400 && setvalue < 500){
			    Http4XXcode.WithLabelValues(string(t.ProcessName), string(t.HostName), string(t.ThreadID), RequestEndpoint).Inc()} else if
	 	    (setvalue >= 200 && setvalue < 300){
			    Http2XXcode.WithLabelValues(string(t.ProcessName), string(t.HostName), string(t.ThreadID), RequestEndpoint).Inc()}
                case "http_5xx":
                    Http5XXcode.WithLabelValues(string(t.ProcessName), string(t.HostName), string(t.ThreadID), RequestEndpoint).Inc()
                case "http_4xx":
                    Http4XXcode.WithLabelValues(string(t.ProcessName), string(t.HostName), string(t.ThreadID), RequestEndpoint).Inc()
                case "http_2xx":
                    Http2XXcode.WithLabelValues(string(t.ProcessName), string(t.HostName), string(t.ThreadID), RequestEndpoint).Inc()
                case "request_count":
                    RequestCount.WithLabelValues(string(t.ProcessName), string(t.HostName), string(t.ThreadID), RequestEndpoint).Inc()
                }
            }

}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("")
}

var (
	Latency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "logstash_latency",
		Help: "Latency of the service.",
	},
		[]string{"processname", "hostname", "processid"},
	)
	RequestTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "logstash_request_time",
		Help: "Nginx request time.",
	},
		[]string{"processname", "hostname", "processid", "request_endpoint"},
	)
	UpstreamResponseTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "logstash_upstream_response_time",
		Help: "Nginx Upstream response time.",
	},
		[]string{"processname", "hostname", "processid", "request_endpoint"},
	)
	BodyBytesSent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "logstash_body_bytes_sent",
		Help: "Nginx Size of body.",
	},
		[]string{"processname", "hostname", "processid", "request_endpoint"},
	)
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "logstash_request_count",
			Help: "Number of Requests.",
		},
		[]string{"processname", "hostname", "processid", "request_endpoint"},
	)
	Http5XXcode = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "logstash_http_5xx_count",
			Help: "Number of 5xx status codes.",
		},
		[]string{"processname", "hostname", "processid", "request_endpoint"},
	)
	Http2XXcode = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "logstash_http_2xx_count",
			Help: "Number of 2xx status codes.",
		},
		[]string{"processname", "hostname", "processid", "request_endpoint"},
	)
	Http4XXcode = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "logstash_http_4xx_count",
			Help: "Number of 4xx status codes.",
		},
		[]string{"processname", "hostname", "processid", "request_endpoint"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(Latency)
	prometheus.MustRegister(RequestTime)
	prometheus.MustRegister(UpstreamResponseTime)
	prometheus.MustRegister(BodyBytesSent)
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(Http5XXcode)
	prometheus.MustRegister(Http2XXcode)
	prometheus.MustRegister(Http4XXcode)
	fmt.Println("Prometheues metrics registered...")
}

var con *conf

func main() {
	con = new(conf)
	con.getConf()
	fmt.Println("Starting Http server...")
	http.HandleFunc("/post", parseGhPost)
	http.HandleFunc("/", indexPage)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		fmt.Println("Failed to make connection" + err.Error())
	}
}