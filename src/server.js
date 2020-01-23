const express = require('express')
const bodyParser = require('body-parser')

const miauthConfig = require('./config')

const { handleError } = require('./utils/error')

// initialized db (singleton)
const db = require('./models')
const app = express()

app.set('view engine', 'html');
app.engine('html', require('hbs').__express);

app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: true }))

app.use('/forgot', require('./routes/forgot.js')(db))
app.use('/admin', require('./routes/admin.js')(db))
app.use('/auth', require('./routes/user.js')(db))

app.use(function(err, req, res, next) {
    if (err)
        handleError(err, res)
});

const host = '0.0.0.0'
const srv = app.listen(miauthConfig.port, host, () => {
    console.log(
        `▀▄▀▄▀▄ [ MiAuth started & listening on port ${miauthConfig.PORT} ] ▄▀▄▀▄▀`
    )
})

module.exports = srv
