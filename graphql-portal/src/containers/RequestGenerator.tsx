import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import BackToHomeBtn from "../component/BackToHomeBtn";
import {Button, Col, Form, InputNumber, Row, Select} from "antd";
import CodeHighlighter from "../component/CodeHighlighter";
import "./RequestGenerator.css";
import {Argument, FormValue, Field} from "./interfaces";

const {generateRandomQuery} = require("ibm-graphql-query-generator")
const {buildSchema, print} = require("graphql")
const {buildClientSchema, printSchema} = require("graphql");
const { Option } = Select;
const query = "{\n" +
    "  \"query\": \"\\n    query IntrospectionQuery {\\n      __schema {\\n        \\n        queryType { name }\\n        mutationType { name }\\n        subscriptionType { name }\\n        types {\\n          ...FullType\\n        }\\n        directives {\\n          name\\n          description\\n          \\n          locations\\n          args {\\n            ...InputValue\\n          }\\n        }\\n      }\\n    }\\n\\n    fragment FullType on __Type {\\n      kind\\n      name\\n      description\\n      \\n      fields(includeDeprecated: true) {\\n        name\\n        description\\n        args {\\n          ...InputValue\\n        }\\n        type {\\n          ...TypeRef\\n        }\\n        isDeprecated\\n        deprecationReason\\n      }\\n      inputFields {\\n        ...InputValue\\n      }\\n      interfaces {\\n        ...TypeRef\\n      }\\n      enumValues(includeDeprecated: true) {\\n        name\\n        description\\n        isDeprecated\\n        deprecationReason\\n      }\\n      possibleTypes {\\n        ...TypeRef\\n      }\\n    }\\n\\n    fragment InputValue on __InputValue {\\n      name\\n      description\\n      type { ...TypeRef }\\n      defaultValue\\n      \\n      \\n    }\\n\\n    fragment TypeRef on __Type {\\n      kind\\n      name\\n      ofType {\\n        kind\\n        name\\n        ofType {\\n          kind\\n          name\\n          ofType {\\n            kind\\n            name\\n            ofType {\\n              kind\\n              name\\n              ofType {\\n                kind\\n                name\\n                ofType {\\n                  kind\\n                  name\\n                  ofType {\\n                    kind\\n                    name\\n                  }\\n                }\\n              }\\n            }\\n          }\\n        }\\n      }\\n    }\\n  \"\n" +
    "}"
const host = 'localhost'
const port = 8000

const formInitialValues = {
    depthProbability: 0.5,
    breadthProbability: 0.5,
    maxDepth: 4,
    requestName: "",
    argumentsToConsider: [],
}

let graphQLJSON: any;

const RequestGenerator: React.FC = () => {
    let {name} = useParams()
    const [form] = Form.useForm();

    const [generatedRequest, setGeneratedRequest] = useState<string>("Request")
    const [generatedVariable, setGeneratedVariable] = useState<Object>({})
    const [availableRequests, setAvailableRequests] = useState<Field[]>([])
    const [availableArguments, setAvailableArguments] = useState<Argument[]>([])

    const onSelectRequest = (requestName: string) => {
        let chosenRequest: any = {}

        form.resetFields(["argumentsToConsider"])
        for (let i = 0; i < availableRequests.length; i++) {
            let request = availableRequests[i]
            if (request.name === requestName) {
                chosenRequest = availableRequests[i]
                setAvailableArguments(request.args)
                break
            }
        }

        let graphQLTypes = graphQLJSON["__schema"]["types"]
        for (let i = 0; i < graphQLTypes.length; i++) {
            if (graphQLTypes[i].name === "Query") {
                graphQLTypes[i].fields = [chosenRequest]
                break
            }
        }
    }

    const onGenerateRequest = (values: FormValue) => {
        console.log(values)
        const configuration = {
            'depthProbability': values.depthProbability,
            'breadthProbability': values.breadthProbability,
            'maxDepth': values.maxDepth,
            'ignoreOptionalArguments': true,
            'argumentsToIgnore': [],
            'argumentsToConsider': values.argumentsToConsider === undefined ? [] : values.argumentsToConsider,
            'providerMap': {'*__*__*': null},
            'considerInterfaces': false,
            'considerUnions': false,
            'pickNestedQueryField': false
        }

        const graphQLSchemaObj = buildClientSchema(graphQLJSON);
        const graphQLSDL = printSchema(graphQLSchemaObj);

        const {queryDocument, variableValues } = generateRandomQuery(
            buildSchema(graphQLSDL),
            configuration
        )

        setGeneratedRequest(print(queryDocument))
        setGeneratedVariable(variableValues)
    };

    useEffect(() => {
        fetch(`http://${host}:${port}/${name}/graphql`, {
            method: 'post',
            headers: {'Content-Type': 'application/json'},
            body: query,
        })
            .then((response) => response.json())
            .then(({ data }) => {
                graphQLJSON = data
                let types = graphQLJSON["__schema"]["types"]

                for (let i = 0; i < types.length; i++) {
                    if (types[i].name !== "Query") {
                        continue
                    }

                    let requests = types[i]["fields"]
                    if (requests.length === 0) {
                        break
                    }

                    setAvailableRequests(requests)
                    setAvailableArguments(requests[0].args)
                    // TODO: trung.bc - hardcode
                    form.setFieldValue("requestName", requests[0].name)
                }
            })
            .catch((err) => {
                console.log(err);
            });
    }, [])

    return <div id="request-generator">
        <Form
            layout="vertical"
            initialValues={formInitialValues}
            onFinish={onGenerateRequest}
            form={form}
        >
            <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
                <Col className="gutter-row" span={4}>
                    <Form.Item
                        label="Depth Probability"
                        name="depthProbability"
                        rules={[
                            { required: true, message: 'Please provide Depth Probability'},
                        ]}
                    >
                        <InputNumber style={{ width: '100%' }} step={0.1} min={0} max={1}/>
                    </Form.Item>
                </Col>
                <Col className="gutter-row" span={4}>
                    <Form.Item
                        label="Breadth Probability"
                        name="breadthProbability"
                        rules={[{ required: true}]}
                    >
                        <InputNumber style={{ width: '100%' }} step={0.1} min={0} max={1} />
                    </Form.Item>
                </Col>
                <Col className="gutter-row" span={4}>
                    <Form.Item
                        label="Max depth"
                        name="maxDepth"
                        rules={[{ required: true}]}
                    >
                        <InputNumber style={{ width: '100%' }} min={1} max={7}/>
                    </Form.Item>
                </Col>
                <Col className="gutter-row" span={4}>
                    <Form.Item
                        name="requestName"
                        label="Request Name"
                        rules={[
                            { required: true, message: 'Please select your request!' },
                        ]}
                    >
                        {/* trung.bc - default option */}
                        <Select placeholder="Please select a request" onSelect={onSelectRequest}>
                            {
                                availableRequests.map((request: any) =>
                                    <Option value={request["name"]} key={request["name"]}>{request["name"]}</Option>
                                )
                            }
                        </Select>
                    </Form.Item>
                </Col>

                <Col className="gutter-row" span={5}>
                    <Form.Item
                        name="argumentsToConsider"
                        label="Arguments To Consider"
                    >
                        <Select mode="multiple">
                            {
                                availableArguments.map((arg: any) =>
                                    <Option value={arg["name"]} key={arg["name"]}>{arg["name"]}</Option>
                                )
                            }
                        </Select>
                    </Form.Item>
                </Col>
                <Col className="gutter-row" span={2} style={{"display": "flex", "justifyContent": "center", alignItems: "center"}}>
                    <Form.Item style={{"margin": 0}}>
                        <Button type="primary" htmlType="submit">
                            Generate
                        </Button>
                    </Form.Item>
                </Col>
            </Row>
        </Form>

        <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }} id="generated-request-variable">
            <Col className="gutter-row" span={16} id="generated-request">
                <div><strong>GraphQL Request</strong></div>
                <CodeHighlighter code={generatedRequest} language="graphql" />
            </Col>
            <Col className="gutter-row" span={8} id="generated-variable">
                <div><strong>GraphQL Variable</strong></div>
                <CodeHighlighter code={JSON.stringify(generatedVariable, null, "\t")} language="json" />
            </Col>
        </Row>

        <BackToHomeBtn/>
    </div>
}

export default RequestGenerator
