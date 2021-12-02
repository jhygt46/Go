'use strict'
const autocannon = require('autocannon');
async function init(){
    const instance = autocannon({
        url: 'http://18.117.89.252/get0',
        connections: 20,
        duration: 60,
        method: 'GET',
        pipelining: 8
    });
    autocannon.track(instance, {renderProgressBar: true});
}
init();