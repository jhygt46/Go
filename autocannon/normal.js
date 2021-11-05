'use strict'
const autocannon = require('autocannon');
async function init(){
    const instance = autocannon({
        url: 'http://18.216.123.130',
        connections: 10,
        duration: 50,
        method: 'GET',
        pipelining: 4
    });
    autocannon.track(instance, {renderProgressBar: true});
}
init();