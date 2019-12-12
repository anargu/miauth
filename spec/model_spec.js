const User = require('../models/user')
const { initDatabase } = require('../db_conn')
const { encodePassword } = require('../utils/auth_utils')
const { REDIS_USERNAME_INDEX, REDIS_EMAIL_INDEX, LOGIN_BY_USERNAME, LOGIN_BY_EMAIL } = require('../constants')

describe('inserting users (only with email)', () => {
    process.env.LOGIN_BY = LOGIN_BY_EMAIL
    let email = 'abc@abc.com'
    let username = 'abc'

    beforeEach(async () => {    
        this.reClient = await initDatabase()
        User.setup(this.reClient)
    })

    it('should save a user', async () => {
        try {
            const hash = await encodePassword('toto')
            
            const user = User.init(email, username, hash)
        
            let result = await user.save()
            expect(result).toBe('OK')
        } catch (error) {
            throw error
        }
    })
    
    afterEach(async () => {
        await this.reClient.del(`${REDIS_USERNAME_INDEX}/${username}`)
        await this.reClient.del(`${REDIS_EMAIL_INDEX}/${email}`)
        await this.reClient.hdel(`users/${email}`)
    })
})


describe('inserting users with username', () => {
    process.env.LOGIN_BY = LOGIN_BY_USERNAME
    let email = 'abc@abc.com'
    let username = 'abc'

    beforeEach(async () => {    
        this.reClient = await initDatabase()
        User.setup(this.reClient)
    })

    it('should save a user', async () => {
        try {
            const hash = await encodePassword('toto')
            
            const user = User.init(email, username, hash)
        
            let result = await user.save()
            expect(result).toBe('OK')
        } catch (error) {
            throw error
        }
    })
    
    afterEach(async () => {
        await this.reClient.del(`${REDIS_USERNAME_INDEX}/${username}`)
        await this.reClient.del(`${REDIS_EMAIL_INDEX}/${email}`)
        await this.reClient.hdel(`users/${email}`)
    })
})

describe('querying users', () => {
    let email = 'abc@abc.com'

    beforeEach(async () => {
        this.reClient = await initDatabase()
        User.setup(this.reClient)
        const hash = await encodePassword('toto')
        let user = User.init(email, hash)
        await user.save()
    })

    it('should find the email previous saved', async () => {
        try {
            let result = await User.find(email)    
            console.log('::: result ::: ', result)
            expect(result.email).toBe(email)
        } catch (error) {
            throw error
        }
    })

    it('not found user', async () => {
        try {
            let email2 = 'another@notfound.com'
            let result = await User.find(email2)
            console.log('::: result ::: ', result)
            expect(result).toBeNull()
        } catch (error) {
            throw error
        }
    })

    afterEach(() => {
        this.reClient.del(`users/${email}`)
    })
})
