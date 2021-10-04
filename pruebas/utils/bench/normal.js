'use strict'
const autocannon = require('autocannon');
async function init(){
    const instance = autocannon({
        url: 'http://localhost/filtro?id=1',
        connections: 10,
        duration: 60,
        method: 'GET',
        pipelining: 8
    });
    autocannon.track(instance, { renderProgressBar: true });
}
init();