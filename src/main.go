package main
import(
      "net/http"
      "fmt"
	  "time"
	  "log"
	  "encoding/json"
	  c "healthcheck/src/config"
	  "github.com/spf13/viper"
	  
	)

func main() {


	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath("/opt/")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var serviceConfig c.Service

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	
	// Set undefined variables
	viper.SetDefault("alert", "false")

	err := viper.Unmarshal(&serviceConfig)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	//calling worker function which regularly hits health for each microservice
	worker(serviceConfig)

	//initiate client
	initiate()

	//HTTP handler
	http.HandleFunc("/", request)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
	  log.Fatal(err)
	 }
}


func initiate() {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 5 * time.Second,
		   DisableKeepAlives: true,
	}
	client = &http.Client{Transport: tr,  
						  Timeout: 10 * time.Second,}
}

func request(w http.ResponseWriter, r *http.Request){
    if r.URL.Path != "/" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

      userJson, err :=json.Marshal(dt)
      if err != nil{
        panic(err)
      }
      fmt.Println(dt)
      w.Header().Set("Contetnt-Type","application/json")
      w.WriteHeader(http.StatusOK)
      w.Write(userJson)

}