# logstash_prom_exporter

A tool to export custom metrics from logs in Logstash to Prometheus.

## Building

    go build

## Usage

    - ./logstash_prom_exporter
    
    - Configure logstash to send the log in json format to logstash_prom_exporter:
      
      Eg logstash output configuraion: 
        
       output {
       
        elasticsearch {
           hosts => ["localhost:9200"]
           ...
           }


      	http {
      	    url => "http://localhost:8000/post"
          	    http_method => "post"
        	   }
        }
        
      But make sure you're getting it in a jason format!
      
    - A log message should contain an element called log_message, which is an array of key value pairs.
    
      Eg:
        
        {"time": "01/Sep/2016:09:08:12 +0000", "remote_addr": "192.168.1.150", "remote_user": "-", "request": "GET /v1/a/b HTTP/1.1", "request_method": "GET", "http_referrer": "http://192.168.1.200/blahblah", "http_user_agent": "Mozilla/5.0 (X11; Linux x86_64)","process_name" : "ab_service","host_name" : "abserver","thread_id" : "undefined","log_message" : [{"key" : "body_bytes_sent", "value" : "86"},{"key" : "request_count", "value" : "true"},{"key" : "status", "value" : "200"}, {"key" : "request_time", "value" : "0.014"},{"key" : "upstream_response_time", "value" : "0.013"}]}
      
    - The key will be metric's name and value will be it's value.
   
    - If the metric is a counter, value must be "true"
    
    - Test by sending curl request to /metrics of the server
      
      Eg:
      
        curl localhost:8000/metrics
    
    - Specify the regex for the http endpoint if you want to track HTTP metrics on prometheus in logex.yml
    
 
     