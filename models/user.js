const { generateUUID } = require('../utils/auth_utils')
const RedisBase = require('./redis_base')

class User {
    static setup(re) {
        RedisBase.setup(re)
    }

    constructor(id = undefined, email = undefined, hash = undefined) {
        this.id = id;
        this.email = email;
        this.hash = hash;
    }

    static init(email, hash) {
        let _user = new User(undefined, email, hash)
        _user.id = generateUUID()
        return _user
    }

    serialize() { 
        return JSON.stringify({ id: this.id, email: this.email, hash: this.hash })        
    }

    static deserialize(rawString) {
        const obj = JSON.parse(rawString)
        let user = new User(obj.id, obj.email, obj.hash)
        return user
    }

    async save() {
        const re = RedisBase.re
        let result = await re.setnx(`users/${this.email}`, this.serialize())
        if(result === 0) {
            throw new Error('user already defined')
        } else {
            return 'OK'
        }
    }

    static async find(email) {
        const re = RedisBase.re
        let result = await re.get(`users/${email}`)
        if (!result) return null
        return User.deserialize(result)
    }
}

module.exports = User
