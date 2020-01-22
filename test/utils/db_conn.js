const util = require('util')
const exec = util.promisify(require('child_process').exec)
const assert = require('assert')
const path = require('path')

module.exports = {
    initializeDatabase: async (options = {}) => {

        const POSTGRES_USER = 'postgres',
            POSTGRES_PASSWORD = 'miauth-test',
            POSTGRES_DB = 'mocha-temporal-test-db',
            POSTGRES_PORT = '5432';
        process.env.POSTGRES_USER = POSTGRES_USER;
        process.env.POSTGRES_PASSWORD = POSTGRES_PASSWORD;
        process.env.POSTGRES_DB = POSTGRES_DB;
        process.env.POSTGRES_PORT = POSTGRES_PORT;
        const POSTGRES_CONN_STRING = `postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_USER}`
        try {
            // // creating temporal database
            let { stdout } = await exec(
                `./node_modules/.bin/sequelize-cli db:create --url '${POSTGRES_CONN_STRING}'`)
            if(options.logging) console.log('*** output ***\n', stdout, '\n*** END output ***')
                        
            // // migration models (tables) to database
            // stdout = (await exec(
            //     `./node_modules/.bin/sequelize-cli db:migrate --url '${POSTGRES_CONN_STRING}' \
            //     --migrations-path ${path.join(__dirname, '../../src/migrations/')}`)
            // ).stdout
            // if(options.logging) console.log('*** output ***\n', stdout, '\n*** END output ***')
        } catch (error) {
            if (error['stderr'].indexOf('already exists') === -1) {
                if(options.logging === true) console.error('*** stderr ***\n', error)
                assert.fail(error)
            }
        }

        // using sequelize sync fn to destroy and create db
        const db = require('../../src/models')
        await db._sequelize.sync({ force: true, logging: false })
        return db
    },
    truncateTables: async (sequelize) => {
        return await sequelize.truncate({ cascade: true, logging: false })
    }
}