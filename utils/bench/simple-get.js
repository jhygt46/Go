var get = require('simple-get');

for(var i=0, ilen=10000; i<ilen; i++){
    get('http://localhost/filtro?id='+i, function (err, res) {
        if (err) throw err;
        console.log(res.statusCode);
    });
}

