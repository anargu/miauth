const User = require('../models/user')
const { initDatabase } = require('../db_conn')
const { encodePassword } = require('../utils/auth_utils')

describe('inserting users', () => {
    let email = 'abc@abc.com'

    beforeEach(async () => {    
        this.reClient = await initDatabase()
        User.setup(this.reClient)
    })

    it('should save a user', async () => {
        try {
            const hash = await encodePassword('toto')
            
            const user = User.init(email, hash)
        
            let result = await user.save()
            expect(result).toBe('OK')
        } catch (error) {
            throw error
        }
    })
    
    afterEach(async () => {
        await this.reClient.del(`users/${email}`)
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
