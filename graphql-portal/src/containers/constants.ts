const query = "{\n" +
    "  \"query\": \"\\n    query IntrospectionQuery {\\n      __schema {\\n        \\n        queryType { name }\\n        mutationType { name }\\n        subscriptionType { name }\\n        types {\\n          ...FullType\\n        }\\n        directives {\\n          name\\n          description\\n          \\n          locations\\n          args {\\n            ...InputValue\\n          }\\n        }\\n      }\\n    }\\n\\n    fragment FullType on __Type {\\n      kind\\n      name\\n      description\\n      \\n      fields(includeDeprecated: true) {\\n        name\\n        description\\n        args {\\n          ...InputValue\\n        }\\n        type {\\n          ...TypeRef\\n        }\\n        isDeprecated\\n        deprecationReason\\n      }\\n      inputFields {\\n        ...InputValue\\n      }\\n      interfaces {\\n        ...TypeRef\\n      }\\n      enumValues(includeDeprecated: true) {\\n        name\\n        description\\n        isDeprecated\\n        deprecationReason\\n      }\\n      possibleTypes {\\n        ...TypeRef\\n      }\\n    }\\n\\n    fragment InputValue on __InputValue {\\n      name\\n      description\\n      type { ...TypeRef }\\n      defaultValue\\n      \\n      \\n    }\\n\\n    fragment TypeRef on __Type {\\n      kind\\n      name\\n      ofType {\\n        kind\\n        name\\n        ofType {\\n          kind\\n          name\\n          ofType {\\n            kind\\n            name\\n            ofType {\\n              kind\\n              name\\n              ofType {\\n                kind\\n                name\\n                ofType {\\n                  kind\\n                  name\\n                  ofType {\\n                    kind\\n                    name\\n                  }\\n                }\\n              }\\n            }\\n          }\\n        }\\n      }\\n    }\\n  \"\n" +
    "}"
const kongProxyURL = `http://${process.env.REACT_APP_KONG_PROXY_HOST}:${process.env.REACT_APP_KONG_PROXY_PORT}`

const formInitialValues = {
    depthProbability: 0.5,
    breadthProbability: 0.5,
    maxDepth: 4,
    fieldName: "",
    argumentsToConsider: [],
}

enum FormLabel {
    DepthProbability = "depthProbability",
    BreadthProbability = "breadthProbability",
    MaxDepth =  "maxDepth",
    FieldName =  "fieldName",
    ArgumentsToConsider = "argumentsToConsider"
}

enum Operation {
    Query = "Query",
    Mutation = "Mutation",
    Subscription = "Subscription"
}

export {query, kongProxyURL, formInitialValues, FormLabel, Operation}