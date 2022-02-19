package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/valyala/fasthttp"
)

type Result struct {
	Prods []Producto `json:"Prods"`
}
type MyHandler struct {
}
type Buffer struct {
	Buf []byte `json:"Buf"`
}
type Empresa struct {
	Lat   float32 `json:"Lat"`
	Lng   float32 `json:"Lng"`
	Prods []Prods `json:"Prods"`
}
type Prods struct {
	Id       uint64   `json:"Id"`
	Ids      uint64   `json:"Ids"`
	Precio   uint64   `json:"Precio"`
	Calidad  uint8    `json:"Calidad"`
	Cantidad uint8    `json:"Calidad"`
	Opciones []Opcion `json:"Opciones"`
	Evals    []Evals  `json:"Opciones"`
}
type Opcion struct {
	Id      uint64   `json:"Id"`
	Valores []uint64 `json:"Valores"`
}
type Evals struct {
	Id      uint64   `json:"Id"`
	Valores []uint64 `json:"Valores"`
}
type Tipo struct {
	Id       uint8 `json:"Id"`
	Precio   uint8 `json:"Precio"`
	Opciones bool  `json:"Opciones"`
}
type ListaProd struct {
	Lista []Producto `json:"Lista"`
}

type ResProd struct {
	Id      uint32 `json:"Id"`
	Nombre  string `json:"Nombre"`
	Nota    uint32 `json:"Nota"`
	Precio  uint32 `json:"Precio"`
	Calidad uint8  `json:"Calidad"`
}
type ResEmp struct {
	Id        uint32 `json:"Id"`
	Nombre    string `json:"Nombre"`
	Distancia uint32 `json:"Distancia"`
}

func (buf *Buffer) GetEmpresaBytes(lat float32, lng float32, idemp uint64, prods []Prods) {

	length := 0

	buf.Buf = AddBytes(buf.Buf, Float32bytes(lat))
	buf.Buf = AddBytes(buf.Buf, Float32bytes(lng))

	resProd := make(map[uint8][]Prods, 256)

	// TIPOS PADRE //
	for _, prod := range prods {

		if prod.Cantidad == 0 || prod.Cantidad == 1 {

			// ID NORMAL
			if prod.Id > 0 && GetBytes(prod.Id) <= 2 && GetBytes(prod.Precio) <= 2 {
				resProd[0] = append(resProd[0], prod)
				length = length + 5
			} else if prod.Id > 0 && GetBytes(prod.Id) == 3 && GetBytes(prod.Precio) <= 2 {
				resProd[1] = append(resProd[1], prod)
				length = length + 6
			} else if prod.Id > 0 && GetBytes(prod.Id) == 4 && GetBytes(prod.Precio) <= 2 {
				resProd[2] = append(resProd[2], prod)
				length = length + 7
			} else if prod.Id > 0 && GetBytes(prod.Id) <= 2 && GetBytes(prod.Precio) == 3 {
				resProd[3] = append(resProd[3], prod)
				length = length + 6
			} else if prod.Id > 0 && GetBytes(prod.Id) == 3 && GetBytes(prod.Precio) == 3 {
				resProd[4] = append(resProd[4], prod)
				length = length + 7
			} else if prod.Id > 0 && GetBytes(prod.Id) == 4 && GetBytes(prod.Precio) == 3 {
				resProd[5] = append(resProd[5], prod)
				length = length + 8
			} else if prod.Id > 0 && GetBytes(prod.Id) <= 2 && GetBytes(prod.Precio) == 4 {
				resProd[6] = append(resProd[6], prod)
				length = length + 7
			} else if prod.Id > 0 && GetBytes(prod.Id) == 3 && GetBytes(prod.Precio) == 4 {
				resProd[7] = append(resProd[7], prod)
				length = length + 8
			} else if prod.Id > 0 && GetBytes(prod.Id) == 4 && GetBytes(prod.Precio) == 4 {
				resProd[8] = append(resProd[8], prod)
				length = length + 9
			}

			// IDS NORMAL
			if prod.Ids > 0 && GetBytes(prod.Ids) <= 2 && GetBytes(prod.Precio) <= 2 {
				resProd[9] = append(resProd[9], prod)
				length = length + 5
			} else if prod.Ids > 0 && GetBytes(prod.Ids) == 3 && GetBytes(prod.Precio) <= 2 {
				resProd[10] = append(resProd[10], prod)
				length = length + 6
			} else if prod.Ids > 0 && GetBytes(prod.Ids) == 4 && GetBytes(prod.Precio) <= 2 {
				resProd[11] = append(resProd[11], prod)
				length = length + 7
			} else if prod.Ids > 0 && GetBytes(prod.Ids) <= 2 && GetBytes(prod.Precio) == 3 {
				resProd[12] = append(resProd[12], prod)
				length = length + 6
			} else if prod.Ids > 0 && GetBytes(prod.Ids) == 3 && GetBytes(prod.Precio) == 3 {
				resProd[13] = append(resProd[13], prod)
				length = length + 7
			} else if prod.Ids > 0 && GetBytes(prod.Ids) == 4 && GetBytes(prod.Precio) == 3 {
				resProd[14] = append(resProd[14], prod)
				length = length + 8
			} else if prod.Ids > 0 && GetBytes(prod.Ids) <= 2 && GetBytes(prod.Precio) == 4 {
				resProd[15] = append(resProd[15], prod)
				length = length + 7
			} else if prod.Ids > 0 && GetBytes(prod.Ids) == 3 && GetBytes(prod.Precio) == 4 {
				resProd[16] = append(resProd[16], prod)
				length = length + 8
			} else if prod.Ids > 0 && GetBytes(prod.Ids) == 4 && GetBytes(prod.Precio) == 4 {
				resProd[17] = append(resProd[17], prod)
				length = length + 9
			}
		} else {

		}
	}

	Prodbuf := []byte{}

	for tipo, v := range resProd {

		Prodbuf = AddBytes(Prodbuf, []byte{tipo})                      // TIPO
		Prodbuf = AddBytes(Prodbuf, big.NewInt(int64(len(v))).Bytes()) // CANTIDAD
		length = length + 1 + len(big.NewInt(int64(len(v))).Bytes())
		for _, prodr := range v {

			if prodr.Id > 0 {
				Prodbuf = AddBytes(Prodbuf, min2bytes(big.NewInt(int64(prodr.Id)).Bytes())) //PROD INFO
			}
			if prodr.Ids > 0 {
				Prodbuf = AddBytes(Prodbuf, min2bytes(big.NewInt(int64(prodr.Ids)).Bytes())) //PROD INFO
			}

			Prodbuf = AddBytes(Prodbuf, min2bytes(big.NewInt(int64(prodr.Precio)).Bytes())) //PROD INFO
			Prodbuf = AddBytes(Prodbuf, []byte{prodr.Calidad})

			switch tipo {
			case 0:
			case 1:
			case 2:
			}

		}
	}

	str := []byte{65, 66, 67, 68, 69, 70}
	length = length + len(big.NewInt(int64(length)).Bytes()) + 1 + len(min2bytes(big.NewInt(int64(idemp)).Bytes())) + 1 + len(str)

	buf.Buf = AddBytes(buf.Buf, big.NewInt(int64(length)).Bytes())           // LARGO EN BYTES
	buf.Buf = AddBytes(buf.Buf, []byte{uint8(len(resProd))})                 // CANTIDAD ARREGLOS
	buf.Buf = AddBytes(buf.Buf, min2bytes(big.NewInt(int64(idemp)).Bytes())) // ID EMPRESA
	buf.Buf = AddBytes(buf.Buf, []byte{6})                                   // CANTIDAD STRING
	buf.Buf = AddBytes(buf.Buf, str)                                         // STRING
	buf.Buf = append(buf.Buf, Prodbuf...)
}
func (prods *ListaProd) ReadBytes(bytes []uint8, distance uint32) {

	var j int = 0
	cantarr := bytes[j : j+1][0]
	idemp, c := GetSize2(bytes[j+1 : j+5])
	j = 1 + c

	cantstring := int(bytes[j : j+1][0])
	nombre := bytes[j+1 : j+cantstring]
	j = j + 1 + cantstring

	fmt.Printf("CANTIDAD ARR: %v | IDEMP: %v | BYTES: %v | CANTSTRING: %v | NOMBRE: %s \n", cantarr, idemp, c, cantstring, string(nombre))

	var id, precio, calidad, desde, mayor, nota uint32

	for i := uint8(0); i < cantarr; i++ {

		tipo := bytes[j : j+1][0]
		cantprod, x := GetSize1(bytes[j+1 : j+4])
		j = j + 1 + x

		for k := uint64(0); k < cantprod; k++ {

			if tipo == 0 {
				id = Bytes2toInt32(bytes[j : j+2])
				precio = Bytes2toInt32(bytes[j+2 : j+4])
				calidad = Bytes1toInt32(bytes[j+4 : j+5])
				j = j + 5
				//fmt.Printf("TIPO: 0 | ID: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 1 {
				id = Bytes3toInt32(bytes[j : j+3])
				precio = Bytes2toInt32(bytes[j+3 : j+5])
				calidad = Bytes1toInt32(bytes[j+5 : j+6])
				j = j + 6
				//fmt.Printf("TIPO: 1 | ID: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 2 {
				id = Bytes4toInt32(bytes[j : j+4])
				precio = Bytes2toInt32(bytes[j+4 : j+6])
				calidad = Bytes1toInt32(bytes[j+6 : j+7])
				j = j + 7
				//fmt.Printf("TIPO: 2 | ID: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 3 {
				id = Bytes2toInt32(bytes[j : j+2])
				precio = Bytes3toInt32(bytes[j+2 : j+5])
				calidad = Bytes1toInt32(bytes[j+5 : j+6])
				j = j + 6
				//fmt.Printf("TIPO: 3 | ID: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 4 {
				id = Bytes3toInt32(bytes[j : j+3])
				precio = Bytes3toInt32(bytes[j+3 : j+6])
				calidad = Bytes1toInt32(bytes[j+6 : j+7])
				j = j + 7
				//fmt.Printf("TIPO: 4 | ID: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 5 {
				id = Bytes4toInt32(bytes[j : j+4])
				precio = Bytes3toInt32(bytes[j+4 : j+7])
				calidad = Bytes1toInt32(bytes[j+7 : j+8])
				j = j + 8
				//fmt.Printf("TIPO: 5 | ID: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 6 {
				id = Bytes2toInt32(bytes[j : j+2])
				precio = Bytes4toInt32(bytes[j+2 : j+6])
				calidad = Bytes1toInt32(bytes[j+6 : j+7])
				j = j + 7
				//fmt.Printf("TIPO: 6 | ID: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 7 {
				id = Bytes3toInt32(bytes[j : j+3])
				precio = Bytes4toInt32(bytes[j+3 : j+7])
				calidad = Bytes1toInt32(bytes[j+7 : j+8])
				j = j + 8
				//fmt.Printf("TIPO: 7 | ID: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 8 {
				id = Bytes4toInt32(bytes[j : j+4])
				precio = Bytes4toInt32(bytes[j+4 : j+8])
				calidad = Bytes1toInt32(bytes[j+8 : j+9])
				j = j + 9
				//fmt.Printf("TIPO: 8 | ID: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}

			if tipo == 9 {
				id = Bytes2toInt32(bytes[j : j+2])
				precio = Bytes2toInt32(bytes[j+2 : j+4])
				calidad = Bytes1toInt32(bytes[j+4 : j+5])
				j = j + 5
				//fmt.Printf("TIPO: 9 | IDS: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 10 {
				id = Bytes3toInt32(bytes[j : j+3])
				precio = Bytes2toInt32(bytes[j+3 : j+5])
				calidad = Bytes1toInt32(bytes[j+5 : j+6])
				j = j + 6
				//fmt.Printf("TIPO: 10 | IDS: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 11 {
				id = Bytes4toInt32(bytes[j : j+4])
				precio = Bytes2toInt32(bytes[j+4 : j+6])
				calidad = Bytes1toInt32(bytes[j+6 : j+7])
				j = j + 7
				//fmt.Printf("TIPO: 11 | IDS: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 12 {
				id = Bytes2toInt32(bytes[j : j+2])
				precio = Bytes3toInt32(bytes[j+2 : j+5])
				calidad = Bytes1toInt32(bytes[j+5 : j+6])
				j = j + 6
				//fmt.Printf("TIPO: 12 | IDS: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 13 {
				id = Bytes3toInt32(bytes[j : j+3])
				precio = Bytes3toInt32(bytes[j+3 : j+6])
				calidad = Bytes1toInt32(bytes[j+6 : j+7])
				j = j + 7
				//fmt.Printf("TIPO: 13 | IDS: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 14 {
				id = Bytes4toInt32(bytes[j : j+4])
				precio = Bytes3toInt32(bytes[j+4 : j+7])
				calidad = Bytes1toInt32(bytes[j+7 : j+8])
				j = j + 8
				//fmt.Printf("TIPO: 14 | IDS: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 15 {
				id = Bytes2toInt32(bytes[j : j+2])
				precio = Bytes4toInt32(bytes[j+2 : j+6])
				calidad = Bytes1toInt32(bytes[j+6 : j+7])
				j = j + 7
				//fmt.Printf("TIPO: 15 | IDS: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 16 {
				id = Bytes3toInt32(bytes[j : j+3])
				precio = Bytes4toInt32(bytes[j+3 : j+7])
				calidad = Bytes1toInt32(bytes[j+7 : j+8])
				j = j + 8
				//fmt.Printf("TIPO: 16 | IDS: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}
			if tipo == 17 {
				id = Bytes4toInt32(bytes[j : j+4])
				precio = Bytes4toInt32(bytes[j+4 : j+8])
				calidad = Bytes1toInt32(bytes[j+8 : j+9])
				j = j + 9
				//fmt.Printf("TIPO: 17 | IDS: %v | PRECIO: %v | CALIDAD: %v \n", id, precio, calidad)
			}

			nota = precio + calidad + distance
			if nota < mayor && nota > desde {
				prods.Lista, mayor = PushProducto(prods.Lista, NewProd(id, "", nota, precio, distance, idemp, string(nombre)))
			}
		}
	}

	fmt.Printf("%v%v%v\n", id, precio, calidad)
}

func main() {

	resprods := ListaProd{Lista: make([]Producto, 0, 100)}
	buf := Buffer{Buf: []byte{}}
	prods := []Prods{

		Prods{Id: 1, Precio: 100, Calidad: 243},
		Prods{Id: 65800, Precio: 101, Calidad: 243},
		Prods{Id: 17000001, Precio: 102, Calidad: 243},
		Prods{Id: 2, Precio: 16000000, Calidad: 243},
		Prods{Id: 65801, Precio: 16000001, Calidad: 243},
		Prods{Id: 17000002, Precio: 16000002, Calidad: 243},
		Prods{Id: 3, Precio: 17000000, Calidad: 243},
		Prods{Id: 65802, Precio: 17000001, Calidad: 243},
		Prods{Id: 17000003, Precio: 17000002, Calidad: 243},

		Prods{Ids: 1, Precio: 100, Calidad: 243},
		Prods{Ids: 65800, Precio: 101, Calidad: 243},
		Prods{Ids: 17000001, Precio: 102, Calidad: 243},
		Prods{Ids: 2, Precio: 16000000, Calidad: 243},
		Prods{Ids: 65801, Precio: 16000001, Calidad: 243},
		Prods{Ids: 17000002, Precio: 16000002, Calidad: 243},
		Prods{Ids: 3, Precio: 17000000, Calidad: 243},
		Prods{Ids: 65802, Precio: 17000001, Calidad: 243},
		Prods{Ids: 17000003, Precio: 17000002, Calidad: 243},
	}
	buf.GetEmpresaBytes(-33.54647, 180.56575, 2, prods)

	now := time.Now()
	length := len(buf.Buf)
	j, e := 0, 0

	for {

		if length <= j {
			break
		}

		size, c := GetSize1(buf.Buf[j+8 : j+11])
		if distance := Distance(-33.44546, 70.44546, Float32frombytes(buf.Buf[j:j+4]), Float32frombytes(buf.Buf[j+4:j+8])); distance > 0 {
			j = j + 8 + c
			resprods.ReadBytes(buf.Buf[j:j+int(size)-c], distance)
		}
		j = j + int(size) - c
		e++

	}

	fmt.Println("time elapse:", time.Since(now))
	fmt.Println(resprods)

	h := &MyHandler{}
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)
}
func GetSize1(buf []byte) (size uint64, count int) {
	size = uint64(buf[0:1][0])
	if size == 255 {
		size = Bytes2toInt64(buf[0:2])
		if size == 65535 {
			size = Bytes3toInt64(buf[0:3])
			return size - 2, 3
		}
		return size - 1, 2
	}
	return size, 1
}
func GetSize2(buf []byte) (size uint64, count int) {
	size = Bytes2toInt64(buf[0:2])
	if size == 65535 {
		size = Bytes3toInt64(buf[0:3])
		return size - 2, 3
	}
	return size - 1, 2
}
func AddBytes(buf []byte, bytes []uint8) []byte {
	for _, x := range bytes {
		buf = append(buf, byte(x))
	}
	return buf
}
func GetBytes(num uint64) uint64 {

	if num <= 255 {
		return 1
	}
	if num <= 65535 {
		return 2
	}
	if num <= 16777215 {
		return 3
	}
	if num <= 4294967295 {
		return 4
	}
	if num <= 1099511627775 {
		return 5
	}
	if num <= 281474976710655 {
		return 6
	}
	if num <= 72057594037927935 {
		return 7
	}
	if num <= 18446744073709551615 {
		return 8
	}
	return 0
}

// BORRAR FUNCION Y REEMPLAZAR PRODUCTO DE LISTA Y DEVOLVER MAYOR
func RemoveProducto(s []Producto, i int) []Producto {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
func PushProducto(Prods []Producto, obj Producto) ([]Producto, uint32) {

	if len(Prods) < 100 {
		return append(Prods, obj), 65535
	} else {
		var posicion int = -1
		var mayor uint16 = obj.Nota
		var nuevomayor uint16 = obj.Nota
		var primero bool = false
		for k, v := range Prods {
			if v.Nota > mayor {
				if !primero {
					nuevomayor = mayor
				}
				mayor = v.Nota
				posicion = k
				primero = true
			}
		}
		if posicion != -1 {
			Prods = RemoveProducto(Prods, posicion)
			return append(Prods, obj), nuevomayor
		}
		return Prods, nuevomayor
	}
}
func Distance(lat1, lng1, lat2, lng2 float32) uint32 {
	first := math.Pow(float64(lat2-lat1), 2)
	second := math.Pow(float64(lng2-lng1), 2)
	return uint32(math.Sqrt(first + second))
}
func NewProd(Id uint32, Nombre string, Nota uint32, Precio uint32, Calidad uint8, Distancia uint32, Idemp uint64, NombreEmp string) Producto {
	return Producto{
		Id:        Id,
		Nombre:    Nombre,
		Nota:      Nota,
		Precio:    Precio,
		Calidad:   Calidad,
		Distancia: Distancia,
		IdEmp:     Idemp,
		NombreEmp: NombreEmp,
	}
}
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	x := ctx.QueryArgs().Peek("id")
	id := Read_uint32bytes(x)

	fmt.Printf("%v %T\n", id, id)

	//fmt.Fprintf(ctx, string(*h.cache[id]))
	fmt.Fprintf(ctx, string(id))
}
func GetMultipleBytes(bytes []uint8, j int) (uint64, int) {

	var res uint64 = uint64(bytes[j : j+1][0])
	if res < 255 {
		return res, 1
	} else {
		res = Bytes2toInt64(bytes[j+1 : j+3])
		if res < 65535 {
			return res + 1, 2
		} else {
			res = Bytes3toInt64(bytes[j+1 : j+4])
			if res < 16777215 {
				return res + 2, 3
			} else {
				res = Bytes4toInt64(bytes[j+1 : j+5])
				if res < 4294967295 {
					return res + 3, 4
				} else {
					res = Bytes5toInt64(bytes[j+1 : j+6])
					if res < 4294967295 {
						return res + 4, 5
					}
				}
			}
		}
	}
	return res, 0
}
func Read_uint32bytes(data []byte) []byte {
	var x uint32
	for _, c := range data {
		x = x*10 + uint32(c-'0')
	}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, x)
	return b
}
func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
func Float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}
func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
func Float32bytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
func Bytes3toInt32(b []uint8) uint32 {
	bytes := make([]byte, 1, 4)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint32(bytes)
}
func Bytes2toInt32(b []uint8) uint32 {
	bytes := make([]byte, 2, 4)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint32(bytes)
}
func Bytes1toInt32(b []uint8) uint32 {
	bytes := make([]byte, 3, 4)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint32(bytes)
}
func Bytes2toInt16(b []uint8) uint16 {
	return binary.BigEndian.Uint16(b)
}
func Bytes1toInt16(b []uint8) uint16 {
	bytes := make([]byte, 1, 2)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint16(bytes)
}
func Bytes4toInt32(b []uint8) uint32 {
	return binary.BigEndian.Uint32(b)
}
func Bytes2toInt64(b []uint8) uint64 {
	bytes := make([]byte, 6, 8)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint64(bytes)
}
func Bytes3toInt64(b []uint8) uint64 {
	bytes := make([]byte, 5, 8)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint64(bytes)
}
func Bytes4toInt64(b []uint8) uint64 {
	bytes := make([]byte, 4, 8)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint64(bytes)
}
func Bytes5toInt64(b []uint8) uint64 {
	bytes := make([]byte, 3, 8)
	bytes = append(bytes, b...)
	return binary.BigEndian.Uint64(bytes)
}
func InttoByte(i int) []byte {
	if i > 0 {
		return append(big.NewInt(int64(i)).Bytes(), byte(1))
	}
	return append(big.NewInt(int64(i)).Bytes(), byte(0))
}
func BytesToInt(b []byte) int {
	if b[len(b)-1] == 0 {
		return -int(big.NewInt(0).SetBytes(b[:len(b)-1]).Int64())
	}
	return int(big.NewInt(0).SetBytes(b[:len(b)-1]).Int64())
}
func min2bytes(bytes []byte) []byte {
	if len(bytes) == 1 {
		b := []byte{0}
		b = append(b, bytes[0])
		return b
	}
	return bytes
}
