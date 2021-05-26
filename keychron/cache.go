package keychron

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	cachePath string = "cache/products"

	productCache = make(map[string]Product)
	productLock  = sync.Mutex{}
)

func addProduct(baseName string, newProduct Product) {
	productLock.Lock()
	defer productLock.Unlock()

	oldProduct, exist := productCache[newProduct.SKU]
	defer func() {
		if oldProduct == newProduct {
			return
		}

		productCache[newProduct.SKU] = newProduct
		go writeConfig(&productCache)
	}()

	if !exist {
		return
	}

	// check out-in stock
	switch oldProduct.Availability {
	case BackOrder, OutOfStock, PreOrder, PreSale, SoldOut:
		switch newProduct.Availability {
		case InStock, InStoreOnly, LimitedAvailability, OnlineOnly:
			Alert(fmt.Sprintf("%v %v is now in stock! %v", baseName, newProduct.Name, newProduct.URL))
		}
	}
}

func readConfig(out *map[string]Product) {
	productLock.Lock()
	defer productLock.Unlock()

	f, err := os.Open(cachePath)
	defer f.Close()
	if err != nil {
		return
	}
	dec := json.NewDecoder(f) // Will read from network.
	dec.Decode(out)
}

func writeConfig(in *map[string]Product) {
	productLock.Lock()
	defer productLock.Unlock()

	f, err := os.Create(cachePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	err = enc.Encode(in)
	if err != nil {
		fmt.Println(err)
	}
}
