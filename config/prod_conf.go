package config

func initProdConf() {
	DBConf = dbConf{
		MySQL:    "root:CFkjAAbb1234@tcp(124.70.134.101:3306)/gim?charset=utf8&parseTime=true&loc=Local",
		RedisIP:  "124.70.134.101:6379",
		RedisPwd: "CFKJAAbb1234",
	}

	WSConf = wsConf{
		WSListenAddr: ":8081",
	}

	GRPCConf = grpcConf{
		GRPCListenAddr: ":9092",
	}
}
