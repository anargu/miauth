const express = require('express')
const initDatabase = require('./db_conn').initDatabase
const settingUpEndpoints = require('./api').settingUpEndpoints
const bodyParser = require('body-parser')
const { AVAILABLE_LOGIN_METHODS } = require('./constants')
const initConfig = require('./config').initConfig

async function main () {

    // env VARIABLES
    // SALT bcrypt salt
    // JWT_EXP_OFFSET exp offset in seconds
    // JWT_SECRET secret for tokenize
    // PORT port to listen express

    // database stuff
    await initDatabase()
    await initConfig()

    // Constants
    const PORT = process.env.PORT || 8080
    const HOST = '0.0.0.0'
    
    const loginBy = process.env.LOGIN_BY || null
    if (loginBy === null || (AVAILABLE_LOGIN_METHODS.indexOf(loginBy) === -1)) {
        throw new Error('incorrect LOGIN_BY environment variable selected or not setted (empty)')
    }
    

    const app = express()
    app.use(bodyParser.json())

    settingUpEndpoints(app)

    app.listen(PORT, HOST)
    console.log(`api listening on port: ${PORT}`)
}
main()