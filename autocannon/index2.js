'use strict'
const autocannon = require('autocannon');
async function init(){
    const instance = autocannon({
        url: 'http://3.141.24.216/get2',
        connections: 10,
        duration: 50,
        method: 'GET',
        pipelining: 4
    });
    autocannon.track(instance, {renderProgressBar: true});
}
init();