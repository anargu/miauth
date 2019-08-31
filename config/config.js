module.exports = {
    development: {
        username: 'dev_user',
        password: 'd3v',
        database: 'db_dev',
        host: process.env.POSTGRES_HOST,
        dialect: 'postgres',
        operatorsAliases: false
    },
    test: {
        username: 'test_user',
        password: 't3st',
        database: 'db_test',
        host: process.env.POSTGRES_HOST,
        dialect: 'postgres',
        operatorsAliases: false
    },
    production: {
        username: process.env.POSTGRES_USER,
        password: process.env.POSTGRES_PASSWORD,
        database: process.env.POSTGRES_DB,
        host: process.env.POSTGRES_HOST,
        dialect: 'postgres',
        operatorsAliases: false
    }
}
