import { statSync, mkdirSync, promises as fs, stat } from 'fs'
import { stringify } from 'querystring'
import { ErrorData } from './errorHandler'

try {
    const stat = statSync('./log/')
    if (!stat.isDirectory()) {
        mkdirSync('./log/')
    }
} catch (error) {
    mkdirSync('./log/')
}

export function log(data: ErrorData | string | object) {
    if ((data as ErrorData).error) {
        logError(data as ErrorData)
    } else if (typeof data !== 'string') {
        const loggable = `${new Date().toLocaleString()}: ${JSON.stringify(data)}`
        fs.writeFile('./log/actions.log', loggable, { flag: 'a' })
    } else if (typeof data === 'string') {
        fs.writeFile('./log/actions.log', `${new Date().toLocaleString()} ${data}`, { flag: 'a' })
    }
}

function logError(data: ErrorData) {
    const errorData = data as ErrorData
    const loggable = `${new Date().toLocaleString()}: ${errorData.code} - ${JSON.stringify(errorData.error)}\n${errorData.stack}\n${JSON.stringify(errorData.message)}\n`
    console.log(loggable)
    fs.writeFile('./log/error.log', loggable, { flag: 'a' })
}