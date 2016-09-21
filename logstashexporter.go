package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
	"github.com/deckarep/golang-set"
	"net/http"
	"strconv"
	"strings"
	"regexp"
	"io/ioutil"
	"log"
)

//Structure of config file
type conf struct {
	Regex string `yaml:"regex"`
	Metrics []Metric `yaml:"metrics"`
}

//Structure of a metric in config
type Metric struct{
	KeyName string `yaml:"keyname"`
	MetricType string `yaml:"metrictype"`
	MetricName string `yaml:"metricname"`
	Help string `yaml:"help"`
	Labels []string `yaml:"labels"`
}

type MapStringInterface map[string]interface{}

//Function to get configuration from yaml
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

//Function to handle a log message received on /post
func parseGhPost(rw http.ResponseWriter, request *http.Request) {
	//A temporary map to store the values for each label
	var labelvaluemap = make(map[string] string)
	//Decode the log message from JSON
	var mapstringinterface MapStringInterface
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&mapstringinterface)
	if err != nil {
		fmt.Println("Error JSON Decoding!", err.Error())
	}

	//Iterate through all labels and check if it is present in the log entry
	for _, label := range AllLabels.ToSlice() {
		if value, ok := mapstringinterface[label.(string)]; ok {
			if _, ok := value.([]interface{}); ok {
				//Array of interfaces
			}else{
				//Adding label to the label-value map for this log entry
				labelvaluemap[label.(string)] = value.(string)
			}
		}else if value, ok := mapstringinterface["request"]; ok {
			//Find out the API endpoint hit using the regex mentioned in "logex.yml"
			var RequestEndpoint string
			Regexp := regexp.MustCompile(con.Regex)
			if(string(value.(string))==""){
				//Request field empty
			}else
			{
				RequestEndpoint = Regexp.FindString(string(value.(string)))
			}
			RequestEndpoint = strings.Replace(RequestEndpoint, " ", "_", -1)
			labelvaluemap["request_endpoint"] = RequestEndpoint

		}else{
				//Label not present
		}
	}

	if value, ok := mapstringinterface["log_message"]; ok {
		if element, ok := value.([]interface{}); ok {
			//Array of interfaces
			for _, message := range element {
				if blahMap, ok := message.(map[string]interface{}); ok {
					var templabelarray []string
					for _, lname:= range LabelsMap[blahMap["key"].(string)]{
						templabelarray = append(templabelarray,labelvaluemap[lname])
					}
					if _, ok := CounterMap[blahMap["key"].(string)]; ok {
						CounterMap[blahMap["key"].(string)].WithLabelValues(templabelarray...).Inc()
				}else if _, ok := GaugeMap[blahMap["key"].(string)]; ok {
					setvalue, err := strconv.ParseFloat(blahMap["value"].(string), 64)
					if err != nil {
						panic(err)
					}
					GaugeMap[blahMap["key"].(string)].WithLabelValues(templabelarray...).Set(setvalue)
				}
				}else {
					fmt.Println("\nNot recognized as a monitoring message")
				}
			}
		}
	}
}

//Configuration object
var con *conf
//A map of Prometheus counters
var CounterMap = make(map[string] *prometheus.CounterVec)
//A map of Prometheus Gauges
var GaugeMap = make(map[string] *prometheus.GaugeVec)
//A map to store the labels each metric is associated with
var LabelsMap = make(map[string] []string)
//A string slice of all the keys names
var KeyNames []string
//A set to store all the labels that are there in the config
var AllLabels = mapset.NewSet()

func main() {
	//Get configuration
	con = new(conf)
	con.getConf()

	fmt.Println("Registering Prometheus Metrics...")
	//Traverse through metrics in "logex.yml" and Create Prometheus Counters and Gauges respectively
	for _,element := range con.Metrics {
		if (element.MetricType=="counter"){
			CounterMap[element.KeyName] = prometheus.NewCounterVec(prometheus.CounterOpts{Name : element.MetricName,Help: element.Help}, element.Labels)
			prometheus.MustRegister(CounterMap[element.KeyName])
			LabelsMap[element.KeyName]= element.Labels
			KeyNames=append(KeyNames,element.KeyName)
			old := element.Labels
			new := make([]interface{}, len(old))
			for i, v := range old {
				new[i] = v
			}
			AllLabels = AllLabels.Union(mapset.NewSetFromSlice(new))
		}else if(element.MetricType=="gauge"){
			GaugeMap[element.KeyName] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name : element.MetricName, Help: element.Help}, element.Labels)
			prometheus.MustRegister(GaugeMap[element.KeyName])
			LabelsMap[element.KeyName]= element.Labels
			KeyNames=append(KeyNames,element.KeyName)
			old := element.Labels
			new := make([]interface{}, len(old))
			for i, v := range old {
				new[i] = v
			}
			AllLabels = AllLabels.Union(mapset.NewSetFromSlice(new))
		}
	}

	fmt.Println("Metrics successfully registered!")
	fmt.Println("Starting Http server and listening on port 8000...")
	http.HandleFunc("/post", parseGhPost)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		fmt.Println("Failed to make connection" + err.Error())
	}
}