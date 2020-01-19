const util = require('util')
const exec = util.promisify(require('child_process').exec)
const assert = require('assert')
const path = require('path')

module.exports = {
    initializeDatabase: async (POSTGRE_DB_URL_CONNECTION, options = {}) => {
        try {
            // creating temporal database
            let { stdout } = await exec(
                `./node_modules/.bin/sequelize-cli db:create --url '${POSTGRE_DB_URL_CONNECTION}'`)
            if(options.logging)
                console.log(
                    '*** output ***\n', stdout, '\n*** END output ***')
            // migration models (tables) to database
            stdout = (await exec(
                `./node_modules/.bin/sequelize-cli db:migrate --url '${POSTGRE_DB_URL_CONNECTION}' --migrations-path ${path.join(__dirname, '../../src/migrations/')}`)
            ).stdout
            if(options.logging)
                console.log(
                    '*** output ***\n', stdout, '\n*** END output ***')
        } catch (error) {
            if(options.logging === true)
                console.error(
                    '*** stderr ***\n', error)
            assert.fail(error)
        }
    },
    removingDatabase: async (POSTGRE_DB_URL_CONNECTION, options = {}) => {
        try {
            const { stdout } = await exec(
                `./node_modules/.bin/sequelize-cli db:drop --url '${POSTGRE_DB_URL_CONNECTION}' --migrations-path ${path.join(__dirname, '../../src/migrations/')}`)
            if(options.logging === true)
                console.log('*** output ***\n', stdout, '\n*** END output ***')
        } catch (error) {
            if(options.logging === true)
                console.error('*** stderr ***\n', error)
            assert.fail(error)
        }
    }
}