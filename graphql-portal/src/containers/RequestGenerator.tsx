import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import BackToHomeBtn from "../component/BackToHomeBtn";
import {Button, Col, Form, InputNumber, Row, Select} from "antd";
import CodeHighlighter from "../component/CodeHighlighter";
import "./RequestGenerator.css";
import {Argument, FormValue, Field, Type} from "./interfaces";
import {kongProxyURL, query, formInitialValues, FormLabel, Operation} from "./constants";

const {generateRandomQuery} = require("ibm-graphql-query-generator")
const {buildSchema, print} = require("graphql")
const {buildClientSchema, printSchema} = require("graphql");
const { Option } = Select;

let graphQLJSON: any;
let graphQLQueryType: Type;

const RequestGenerator: React.FC = () => {
    let {name} = useParams()
    const [form] = Form.useForm();

    const [generatedRequest, setGeneratedRequest] = useState<string>("Request")
    const [generatedVariable, setGeneratedVariable] = useState<Object>({})
    const [availableFields, setAvailableFields] = useState<Field[]>([])
    const [availableArguments, setAvailableArguments] = useState<Argument[]>([])

    const onSelectField = (fieldName: string) => {
        let chosenField: Field

        form.resetFields([FormLabel.ArgumentsToConsider])
        for (let i = 0; i < availableFields.length; i++) {
            let field = availableFields[i]
            if (field.name === fieldName) {
                chosenField = availableFields[i]
                setAvailableArguments(field.args)
                break
            }
        }

        // @ts-ignore
        if (chosenField === undefined) {
            return
        }

        graphQLQueryType.fields = [chosenField!]
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
        fetch(`${kongProxyURL}/${name}/graphql`, {
            method: 'post',
            headers: {'Content-Type': 'application/json'},
            body: query,
        })
            .then((response) => response.json())
            .then(({ data }) => {
                graphQLJSON = data
                let types = graphQLJSON["__schema"]["types"]

                for (let i = 0; i < types.length; i++) {
                    if (types[i].name !== Operation.Query) {
                        continue
                    }

                    graphQLQueryType = types[i]

                    let fields = types[i].fields
                    if (fields.length === 0) {
                        break
                    }

                    setAvailableFields(fields)
                    setAvailableArguments(fields[0].args)
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
                        name={FormLabel.DepthProbability}
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
                        name={FormLabel.BreadthProbability}
                        rules={[{ required: true, message: 'Please provide Breadth Probability'}]}
                    >
                        <InputNumber style={{ width: '100%' }} step={0.1} min={0} max={1} />
                    </Form.Item>
                </Col>
                <Col className="gutter-row" span={4}>
                    <Form.Item
                        label="Max depth"
                        name={FormLabel.MaxDepth}
                        rules={[{ required: true, message: 'Please provide Max depth'}]}
                    >
                        <InputNumber style={{ width: '100%' }} min={2} max={5}/>
                    </Form.Item>
                </Col>
                <Col className="gutter-row" span={4}>
                    <Form.Item
                        label="Field name"
                        name={FormLabel.FieldName}
                        rules={[
                            { required: true, message: 'Please select your field!' },
                        ]}
                    >
                        <Select placeholder="Please select a field" onSelect={onSelectField}>
                            {
                                availableFields.map((field: Field) =>
                                    <Option value={field.name} key={field.name}>{field.name}</Option>
                                )
                            }
                        </Select>
                    </Form.Item>
                </Col>

                <Col className="gutter-row" span={5}>
                    <Form.Item
                        name={FormLabel.ArgumentsToConsider}
                        label="Arguments to consider"
                    >
                        <Select mode="multiple">
                            {
                                availableArguments.map((arg: Argument) =>
                                    <Option value={arg.name} key={arg.name}>{arg.name}</Option>
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
