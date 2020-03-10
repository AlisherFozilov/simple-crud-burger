package app

import (
	"github.com/AlisherFozilov/crud/pkg/crud/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func (receiver *server) handleBurgersList() func(http.ResponseWriter, *http.Request) {
	tpl, err := template.ParseFiles(filepath.Join(receiver.templatesPath, "index.gohtml"))
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := receiver.burgersSvc.BurgersList()
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			Burgers []struct {
				Id        int64
				Name      string
				BeforeDot int
				AfterDot  int
			}
		}{
			Title: "McDonalds",
		}

		listWithConvertedPrices := make([]struct {
			Id        int64
			Name      string
			BeforeDot int
			AfterDot  int
		}, len(list))

		for index, _ := range list {
			listWithConvertedPrices[index].Id = list[index].Id
			listWithConvertedPrices[index].Name = list[index].Name
			listWithConvertedPrices[index].BeforeDot = list[index].Price / 100
			listWithConvertedPrices[index].AfterDot = list[index].Price % 100
		}

		data.Burgers = listWithConvertedPrices
		err = tpl.Execute(writer, data)
		if err != nil {
			log.Printf("can't execute template with data %v, %s", data, err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (receiver *server) handleBurgersSave() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			log.Print(err)
			return
		}

		name, ok := request.PostForm["name"]
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			log.Print("bad request: ", request.RemoteAddr)
			return
		}

		price, ok := request.PostForm["price"]
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			log.Print("bad request: ", request.RemoteAddr)
			return
		}

		dotIndex := strings.Index(price[0], ".")
		var priceString string
		if dotIndex != -1 {
			priceString = price[0][:dotIndex] + price[0][dotIndex+1:]
		} else {
			priceString = price[0] + "00"
		}
		priceInt, err := strconv.Atoi(priceString)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			log.Print("bad request: ", request.RemoteAddr)
			return
		}
		if priceInt <= 0 {
			_, err := writer.Write([]byte("price can't be zero or less"))
			if err != nil {
				log.Println("can't write to connection: ", err)
			}
			log.Printf("price can't be zero or less")
			return
		}

		burger := models.Burger{
			Id:      0,
			Name:    name[0],
			Price:   priceInt,
			Removed: false,
		}
		err = receiver.burgersSvc.Save(burger)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 400 - Bad Request
			log.Printf("can't save burger %v: %s", burger, err)
			return
		}
		// TODO: посмотреть, можно ли переделать на GET
		http.Redirect(writer, request, "/", http.StatusMovedPermanently)
		return
	}
}

func (receiver *server) handleBurgersRemove() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		err := request.ParseForm()
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		id, ok := request.PostForm["id"]
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		idInt, err := strconv.Atoi(id[0])
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}

		err = receiver.burgersSvc.RemoveById(int64(idInt))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 400 - Bad Request
			log.Print(err)
			return
		}
		// TODO: посмотреть, можно ли переделать на GET
		http.Redirect(writer, request, "/", http.StatusMovedPermanently)
		return
	}
}

func (receiver *server) handleFavicon() func(http.ResponseWriter, *http.Request) {
	// TODO: handle concurrency
	file, err := ioutil.ReadFile(filepath.Join(receiver.assetsPath, "favicon.ico"))
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(file)
		if err != nil {
			log.Print("can't send favicon: ", err)
		}
	}
}

func (receiver server) handleNotFound() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
