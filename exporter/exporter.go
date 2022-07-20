package exporter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
)

var (
	requestCpu = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_allocation_request_cpu",
			Help: "total request cpu of a node",
		},
		[]string{"id"},
	)

	requestMem = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_allocation_request_mem",
			Help: "total request mem of a node",
		},
		[]string{"id"},
	)
)

var (
	clientset *kubernetes.Clientset
)

func register() {
	prometheus.MustRegister(requestCpu)
	prometheus.MustRegister(requestMem)
}

func initClientset() {
	// creates the in-cluster config
	config, err := clientcmd.BuildConfigFromFlags("","/root/.kube/config/kubectl.kubeconfig")
	//config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func init() {
	register()
	initClientset()
}

type Exporter struct {
	addr string
}

// ServeHTTP handler
func (e *Exporter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	e.reportAllocationData()
	promhttp.Handler().ServeHTTP(writer, request)
}

func NewExporter(addr string) *Exporter {
	return &Exporter{
		addr: addr,
	}
}

// RunServer starts HTTP server loop
func (e *Exporter) RunServer() {
	http.Handle("/", http.HandlerFunc(showIndex))
	http.Handle("/metrics", e)

	log.Printf("Providing metrics at http://%s/metrics", e.addr)
	err := http.ListenAndServe(e.addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// showIndex serves index page
func showIndex(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")
	res :=
		`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
	<meta name="viewport" content="width=device-width">
	<title>Node Allocation Prometheus Exporter</title>
</head>
<body>
<h1>Node Allocation Prometheus Exporter</h1>
<p>
	<a href="/metrics">Metrics</a>
</p>
<p>
	<a href="https://github.com/justadogistaken/pmu_exporter">Homepage</a>
</p>
</body>
</html>
`
	fmt.Fprint(w, res)
}

func (*Exporter) reportAllocationData() {
	uploads := getAllocationData()
	for _, u := range uploads {
		requestCpu.WithLabelValues(u.ID).Set(u.requestCpu)
		requestMem.WithLabelValues(u.ID).Set(u.requestMem)
	}
}
