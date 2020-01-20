const util = require('util')
const exec = util.promisify(require('child_process').exec)
const assert = require('assert')
const path = require('path')

module.exports = {
    initializeDatabase: async (POSTGRESQL_CONN_STRING, options = {}) => {
        process.env.POSTGRESQL_CONN_STRING = POSTGRESQL_CONN_STRING        
        try {
            // // creating temporal database
            let { stdout } = await exec(
                `./node_modules/.bin/sequelize-cli db:create --url '${POSTGRESQL_CONN_STRING}'`)
            if(options.logging) console.log('*** output ***\n', stdout, '\n*** END output ***')
                        
            // // migration models (tables) to database
            // stdout = (await exec(
            //     `./node_modules/.bin/sequelize-cli db:migrate --url '${POSTGRESQL_CONN_STRING}' \
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