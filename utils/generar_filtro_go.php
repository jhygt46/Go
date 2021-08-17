<?php

error_reporting(0);
set_time_limit(0);

$options = getopt("f:p:s:");
var_dump($options);

/*
$tiempo1 = microtime(true);

writeFileGo(1000, 'go/filtros/filtros_go.json');
$tiempo2 = microtime(true);
$diff1 = $tiempo2 - $tiempo1;
echo "El tiempo de ejecución del archivo ha sido de " . $diff1 . " segundos\n";


writeFileNode(1000000, 'nodejs/data/filtros_node.json');
$tiempo3 = microtime(true);
$diff2 = $tiempo3 - $tiempo2;
echo "El tiempo de ejecución del archivo ha sido de " . $diff2 . " segundos\n";


writeFileAutoComplete('go/autocomplete/autocomplete_go.json', 2, 3);
$tiempo3 = microtime(true);
$diff3 = $tiempo3 - $tiempo2;
echo "El tiempo de ejecución del archivo ha sido de " . $diff3 . " segundos\n";


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
*/
?>