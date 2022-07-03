package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

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
	// если все условия не заполнены - ошибка
	if in.Type == "" &&
		in.Genres == nil &&
		in.StartYear == "" &&
		in.EndYear == "" &&
		in.Keyword == "" &&
		in.Countries == nil {
		logrus.Error("got query with no conditions")
		return nil, errors.New("no conditions")
	}

	// вычисляем хэш условий запроса
	hash := getMD5Hash(in.String())

	// проверяем кэшированые в редисе данные по ключу-хэшу
	cachedRes, err := q.rc.Get(ctx, &gen.Key{Key: hash})
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			logrus.Infof("no cached results, %v", err)
		} else {
			logrus.Errorf("got error from redis while Get(), %v", err)
		}
	} else {
		var cachedVals FilmsToCache
		err = json.Unmarshal([]byte(cachedRes.Val), &cachedVals)
		if err != nil {
			logrus.Errorf("got error while unmarshalling cached values, %v", err)
		}

		var res gen.ParsedResults

		for i := 0; i < len(cachedVals.Films); i++ {
			res.Data = append(res.Data, &gen.ParsedData{
				Name: cachedVals.Films[i].Name,
				Ref:  cachedVals.Films[i].Ref,
				Img:  cachedVals.Films[i].Img,
			})
		}
		logrus.Info("parsing result succesfully got from cache")
		return &res, nil
	}

	//спрашиваем у парсера
	res, err := q.pc.ParseData(ctx, in)
	if err != nil {
		logrus.Errorf("got error from parser, %v", err)
		return nil, err
	}

	// кэшируем
	// преобразуем в json, чтобы потом распарсить
	var cachedVals FilmsToCache
	for i := range res.Data {
		cachedVals.Films = append(cachedVals.Films,
			DataToCache{
				Name: res.Data[i].Name,
				Ref:  res.Data[i].Ref,
				Img:  res.Data[i].Img,
			})
	}

	json, err := json.Marshal(cachedVals)
	if err != nil {
		logrus.Errorf("cannot marshall struct to json to store in cache, %v", err)
	}

	ok, err := q.rc.Set(ctx, &gen.Input{Key: hash, Val: string(json)})
	if err != nil {
		logrus.Errorf("got error from redis while Set(), %v", err)
	}

	if !ok.Ok {
		logrus.Errorf("cannot cache value, %v", err)
	} else {
		logrus.Info("parsing result succesfully cached")
	}

	return res, nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
