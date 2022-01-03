'use strict'
const autocannon = require('autocannon');
async function init(){
    const instance = autocannon({
        url: 'http://172.31.38.111/get0',
        connections: 20,
        duration: 10,
        method: 'GET',
        pipelining: 8
    });
    autocannon.track(instance, {renderProgressBar: true});
}
init();