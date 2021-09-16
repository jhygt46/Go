<?php

error_reporting(0);
set_time_limit(0);

$longopts  = array("help::");
$opt = getopt("f:t:c:", $longopts);


if(isset($opt["help"])){
    echo "-f nombre y ruta del archivo\n";
    echo "-t tipo => 1 filtros GO / 2 autocomplete / 3 filtros Go separado\n";
    echo "-c cantidad";
    exit;
}

if(isset($opt["t"])){

    if(isset($opt["f"])){
        if(str_contains($opt["f"], '/')) {
            $file = $opt["f"];
        }else{
            $file = "files/".$opt["f"];
        }
    }else{
        switch ($opt["t"]) {
            case 1:
                $file = "files/filtros_go.json";
                break;
            case 2:
                $file = "files/autocomplete_go.json";
                break;
            case 3:
                $file = "filtros";
                break;
            default:
                echo "Tipo ".$opt["t"]." no existe\n";
        }
    }

    switch ($opt["t"]) {
        case 1:
            if(isset($opt["c"])){
                echo "Creando ".$file." ....\n";
                $inicio = microtime(true);
                writeFileGo($opt["c"], $file);
                $fin = microtime(true);
                $diff = $fin - $inicio;
                echo "El tiempo de ejecución del archivo ha sido de " . $diff . " segundos\n";
            }else{
                echo "Debe seleccionar la cantidad ej -c 1000\n";
            }
            break;
        case 2:
            if(isset($opt["c"])){
                echo "Creando ".$file." ....\n";
                $inicio = microtime(true);
                writeFileAutoComplete($file, 2, $opt["c"]);
                $fin = microtime(true);
                $diff = $fin - $inicio;
                echo "El tiempo de ejecución del archivo ha sido de " . $diff . " segundos\n";
            }else{
                echo "Debe seleccionar la cantidad ej -c 3\n";
            }
            break;
        case 3:
            if(isset($opt["c"])){
                echo "Creando ".$opt["c"]." archivos en ".$file." ....\n";
                $inicio = microtime(true);
                writeFileDiffGo($file, $opt["c"]);
                $fin = microtime(true);
                $diff = $fin - $inicio;
                echo "El tiempo de ejecución del archivo ha sido de " . $diff . " segundos\n";
            }else{
                echo "Debe seleccionar la cantidad ej -c 3\n";
            }
            break;
        default:
            echo "Debe Seleccionar el tipo\n";
    }

}



function writeFileAutoComplete($file, $jmin, $jmax){

    $alpha = "abcdefghijklmnopqrstuvwxyz";
    $count = strlen($alpha);
    $w = 0;

    $data = '[';
    file_put_contents($file, $data, FILE_APPEND);

    for($j=$jmin; $j<=$jmax; $j++){

        $x = $count**$j;
        for($i=0; $i<$x; $i++){

            $aux = convert($i, $count, $j);
            $d = "";
            for($z=0; $z<count($aux); $z++){
                $d .= $alpha[$aux[$z]];
            }

            if($w == 0){
                $data = '{"Id": "'.$d.'","Data": [{"T": 1, "I": 1, "P": "cde"},{"T": 1, "I": 2, "P": "rdf"},{"T": 1, "I": 3, "P": "edr"},{"T": 1, "I": 4, "P": "dfe"}]}';
                $w = 1;
            }else{
                $data = ',{"Id": "'.$d.'","Data": [{"T": 1, "I": 1, "P": "cde"},{"T": 1, "I": 2, "P": "rdf"},{"T": 1, "I": 3, "P": "edr"},{"T": 1, "I": 4, "P": "dfe"}]}';
            }
            file_put_contents($file, $data, FILE_APPEND);

        }
    }

    $data = ']';
    file_put_contents($file, $data, FILE_APPEND);


}
function convert($num, $base, $min){

    $aux = $num;
    $arr = [];
    if($aux > $base){
        while($aux > $base){

            $arr[] = $aux % $base;
            $aux = floor($aux/$base);
            if($aux == $base){
                $arr[] = 0;
                $arr[] = 1;
            }
            if($aux < $base){
                $arr[] = $aux;
            }

        }
    }else{
        if($aux == $base){
            $arr[] = 0;
            $arr[] = 1;
        }
        if($aux < $base){
            $arr[] = $aux % $base;
        }
    }
    for($z=count($arr); $z<$min; $z++){
        $arr[] = 0;
    }
    return array_reverse($arr);

}
function writeFileGo($len, $file){

    $data = '[{"Id":1,"Data":{"C":[{ "T": 1, "N": "Nacionalidad", "V": ["Chilena", "Argentina", "Brasileña", "Uruguaya"] }, { "T": 2, "N": "Servicios", "V": ["Americana", "Rusa", "Bailarina", "Masaje"] },{ "T": 3, "N": "Edad" }],"E": [{ "T": 1, "N": "Rostro" },{ "T": 1, "N": "Senos" },{ "T": 1, "N": "Trasero" }]}}';
    file_put_contents($file, $data, FILE_APPEND);

    for($i=2; $i<=$len; $i++){
        $data = ',{"Id":'.$i.',"Data":{"C":[{ "T": 1, "N": "Nacionalidad", "V": ["Chilena", "Argentina", "Brasileña", "Uruguaya"] }, { "T": 2, "N": "Servicios", "V": ["Americana", "Rusa", "Bailarina", "Masaje"] },{ "T": 3, "N": "Edad" }],"E": [{ "T": 1, "N": "Rostro" },{ "T": 1, "N": "Senos" },{ "T": 1, "N": "Trasero" }]}}';
        file_put_contents($file, $data, FILE_APPEND);
    }

    $data = ']';
    file_put_contents($file, $data, FILE_APPEND);

}
function writeFileDiffGo($folder, $len){

    for($i=1; $i<=$len; $i++){
        $data = '{"Id":'.$i.',"Data":{"C":[{ "T": 1, "N": "Nacionalidad", "V": ["Chilena", "Argentina", "Brasileña", "Uruguaya"] }, { "T": 2, "N": "Servicios", "V": ["Americana", "Rusa", "Bailarina", "Masaje"] },{ "T": 3, "N": "Edad" }],"E": [{ "T": 1, "N": "Rostro" },{ "T": 1, "N": "Senos" },{ "T": 1, "N": "Trasero" }]}}';
        file_put_contents($folder."/".$i.".json", $data);
    }

}
function writeFileNode($len, $file){

    $data = '{"1":{"C":[{"T":1,"N":"Nacionalidad","V":["Chilena","Argentina","Brasileña","Uruguaya"]},{"T":2,"N":"Servicios","V":["Americana","Rusa","Bailarina","Masaje"]},{"T":3,"N":"Edad"}],"E":[{"T":1,"N":"Rostro"},{"T":1,"N":"Senos"},{"T":1,"N":"Trasero"}]}';
    file_put_contents($file, $data, FILE_APPEND);
    for($i=2; $i<=$len; $i++){
        $data = ',"'.$i.'":{"C":[{"T":1,"N":"Nacionalidad","V":["Chilena","Argentina","Brasileña","Uruguaya"]},{"T":2,"N":"Servicios","V":["Americana","Rusa","Bailarina","Masaje"]},{"T":3,"N":"Edad"}],"E":[{"T":1,"N":"Rostro"},{"T":1,"N":"Senos"},{"T":1,"N":"Trasero"}]}';
        file_put_contents($file, $data, FILE_APPEND);
    }
    $data = '}';
    file_put_contents($file, $data, FILE_APPEND);

}

?>