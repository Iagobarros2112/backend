// aqui vou definir o arquivo principal que vai
//  gerenciar todas as funçoes
// do projeto

package main

import (
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	orderModel "backend/internal/order/model"
	productModel "backend/internal/product/model"

	httpServer "backend/internal/server/http"
	userModel "backend/internal/user/model"
	"backend/pkg/config"
	"backend/pkg/dbs"
	"backend/pkg/redis"
)

func main() {

	// criei um ambiente virtual
	// na qual a aplicaçao vai funcionar
	cfg := config.LoadConfig()
	logger.Initialize(cfg.Environment)

	//criei um banco de dados novo e verifico se ele existe
	// se nao mostre um erro na tela

	db, err := dbs.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
	}

	//nesse banco de dados que foi criado
	//preencha ele com os parametros que defini nos models
	// se isso falhar, mostre um erro na tela

	err = db.AutoMigrate(&userModel.User{}, &productModel.Product{}, orderModel.Order{}, orderModel.OrderLine{})
	if err != nil {
		logger.Fatal("Database migration fail", err)
	}

	// iniciei um validador de dados
	// para verificar se os dados que estao sendo inseridos
	//estao corretos

	validator := validation.New()

	cache := redis.New(redis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})

	// criei uma funcao assincrona que vai rodar o servidor http
	// contendo o validador, o banco de dados e o cache do redis
	//se houver algum erro,mostre na tela

	go func() {
		httpSvr := httpServer.NewServer(validator, db, cache)
		if err = httpSvr.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

}
