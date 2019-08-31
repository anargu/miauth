const express = require('express')
const initDatabase = require('./db_conn').initDatabase
const settingUpEndpoints = require('./api').settingUpEndpoints

// env VARIABLES
// SALT bcrypt salt
// JWT_EXP_OFFSET exp offset in seconds
// JWT_SECRET secret for tokenize
// PORT port to listen express

// database stuff
initDatabase()

// Constants
const PORT = process.env.PORT || 8080
const HOST = '0.0.0.0'

const app = express()
settingUpEndpoints(app)

app.listen(PORT, HOST)
console.log(`api listening on port: ${PORT}`)
