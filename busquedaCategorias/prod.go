package main

import (
	"os"
	"fmt"
	"time"
	"runtime"
	//"strconv"
	"io/ioutil"
	"encoding/json"
    "github.com/valyala/fasthttp"
    "github.com/dgraph-io/ristretto"
)



// OBJETO CACHE CATCUAD
type CatCuad struct {
	Hits int `json:"Hits"`
	Date time.Time `json:"Date"`
	Locales []Local `json:"Local"`
}
type Locales struct {
	info *InfoLocal
	prodscache []*ProductosCache
	prodsunicos []ProductosUnicosLocal
}
type InfoLocal struct {
	idloc uint32 `json:"idloc"`
	lat float32 `json:"lat"`
	lng float32 `json:"lng"`
	nombre string `json:"nombre"`
}
type ProductosUnicosLocal struct {
	idpro uint32 `json:"idpro"`
	precio uint32 `json:"precio"`
	Calidad uint8 `json:"Calidad"`
	Campos [][]uint16 `json:"Campos"`
	Evals [][]uint16 `json:"Evals"`
}
type ProductosCache struct {
	idpro uint32 `json:"idpro"`
	precio uint32 `json:"precio"`
	Calidad uint8 `json:"Calidad"`
	Campos [][]uint16 `json:"Campos"`
	Evals [][]uint16 `json:"Evals"`
}



// OBJETO RESPUESTA
type Respuesta struct {
	Locales []ResLocal `json:"Locales"`
	Productos []ResProductos `json:"Productos"`
}
type ResLocal struct {
	Idloc int32 `json:"Idloc"`
	Lat float32 `json:"Lat"`
	Lng float32 `json:"Lng"`
	Nombre float32 `json:"Nombre"`
}
type ResProductos struct {
	Idpro uint32 `json:"Idpro"`
	Precio uint32 `json:"Precio"`
	Calidad uint8 `json:"Calidad"`
	Puntaje uint8 `json:"Puntaje"`
	Campos []uint16 `json:"Campos"`
	Evals []uint16 `json:"Evals"`
}

// OBJETO PARA LEER DE ARCHIVO
type CatCuadFile struct {
	Locales []Local `json:"Locales"`
}
type Local struct {
	Idloc int32 `json:"Idloc"`
	Lat float32 `json:"Lat"`
	Lng float32 `json:"Lng"`
	Nombre string `json:"Nombre"`
	Prods []ProductosCompartidoLocal `json:"Prods"`
}
type ProductosCompartidoLocal struct {
	Idpro uint32 `json:"Idpro"`
	Precio uint32 `json:"Precio"`
	Calidad uint8 `json:"Calidad"`
	Campos [][]uint16 `json:"Campos"`
	Evals [][]uint16 `json:"Evals"`
}
type CacheHits struct {
	Hits int `json:"Hits"`
	Date time.Time `json:"Date"`
}

type MyHandler struct {
	cache *ristretto.Cache
}



var ListaProductos []ProductosCache
var ListaLocales []InfoLocal

func main() {

	//cacheMode := true
	balancerMode := true

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil { panic(err) }

	if balancerMode {
		// CONTAR POR RAGOS
	}

	pass := &MyHandler{ cache: cache, }
    fasthttp.ListenAndServe(":80", pass.HandleFastHTTP)
	
}
/*
func Process(id string, res *Respuesta){

	cacheMode := true
	value, found := h.cache.Get("c"+sid)
	if found {
		if cacheMode {
			// ACTUALIZA NUMERO CACHE
		}
		// MUESTRA CACHE ENCONTRADO
		fmt.Println(value)
	}else{
		if cacheMode {
			// CACHE AUTOMATICO
			value, found := h.cache.Get("a"+sid)
			if !found {
				//SE CREA CACHE
				h.cache.SetWithTTL("a"+sid, CacheHits{ 1, time.Now() }, 1, 10*time.Second)
				readCategoriaCuad(127, 186, false)
			}else{
				if value.(CacheHits).Hits < 1235 {
					//SE ACTUALIZA
					//fmt.Println(value.(CacheHits).Hits)
					h.cache.SetWithTTL("a"+sid, CacheHits{ Hits: 1236, Date: time.Now() }, 1, 10*time.Second)
					readCategoriaCuad(127, 186, false)
				}else{
					//LEE ARCHIVO 127/186 Y LO MUSTRA Y GUARDA EN CACHE
					readCategoriaCuad(127, 186, true)
					h.cache.Del("a"+sid)
				}
			}

		}else{
			readCategoriaCuad(127, 186, false)
		}
	}
	
	id := "c127-186"
	value, found := cache.Get(id)
	if found {
		if cacheMode {
			// ACTUALIZA NUMERO CACHE
		}
		// MUESTRA CACHE ENCONTRADO
		fmt.Println(value)
	}else{
		if cacheMode {
			// CACHE AUTOMATICO
			idw := "a127-186"
			value, found := cache.Get(idw)
			if !found {
				//SE CREA CACHE
				cache.SetWithTTL(id, CacheHits{ 1, time.Now() }, 1, 10*time.Second)
				//json.NewEncoder(w).Encode(readCategoriaCuad(127, 186, false))
			}else{
				if value.(CacheHits).Hits < 1235 {
					//SE ACTUALIZA
					//LEE ARCHIVO 127/186 Y LO MUSTRA
					//json.NewEncoder(w).Encode(readCategoriaCuad(127, 186, false))
					fmt.Println(value.(CacheHits).Hits)
					fmt.Println(value.(CacheHits).Date)
					cache.SetWithTTL(idw, CacheHits{ Hits: 1, Date: time.Now() }, 1, 10*time.Second)

				}else{
					//LEE ARCHIVO 127/186 Y LO MUSTRA Y GUARDA EN CACHE
					filecache := readCategoriaCuad(127, 186, true)
					fmt.Println(filecache)
					
					newcache := &CatCuad{ Hits: 0, Date: time.Now() }
					for i := range filecache.Locales {
						newcache.Locales{ info: getPointerLocal(filecache.Locales[i]) }
						for j := range filecache.Locales[j].prodscache {

						}
						for j := range filecache.Locales[j].prodsunicos {

						}
					}
					cache.Set("127/186", newcache, 1)
					
					//json.NewEncoder(w).Encode(newcache)
				}

			}
		}
	}
	

}
*/

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	
	var res Respuesta
	cat := string(ctx.QueryArgs().Peek("cat"))
	cuad := string(ctx.QueryArgs().Peek("cuad"))
	
	for _, cuads := range getCuads(cuad){
		value, found := h.cache.Get("c"+cat+"-"+cuads)
		if found {
			fmt.Println(value)
		}else{
			readCategoriaCuad(cat, cuads, false, &res)
		}
	}
	
	json.NewEncoder(ctx).Encode(res)
    //fmt.Fprintf(ctx, "Ok.");
	
}
func getCuads(cuad string) [2]string{

	return [2]string{"10", "25"}

}
func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func readCategoriaCuad(cat string, cuad string, cache bool, res *Respuesta) {

	file := "catcuad/"+cat+"/"+cuad+".json"
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var obj CatCuadFile
	json.Unmarshal(byteValue, &obj)
	fmt.Println(obj)
	

	if cache {
		// CREA OBJETO CACHE
		// PUNTERO DE LOCALES Y PRODUCTOS 
	}


}
/*
func getPointerLocal(local InfoLocal) *InfoLocal{
	for i := range ListaLocales {
		if ListaLocales[i].Idloc == local.Idloc {
			return &ListaLocales[i]
		}
	}
	ListaLocales = append(ListaLocales, local)
	return &ListaLocales[len(ListaLocales) - 1]
}
func getPointerProducto(pro ProductosCache) *ProductosCache{
	for i := range ListaProductos {
		if ListaProductos[i].Idpro == pro.Idpro {
			return &ListaProductos[i]
		}
	}
	ListaProductos = append(ListaProductos, pro)
	return &ListaProductos[len(ListaProductos) - 1]
}
*/
