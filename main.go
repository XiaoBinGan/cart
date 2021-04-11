package main

import (
	"github.com/XiaoBinGan/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	"github.com/opentracing/opentracing-go"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	 service2 "github.com/XiaoBinGan/cart/domain/service"
	"github.com/XiaoBinGan/cart/domain/repository"
	"github.com/XiaoBinGan/cart/handler"
	cart "github.com/XiaoBinGan/cart/proto/cart"
)


var  QPS = 100//每秒的访问

func main()  {
	//config center
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err!=nil{
		log.Error(err)
	}
	//registry center
	consul := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	//jaeger tracing
	//初始化
	newTracer, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
    if err!=nil{
    	log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(newTracer)
	//get mysql from consul
	mysqlConfig := common.GetMysqlFromConsul(consulConfig, "mysql")
	//user gorm to connect mysql
	db, err := gorm.Open("mysql", mysqlConfig.User+":"+mysqlConfig.Pwd+"@/"+mysqlConfig.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err!=nil{
		log.Error(err)
	}
	defer db.Close()
	//禁止副表
	db.SingularTable(true)
	//init table
	//if err = repository.NewCartRepository(db).InitTable();err!=nil{
	//	log.Error(err)
	//}
	//micro New service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		micro.Address("0.0.0.0:8087"),
		//registry center
		micro.Registry(consul),
		//jaeger tracing
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)
	//service init
	service.Init()

	//new data cart service
	CartDataService := service2.NewCartDataService(repository.NewCartRepository(db))
	//new handler.Cart
	if err := cart.RegisterCartHandler(service.Server(), &handler.Cart{CartDataService: CartDataService});err!=nil{
		log.Error(err)
	}


	//service run
	if err := service.Run();err!=nil{
		log.Error(err)
	}
}