'use strict'
const autocannon = require('autocannon');
async function init(){
    const instance = autocannon({
        url: 'http://172.31.38.111/get0',
        connections: 20,
        duration: 60,
        method: 'GET',
        pipelining: 8,
        workers: 4
    });
    autocannon.track(instance, {renderProgressBar: true});
}
init();