
const uuidChars = 'abcdefghijklmnopqrstuvwyxzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'.split('')

export function generateUuid(length = 6) {
    const output: string[] = []

    while (output.length < length) {
        output.push(uuidChars[Math.round(Math.random() * uuidChars.length)])
    }

    return output.join('')
}

export class Debouncer {
    private _debounceTimeout?: ReturnType<typeof setTimeout>
    private _wait: number = 250

    constructor(wait?: number) {
        if (wait) {
            this._wait = wait
        }
    }

    debounce<T extends (...args: any) => any>(func: Debounceable<T>, args: Parameters<T>, immediate?: boolean) {
        if (this._debounceTimeout || immediate) {
            clearTimeout(this._debounceTimeout!)
        }
        if (immediate) {
            func(...args)
        }
        this._debounceTimeout = setTimeout(() => func(...args), this._wait) as any as NodeJS.Timeout
    }
}

type Debounceable<T extends (...args: any) => any> = T

