const { generateUUID } = require('../utils/auth_utils')
const RedisBase = require('./redis_base')
const { REDIS_EMAIL_INDEX, REDIS_USERNAME_INDEX, LOGIN_BY_USERNAME } = require('../constants')

class User {
    static setup(re) {
        RedisBase.setup(re)
    }

    constructor(id = undefined, email = undefined, username = undefined, hash = undefined) {
        this.id = id;
        this.email = email;
        if(process.env.LOGIN_BY === LOGIN_BY_USERNAME) {
            this.username = username;
        }
        this.hash = hash;
    }

    static init(email, username, hash) {
        let _user = new User(undefined, email, username, hash)
        _user.id = generateUUID()
        return _user
    }

    toMap() { 
        let _map = new Map()
        _map.set('id', this.id)
        _map.set('email', this.email)
        if(this.username !== undefined) {
            _map.set('username', this.username)
        }
        _map.set('hash', this.hash)
        return _map
    }

    fromMap(obj) { 
        let user = new User(obj.id, obj.email, obj.username, obj.hash)
        return user
    }

    async save() {
        const re = RedisBase.re

        try {
            const results = await re.multi()
            .setnx(`${REDIS_EMAIL_INDEX}/${this.email}`, this.email)
            .setnx(`${REDIS_USERNAME_INDEX}/${this.username}`, this.email)
            .hset(`users/${this.email}`, this.toMap())
            .exec()
            
            console.log(`results are ${results}`)

            return 'OK'
        } catch (error) {
            throw new Error('user already defined')
        }
    }

    static async find(value) {
        const loginBy = process.env.LOGIN_BY
        const re = RedisBase.re
        // getting user identifier (email)
        let identifier = await re.get(`${loginBy}/${value}`)
        if (!identifier) return null
        let result =  await re.hgetall(`users/${identifier}`)
        return User.fromMap(result)
    }

    static async findByEmail(email) {
        const re = RedisBase.re
        let result = await re.hgetall(`users/${email}`)
        if (!result) return null
        return User.fromMap(result)
    }
}

module.exports = User
