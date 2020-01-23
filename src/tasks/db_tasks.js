const util = require('util')
const path = require('path')
const exec = util.promisify(require('child_process').exec)

const { 
    POSTGRES_USER,
    POSTGRES_PASSWORD,
    POSTGRES_DB,
    POSTGRES_PORT } = process.env

module.exports = {
    execDatabaseUpdateCommands: async () => {
        // database conn string to existing default db (user)
        const POSTGRES_CONN_STRING = `postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}`
        // POSTGRES_CONN_STRING
        try {
            // // creating temporal database
            await exec(
                `./node_modules/.bin/sequelize-cli db:create --url '${POSTGRES_CONN_STRING}'`)
            console.log('*** database created ***\n')
        } catch (error) {
            if (error['stderr'].indexOf('already exists') === -1) {
                throw error
            }
            console.log('*** database already exists ***\n')
        }

        // migration models (tables) to database
        let stdout = (await exec(
            `./node_modules/.bin/sequelize-cli db:migrate --url '${POSTGRES_CONN_STRING}' \
            --migrations-path ${path.join(__dirname, '../migrations/')}`)
        ).stdout

        console.log('*** database migration: OUTPUT ***')
        console.log(stdout)
        console.log('\n*** END output ***')
    }
}