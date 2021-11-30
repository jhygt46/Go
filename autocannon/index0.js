'use strict'
const autocannon = require('autocannon');
async function init(){
    const instance = autocannon({
        url: 'http://13.58.177.170/get',
        connections: 20,
        duration: 60,
        method: 'GET',
        pipelining: 8
    });
    autocannon.track(instance, {renderProgressBar: true});
}
init();