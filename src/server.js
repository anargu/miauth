const express = require('express')
const bodyParser = require('body-parser')

const miauthConfig = require('./config')

const { handleError } = require('./middlewares/error_handler')

const { execDatabaseUpdateCommands } = require('./tasks/db_tasks.js')

module.exports = (async function startServer () {
    try {
        await execDatabaseUpdateCommands()        
    } catch (error) {
        throw error
    }

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

    app.use(handleError);

    const host = '0.0.0.0'
    const srv = app.listen(miauthConfig.port, host, () => {
        console.log(
            `▀▄▀▄▀▄ [ MiAuth started & listening on port ${miauthConfig.port} ] ▄▀▄▀▄▀`
        )
    })

    return srv
})
