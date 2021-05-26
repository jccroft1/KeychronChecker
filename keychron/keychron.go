package keychron

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocolly/colly"
	"github.com/jccroft1/KeychronChecker/storage"
)

const (
	BackOrder           string = "http://schema.org/BackOrder"
	Discontinued        string = "http://schema.org/Discontinued"
	InStock             string = "http://schema.org/InStock"
	InStoreOnly         string = "http://schema.org/InStoreOnly"
	LimitedAvailability string = "http://schema.org/LimitedAvailability"
	OnlineOnly          string = "http://schema.org/OnlineOnly"
	OutOfStock          string = "http://schema.org/OutOfStock"
	PreOrder            string = "http://schema.org/PreOrder"
	PreSale             string = "http://schema.org/PreSale"
	SoldOut             string = "http://schema.org/SoldOut"
)

type Product struct {
	// Type         string `json:"@type"`
	// Price        string `json:"price"`
	SKU          string `json:"sku"`
	URL          string `json:"url"`
	Availability string `json:"availability"`
	Name         string `json:"name"`
}

type ProductInfo struct {
	// Type     string `json:"@type"`
	Name     string    `json:"name"`
	Products []Product `json:"offers"`
}

var (
	delay time.Duration = 15 * time.Minute

	allKeyboards = []string{
		"https://www.keychron.com/products/keychron-k12-wireless-mechanical-keyboard",
		"https://www.keychron.com/products/keychron-k1-wireless-mechanical-keyboard",
		"https://www.keychron.com/products/keychron-k2-wireless-mechanical-keyboard",
		"https://www.keychron.com/products/keychron-k2-hot-swappable-wireless-mechanical-keyboard",
		"https://www.keychron.com/products/keychron-k3-wireless-mechanical-keyboard",
		"https://www.keychron.com/products/keychron-k4-wireless-mechanical-keyboard-version-2",
		"https://www.keychron.com/products/keychron-k6-wireless-mechanical-keyboard",
		"https://www.keychron.com/products/keychron-k8-tenkeyless-wireless-mechanical-keyboard",
		"https://www.keychron.com/products/keychron-c1-wired-mechanical-keyboard",
		"https://www.keychron.com/products/keychron-c2-wired-mechanical-keyboard",
		"https://www.keychron.com/products/keychron-k1-wireless-mechanical-keyboard-japan-jis-layout-version-4",
	}

	Alert = func(message string) error {
		fmt.Println(message)
		return nil
	}
)

func newData(e *colly.HTMLElement) {
	var info ProductInfo
	if err := json.Unmarshal([]byte(e.Text), &info); err != nil {
		return
	}

	if len(info.Products) == 0 {
		return
	}

	for _, p := range info.Products {
		addProduct(info.Name, p)
	}
	return
}

func Start() {
	readConfig(&productCache)

	s := storage.NoCacheStorage{}
	c := colly.NewCollector()
	c.SetStorage(&s)

	c.OnHTML(`head > script[type="application/ld+json"]`, newData)

	for _, link := range allKeyboards {
		c.Visit(link)
		time.Sleep(time.Duration(1) * time.Second)
	}

	ticker := time.NewTicker(delay)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for _, link := range allKeyboards {
					c.Visit(link)
					time.Sleep(time.Duration(5) * time.Second)
				}
			}
		}
	}()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	<-term
	writeConfig(&productCache)
}
