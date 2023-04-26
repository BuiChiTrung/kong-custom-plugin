const { buildClientSchema, printSchema } = require("graphql");
const fs = require("fs");
const fetch = require("node-fetch");


let query = "{\n" +
    "  \"query\": \"\\n    query IntrospectionQuery {\\n      __schema {\\n        \\n        queryType { name }\\n        mutationType { name }\\n        subscriptionType { name }\\n        types {\\n          ...FullType\\n        }\\n        directives {\\n          name\\n          description\\n          \\n          locations\\n          args {\\n            ...InputValue\\n          }\\n        }\\n      }\\n    }\\n\\n    fragment FullType on __Type {\\n      kind\\n      name\\n      description\\n      \\n      fields(includeDeprecated: true) {\\n        name\\n        description\\n        args {\\n          ...InputValue\\n        }\\n        type {\\n          ...TypeRef\\n        }\\n        isDeprecated\\n        deprecationReason\\n      }\\n      inputFields {\\n        ...InputValue\\n      }\\n      interfaces {\\n        ...TypeRef\\n      }\\n      enumValues(includeDeprecated: true) {\\n        name\\n        description\\n        isDeprecated\\n        deprecationReason\\n      }\\n      possibleTypes {\\n        ...TypeRef\\n      }\\n    }\\n\\n    fragment InputValue on __InputValue {\\n      name\\n      description\\n      type { ...TypeRef }\\n      defaultValue\\n      \\n      \\n    }\\n\\n    fragment TypeRef on __Type {\\n      kind\\n      name\\n      ofType {\\n        kind\\n        name\\n        ofType {\\n          kind\\n          name\\n          ofType {\\n            kind\\n            name\\n            ofType {\\n              kind\\n              name\\n              ofType {\\n                kind\\n                name\\n                ofType {\\n                  kind\\n                  name\\n                  ofType {\\n                    kind\\n                    name\\n                  }\\n                }\\n              }\\n            }\\n          }\\n        }\\n      }\\n    }\\n  \"\n" +
    "}"
let host = 'localhost'
let name = 'country'

fetch(`http://${host}:8000/${name}/graphql`, {
    method: 'post',
    headers: { 'Content-Type': 'application/json' },
    body: query,
})
    .then((response) => {
        return response.json()
    })
    .then(response => {
       console.log(response.data)
       // const introspectionSchemaResult = JSON.parse(fs.readFileSync("result.json"));
       const graphqlSchemaObj = buildClientSchema(response.data);
       const sdlString = printSchema(graphqlSchemaObj);

       console.log(sdlString)
    })
    .catch((err) => {
        console.log(err);
    });


