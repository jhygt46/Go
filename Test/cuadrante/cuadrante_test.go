package cuad

import (
	"fmt"
	//"github.com/JamesMilnerUK/pip-go"
	jsoniter "github.com/json-iterator/go"
	//"encoding/json"
	"testing"
)

type Empresa1 struct {
	P Posicion  `json:"P"`
	I []ProdsId `json:"I"`
}
type Posicion struct {
	A float32 `json:"A"`
	N float32 `json:"N"`
}
type ProdsId struct {
	I uint32 `json:"I"`
	P uint32 `json:"P"`
}

/*
type Cuadrantes struct {
	Minlat             float32 `json:"Minlat"`
	Maxlat             float32 `json:"Minlat"`
	Minlng             float32 `json:"Minlat"`
	Maxlng             float32 `json:"Minlat"`
	DimensionCuadrante float32 `json:"Minlat"`
}
*/
var (
	//Cuad = Cuadrantes{Minlat: 30.0, Maxlat: 34.0, Minlng: 70.0, Maxlng: 76.0, DimensionCuadrante: 0.25}
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func main() {
	fmt.Println("HOLA MUNDO")
}
func x(m []uint8) {
	//fmt.Println(m)
}
func y(m ProdsId) {
	//fmt.Println(m)
}

func BenchmarkTestMarshal(b *testing.B) {
	emp := Empresa1{P: Posicion{A: 33.0, N: 67}, I: []ProdsId{ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}}}
	for i := 0; i < b.N; i++ {
		u, _ := json.Marshal(emp)
		x(u)
	}
}
func BenchmarkTestUnMarshal(b *testing.B) {

	emp := Empresa1{P: Posicion{A: 33.0, N: 67}, I: []ProdsId{ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}, ProdsId{I: 3465, P: 45875}}}
	u, _ := json.Marshal(emp)
	res := Empresa1{}

	//fmt.Println("BYTE: ", u[0])

	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(u, &res); err == nil {
			for _, v := range res.I {
				y(v)
			}
		} else {
			fmt.Println(err)
		}
	}
}

/*
func BenchmarkTestCuad1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get_Cuad(31.75, 75.50)
	}
}
func BenchmarkTestCuad2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get_Cuad(30.0, 70.0)
	}
}
func BenchmarkTestCuad3(b *testing.B) {

	rectangle := pip.Polygon{
		Points: []pip.Point{
			pip.Point{X: 1.0, Y: 1.0},
			pip.Point{X: 1.0, Y: 2.0},
			pip.Point{X: 2.0, Y: 2.0},
			pip.Point{X: 2.0, Y: 1.0},
		},
	}
	pt1 := pip.Point{X: 1.1, Y: 1.1}

	for i := 0; i < b.N; i++ {
		pip.PointInPolygon(pt1, rectangle)
	}
}
func Get_Cuad(lat, lng float32) uint16 {
	x := (lat - Cuad.Minlat) / Cuad.DimensionCuadrante
	y := ((lng - Cuad.Minlng) / Cuad.DimensionCuadrante) * (Cuad.Maxlat - Cuad.Minlat) / Cuad.DimensionCuadrante
	return uint16(x + y)
}
*/
//go test -benchmem -bench=.
