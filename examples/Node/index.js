var http = require('http')
const crypto = require('crypto')

const key = "sqlrestTestKey"

const sqlRestHost = 'localhost'
const sqlRestPort = 5050

function main() {
    const acceptedInput = 'Accepted inputs are: "ping", "connect", "query", and "procedure"'

    if (!process.argv[2]) {
        console.log('You must provide a function to execute')
        console.log(acceptedInput)
        return
    }

    switch(process.argv[2]) {
        case 'ping':
            ping()
            break
        case 'connect':
            connect()
            break
        case 'query':
            query()
            break
        case 'procedure':
            procedure()
            break
        default:
            console.log('Unknown input')
            console.log(acceptedInput)
    }
}

function createHmac(message) {
    const hmac = crypto.createHmac('sha256', key)
    hmac.update(message)
    return hmac.digest('hex')
}

function setAuthHeader(options, message) {
    const realm = "testing-func"
    let hmac = ''
    if (message) {
        hmac = createHmac(message)
    }
    const nonce = crypto.randomBytes(16).toString("hex")
    const timestamp = new Date().getTime()
    const authHeader = {
        Authorization: `${realm}:${hmac}:${nonce}:${timestamp}`
    }

    if (options.headers) {
        Object.assign(options.headers, authHeader)
    } else {
        options.headers = authHeader
    }

    return options
}

function ping() {
    let options = {
        hostname: sqlRestHost,
        port: sqlRestPort,
        path: '/ping'
    }

    options = setAuthHeader(options)

    http.get(options, (res) => {        
        res.on('data', (chunk) => {
            console.log('Response:', chunk.toString())
        })
    })
}

function connect() {
    let options = {
        hostname: sqlRestHost,
        port: sqlRestPort,
        path: '/connect'
    }

    options = setAuthHeader(options)

    http.get(options, (res) => {        
        res.on('data', (chunk) => {
            console.log('Response:', chunk.toString())
        })
    })
}

function query() {
    const data = JSON.stringify({
        query: "SELECT TOP 3 * FROM Flights.dbo.Airlines"
    })

    let options = {
        host: sqlRestHost,
        port: sqlRestPort,
        path: '/v1/query',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Content-Length': Buffer.byteLength(data)
        }
    }

    options = setAuthHeader(options, data)

    var req = http.request(options, (res) => {
        if (res.statusCode !== 200) {
            console.log('Non-200 status code')
        }

        res.setEncoding('utf8')
        res.on('data', function (chunk) {
            const response = JSON.parse(chunk)
            console.log(response)
        })
    })

    req.write(data)
    req.end()
}

function procedure() {
    const data = JSON.stringify({
        name: "Flights.dbo.AirportsByAirline",
        parameters: {
            airlineId: 109
        },
        executeOnly: false
    })

    let options = {
        host: sqlRestHost,
        port: sqlRestPort,
        path: '/v1/procedure',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Content-Length': Buffer.byteLength(data)
        }
    }

    options = setAuthHeader(options, data)

    var req = http.request(options, (res) => {
        if (res.statusCode !== 200) {
            console.log('Non-200 status code')
        }

        res.setEncoding('utf8')
        res.on('data', function (chunk) {
            const response = JSON.parse(chunk)
            console.log(response)
        })
    })

    req.write(data)
    req.end()
}

main()