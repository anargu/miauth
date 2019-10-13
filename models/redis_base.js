
class RedisBase {
    static setup(re) {
        this.re = re
    }

    get re() {
        return this.re
    }
}

module.exports = RedisBase