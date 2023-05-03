import {Space, Table, Tag} from "antd";
import axios from "axios";
import React, {useEffect, useState} from "react";
import {Link} from "react-router-dom";
import {kongProxyURL} from "./constants";

interface Service {
    id: string
    name: string,
    host: string,
    tags: string[],
    enabled: boolean,
    protocol: string,
}

const columns = [
    {
        title: 'Name',
        dataIndex: 'name',
        key: 'name',
        render: (text: string) => <strong>{text}</strong>,
        sorter: (a: Service, b: Service) => a.name.localeCompare(b.name),
    },
    {
        title: 'Protocol',
        dataIndex: 'protocol',
        key: 'protocol',
        sorter: (a: Service, b: Service) => a.protocol.localeCompare(b.protocol),
        render: (_ : any, { protocol } : Service) => {
            let color = protocol === 'http' ? 'volcano' : 'green';
            return (
                <Tag color={color} key={protocol}>
                    {protocol.toUpperCase()}
                </Tag>
            );
        },
    },
    {
        title: 'Host',
        dataIndex: 'host',
        key: 'host',
        sorter: (a: Service, b: Service) => a.host.localeCompare(b.host),
    },
    {
        title: 'Action',
        key: 'action',
        render: (_: any, record: Service) => (
            <Space size="middle">
                <Link to={`${record.name}/playground`}>Playground</Link>
                <Link to={`${record.name}/visualize`}>Visualize</Link>
                <Link to={`${record.name}/request-generator`}>Request Generator</Link>
            </Space>
        ),
    },
];

const graphQLServiceTag = "graphql"

const ServiceList: React.FC = () => {
    const [services, setServices] = useState<Service[]>([])

    useEffect(() => {
        const config = {
            method: 'get',
            maxBodyLength: Infinity,
            url: `${kongProxyURL}/services`,
        };

        axios(config)
            .then(function (response: { data: any; }) {
                return response.data
            })
            .then(({data}) => {
                data = data.filter((service: Service) => (service.tags != null && service.tags.includes(graphQLServiceTag) && service.enabled))
                setServices(data)
            })
            .catch(function (error: any) {
                console.log(error);
            });
    }, [])

    return (
        <Table dataSource={services} columns={columns}/>
    );
}

export default ServiceList
