interface FormValue {
    depthProbability: number
    breadthProbability: number
    maxDepth: number
    argumentsToConsider: string[]
    requestName: string
}

interface Type {
    name: string
    fields: Field
}

interface Field {
    name: string
    args: Argument[]
}

interface Argument {
    name: string
}

export type { FormValue, Type, Field, Argument }
