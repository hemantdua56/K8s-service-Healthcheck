package main

import(
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"gopkg.in/gomail.v2"
	"log"
	"strconv"
	"strings"
	"net/http"
	"fmt"
	"time"
	"sort"
	"sync"
	"html/template"
	c "healthcheck/src/config"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	err error
	templ *template.Template
)

type Data struct {
	Microservice string
	Build string
	ReadyStatus string
	Healthcode string
	Healthstatus string
}

var dt []Data
var td []Data
var client *http.Client


func worker(Services c.Service ){

	x := make(map[string]int)

   //infinite cycle of getting health every 15 seconds
	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
		   select {
			case <- ticker.C:

					  td=nil

					  config, err := rest.InClusterConfig()
					  if err != nil {
						  panic(err.Error())
					  }
					  clientset, err := kubernetes.NewForConfig(config)
					  if err != nil {
						  log.Fatal(err)
						  fmt.Println(err.Error())
					  }


					  //Multi threading all the request for all microservices
					  var wg sync.WaitGroup
					  wg.Add(len(Services.Microservice))

					  for i := 0; i < len(Services.Microservice); i++ {
					  go func(i int) {
						  defer wg.Done()
		  
						  
						 
						
						//   fmt.Println(t)

						  //getting build version
						  deployment, err := clientset.AppsV1().Deployments(viper.GetString("namespace")).Get(Services.Microservice[i].Service,metav1.GetOptions{})
							  if err != nil {
										fmt.Println(err.Error())
										fallback, errr := clientset.AppsV1().Deployments(viper.GetString("namespace")).Get("healthcheck",metav1.GetOptions{})
                                          if errr!= nil{
                                            panic(errr)
                                          }
                                         deployment = fallback

							  }
							readyReplicas := strconv.Itoa(int(deployment.Status.ReadyReplicas))
							replicas := strconv.Itoa(int(deployment.Status.Replicas))

							ready := readyReplicas+"/"+replicas							  


							tag := deployment.Spec.Template.Spec.Containers[0].Image
							token := strings.Split(tag, ":")
							build := token[1]
							if err != nil {
									build = "deployment not found"  
									ready = "0/0"
										}




						  // generating health and ready URL
						  healthUrl := "http://"+Services.Microservice[i].Service+":"+Services.Microservice[i].Port+"/"+Services.Microservice[i].Context_Path+Services.Microservice[i].Path

						  //calling health and ready status
						  health_response := healthCheck(healthUrl)

						  fmt.Println("http://"+Services.Microservice[i].Service+":"+Services.Microservice[i].Port+"/"+Services.Microservice[i].Context_Path+Services.Microservice[i].Path+" :" +health_response[1])

						  //check if health is down for email notification
						  if health_response[0] != "200" {
							x[Services.Microservice[i].Service]=x[Services.Microservice[i].Service]+1
							fmt.Println(x[Services.Microservice[i].Service])
						  } else {
							x[Services.Microservice[i].Service]=0
						  }

						  //
						  if x[Services.Microservice[i].Service]==60 && viper.GetString("alert")=="true" {
								  m := gomail.NewMessage()
								  m.SetHeader("From", viper.GetString("sender"))
								  m.SetHeader("To", viper.GetString("recipientList"))
								  m.SetHeader("Subject", "DOWN DOWN DOWN!")
								  m.SetBody("text/html", fmt.Sprintf("Hello %s is Down on %s environment!", Services.Microservice[i].Service, viper.GetString("namespace")))

								  d := gomail.NewDialer("email-smtp.us-east-1.amazonaws.com", 587, viper.GetString("smtpUser"), viper.GetString("smtpPass"))

								  // Send the email to Pearson-GLP-Realsteel@globallogic.com
								  if err := d.DialAndSend(m); err != nil {
									  panic(err)
								  }
								  x[Services.Microservice[i].Service]=-120
						  }

						  //rendering output
						  td = append(td, Data{Microservice: Services.Microservice[i].Service, Build: build, ReadyStatus: ready, Healthcode: health_response[0], Healthstatus: health_response[1]})

						  }(i)
						}
						wg.Wait()
						fmt.Println("Main: Completed")

						//sorting the results
						sort.SliceStable(td, func(i, j int) bool {
						return td[i].Microservice < td[j].Microservice
						})
						fmt.Println(td)
						dt=td
			case <- quit:
				ticker.Stop()
				return
			}
		}
  }()
}

//function to request health
func healthCheck(url string) [2]string {


    var resp_string [2]string
    resp, err := client.Get(url)
        if err != nil {
          resp_string[0] = "500"
          resp_string[1] = "Timeout"
           return resp_string
         }
    defer resp.Body.Close()
    response := resp.StatusCode
    responsetext := http.StatusText(resp.StatusCode)
    resp_string[0] = strconv.Itoa(response)
    resp_string[1] = responsetext
    return resp_string
}
