package config

func initDevConf() {
	DBConf = dbConf{
		MySQL:    "root:1111112@tcp(localhost:3306)/golang_im?charset=utf8&parseTime=true&loc=Local",
		RedisIP:  "127.0.0.1:6379",
		RedisPwd: "",
	}

	WSConf = wsConf{
		WSListenAddr: ":9091",
	}

	GRPCConf = grpcConf{
		GRPCListenAddr: ":9092",
	}
}
