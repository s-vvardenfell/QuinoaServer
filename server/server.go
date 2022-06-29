package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"github.com/s-vvardenfell/QuinoaServer/config"
	gen "github.com/s-vvardenfell/QuinoaServer/generated"
	"github.com/sirupsen/logrus"
)

type FilmsToCache struct {
	Films []DataToCache `json:"films"`
}

type DataToCache struct {
	Name string `json:"name"`
	Ref  string `json:"ref"`
	Img  string `json:"img"`
}

type QuinoaMainServer struct {
	gen.UnimplementedMainServiceServer
	rc *RedisClient
	pc *ParserClient
}

func NewServer(cnfg config.Config) *QuinoaMainServer {
	return &QuinoaMainServer{
		rc: NewRedisClient(cnfg.ServerHost, cnfg.RedisServPort),
		pc: NewParserClient(cnfg.ServerHost, "50053"), //TODO занести в конфиг
	}
}

func (q *QuinoaMainServer) GetParsedData(
	ctx context.Context, in *gen.Conditions) (*gen.ParsedResults, error) {
	// res := generated.ParsedResult{
	// 	Data: []*generated.ParsedData{{Name: "FilmName", Ref: "FilmRef", Img: "SmallImage"}},
	// }

	// вычисляем хэш
	// hash := getMD5Hash(in.String())
	// fmt.Println("HASH: ", hash)

	// // hash = "d1ec4ea6edb766e3cde6534ea442f921" //TEST
	// // проверяем кэшированые в редисе данные по этому ключу-хэшу
	// cachedRes, err := q.rc.Get(ctx, &gen.Key{Key: hash})

	// fmt.Println("DATA FROM REDIS:<<<<", cachedRes, ">>>>")

	// // если ошибка, выводим в консоль и работаем дальше
	// if err != nil {
	// 	logrus.Warningf("got error from redis while Get(), %v", err)
	// } else { //если ошибки нет и есть данные в редисе
	// 	// здесь размаршалливаем в структуру FilmsToCache,
	// 	// конвертим в ParsedResults и возвращаем
	// 	return &gen.ParsedResults{
	// 		Data: []*gen.ParsedData{{
	// 			Name: "STUB-Name response from redis",
	// 			Ref:  "STUB-Ref response from redis",
	// 			Img:  "STUB-Img response from redis",
	// 		},
	// 		},
	// 	}, nil
	// }

	////////////////////////спрашиваем у парсера///////////////////////////
	res, err := q.pc.ParseData(ctx, in)
	if err != nil {
		logrus.Errorf("got error from parser, %v", err)
		return nil, err
	}
	//////////////////////////////////////////////////////////////////////

	// pd := make([]*gen.ParsedData, 0)
	// pd = append(pd, &gen.ParsedData{ //ЗАГЛУШКА
	// 	Name: "Имя: " + in.Type + in.Genres[0],
	// 	Ref:  "Ссылка: " + in.StartYear + in.EndYear,
	// 	Img:  "Картинка: " + in.Keyword + in.Countries[0]})

	// res := gen.ParsedResults{
	// 	Data: pd,
	// }

	// кэшируем в редис
	// маршаллим в json, чтобы потом удобно получить
	// var cache FilmsToCache
	// cache.Films = make([]DataToCache, 0, len(pd))

	// for i := range pd {
	// 	cache.Films = append(cache.Films, DataToCache{
	// 		Name: pd[i].Name,
	// 		Ref:  pd[i].Ref,
	// 		Img:  pd[i].Img,
	// 	})
	// }

	// jdata, err := json.Marshal(cache)
	// if err != nil {
	// 	logrus.Errorf("cannot convert FilmsToCache struct to json to store in Redis")
	// 	//TODO а как обрабатываетс???
	// }
	// // кэшируем
	// ok, err := q.rc.Set(ctx, &gen.Input{Key: hash, Val: string(jdata)})
	// if err != nil {
	// 	fmt.Printf("%T", err)
	// 	logrus.Warningf("got error from redis while Set(), %v", err)
	// }

	// проверяем, что записано
	// if !ok.Ok {
	// 	logrus.Warningf("it's not OK, %v", err)
	// } else {
	// 	logrus.Warningf("it's OK, %v", err)
	// }

	return res, nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
