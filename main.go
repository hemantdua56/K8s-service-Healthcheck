package main
import(
      "k8s.io/client-go/kubernetes"
      "k8s.io/client-go/rest"
      "gopkg.in/gomail.v2"
      "log"
      "strconv"
      "strings"
      "net/http"
      "os"
      "fmt"
      "time"
      "sort"
      "sync"
      "html/template"
      "encoding/json"
      metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
      err error
      templ *template.Template
)

type Data struct {
      Microservice string
      Build string
      Healthcode string
      Healthstatus string
}


func main() {
      //calling worker function which regularly hits health for each microservice
      worker()

      //initiate client
  		initiate()

      //HTTP handler
      http.HandleFunc("/", hello)

      fmt.Printf("Starting server for testing HTTP POST...\n")
      if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
       }
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


func worker(){

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
                        //list of microservices
                        microservices := []string{"<service list seperated by a ,>"}

                        //Getting namespace from environment variables
                        Namespace := os.Getenv("namespace")
                        user := os.Getenv("smtp-user")
                        password := os.Getenv("smtp-pass")
                        alert := os.Getenv("alert")

                        //Multi threading all the request for all microservices
                        var wg sync.WaitGroup
                        wg.Add(len(microservices))

                        for i := 0; i < len(microservices); i++ {
                        go func(i int) {
                            defer wg.Done()
                            svc := microservices[i]

                            //getting service port
                            service, err := clientset.CoreV1().Services(Namespace).Get(svc,metav1.GetOptions{})
                             if err != nil {
                                       fmt.Println(err.Error())
					   service, err = clientset.CoreV1().Services(Namespace).Get("lap",metav1.GetOptions{})
                                          }
                            po := service.Spec.Ports[0].Port
                            port := int64(po)
                            t := strconv.FormatInt(port,10)
                            fmt.Println(t)

                            //getting build version
                            deployment, err := clientset.AppsV1().Deployments(Namespace).Get(svc,metav1.GetOptions{})
                                if err != nil {
                                         fmt.Println(err.Error())
                                          fallback, errr := clientset.AppsV1().Deployments(Namespace).Get("lap",metav1.GetOptions{})
                                          if errr!= nil{
                                            panic(errr)
                                          }
                                          deployment = fallback
                                }
                            tag := deployment.Spec.Template.Spec.Containers[0].Image
                            s := strings.Split(tag, ":")
                            build := s[1]
                             if err != nil {
									 build = "deployment not found"
                                        }
                            fmt.Println(build)


                            // generating health and ready URL
                            healthUrl := "http://"+svc+":"+t+"/"+svc+"/health"

                            //calling health and ready status
                            health_response := healthCheck(healthUrl)

                            fmt.Println("http://"+svc+":"+t+"/"+svc+"/health :" +health_response[1])

                            //check if health is down for email notification
                            if health_response[0] != "200" {
                              x[svc]=x[svc]+1
                              fmt.Println(x[svc])
                            } else {
                              x[svc]=0
                            }

                            //
                            if x[svc]==60 && alert=="true" {
                                    m := gomail.NewMessage()
                                    m.SetHeader("From", "alertmanager@alerts.com")
                                    m.SetHeader("To", "")
                                    m.SetHeader("Subject", "")
                                    m.SetBody("text/html", fmt.Sprintf("Hello %s is Down on %s environment!", svc,Namespace))

                                    d := gomail.NewDialer("email-smtp.us-east-1.amazonaws.com", 587, user, password)

                                    // Send the email to Pearson-GLP-Realsteel@globallogic.com
                                    if err := d.DialAndSend(m); err != nil {
                                        panic(err)
                                    }
                                    x[svc]=-120
                            }

                            //rendering output
                            td = append(td, Data{Microservice: svc, Build: build, Healthcode: health_response[0], Healthstatus: health_response[1]})

                            }(i)
                          }
                          wg.Wait()
                          fmt.Println("Main: Completed")
                          fmt.Println(td)

                          //sorting the results
                          sort.SliceStable(td, func(i, j int) bool {
                          return td[i].Microservice < td[j].Microservice
                          })
                          fmt.Println("here i am")
                          fmt.Println(td)
                          dt=td
              case <- quit:
                  ticker.Stop()
                  return
              }
          }
    }()
}

var dt []Data
var td []Data

var client *http.Client

func initiate() {
     tr := &http.Transport{
         MaxIdleConnsPerHost: 1024,
         TLSHandshakeTimeout: 5 * time.Second,
       	 DisableKeepAlives: true,
     }
     client = &http.Client{Transport: tr,  
                           Timeout: 10 * time.Second,}
 }


func hello(w http.ResponseWriter, r *http.Request){
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
