const Redis = require('ioredis')
const { REDIS_PASSWORD, REDIS_HOST, REDIS_PORT, REDIS_DB } = process.env
const { initializeModels } = require('./models')

let re

async function initDatabase () {
    re = new Redis({
        port: REDIS_PORT || 6379,
        host: REDIS_HOST || 'localhost',
        password: REDIS_PASSWORD || 'put_here_your_long_password',
        db: parseInt(REDIS_DB || 0)
        // db 0 // default
    })

    initializeModels(re)
    return re
}

module.exports = {
    initDatabase,
    re
}
