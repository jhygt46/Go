'use strict'
const autocannon = require('autocannon');
async function init(){
    const instance = autocannon({
        url: 'http://localhost/filtro?id=1',
        connections: 10,
        duration: 5,
        method: 'GET',
        pipelining: 4
    });
    autocannon.track(instance, { renderProgressBar: true });
}
init();