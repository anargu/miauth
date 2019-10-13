const Redis = require('ioredis');

// const { REDIS_PASSWORD, REDIS_HOST, REDIS_PORT } = process.env
const REDIS_PASSWORD = "put_here_your_long_password"
const REDIS_HOST = "127.0.0.1"
const REDIS_PORT = 6379

describe ('database', function () {
    it('should receive correctly env values', () => {
        console.log(`REDIS_HOST: ${REDIS_HOST} || REDIS_PASSWORD ${REDIS_PASSWORD} || REDIS_PORT ${REDIS_PORT}`)

        expect(REDIS_PORT).not.toBeNull()
        expect(REDIS_PORT).not.toBeUndefined()

        expect(REDIS_HOST).not.toBeNull()
        expect(REDIS_HOST).not.toBeUndefined()

        expect(REDIS_PASSWORD).not.toBeNull()
        expect(REDIS_PASSWORD).not.toBeUndefined()
    })

    it('shoud connect to db', async function () {
        let re;
        try {
            re = new Redis({
                port: REDIS_PORT,
                host: REDIS_HOST,
                password: REDIS_PASSWORD,
                // db: 0 // default
            })
            re.on('error', (e) => {
                console.log('ERROR ::: ', e.toString())
                expect(e.toString())
                .toBe('Error: connect ECONNREFUSED 127.0.0.1:6379')
            })
            re.on('connect', () => {
                console.log('::: CONNECTED :::')
            })
        } catch (error) {
            console.error('Unable to connect to the database:', error)
            throw error
        } finally {
            console.log('Connection has been established successfully.')
            re.quit()
            console.log('redis quited')
        }
    })
})
