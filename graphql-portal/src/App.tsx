import React, {useState} from 'react';
import {CodeSandboxCircleFilled, BarChartOutlined, MenuFoldOutlined, MenuUnfoldOutlined, FormOutlined} from '@ant-design/icons';
import {Button, Layout, Menu, theme} from 'antd';
import 'graphiql/graphiql.min.css';
import { Link } from "react-router-dom";
import ServiceList from "./containers/ServiceList";

const {Header, Sider, Content} = Layout;

const items = [
    {
        key: '1',
        path: "/",
        icon: <CodeSandboxCircleFilled style={{ fontSize: '110%'}}/>,
        label: 'Schema Manager',
    },
    {
        key: '2',
        path: "https://buichitrung.grafana.net/d/mY9p7dQmz/kong-official?orgId=1&from=now-3h&to=now",
        icon: <BarChartOutlined style={{ fontSize: '110%'}} />,
        label: 'Service Chart',
    },
    {
        key: '3',
        path: 'https://buichitrung.grafana.net/explore?orgId=1&left=%7B%22datasource%22:%22grafanacloud-logs%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%7Blevel%3D%5C%22debug%5C%22,job%3D%5C%22proxy-cache%5C%22%7D%22,%22queryType%22:%22range%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22grafanacloud-logs%22%7D,%22editorMode%22:%22code%22%7D%5D,%22range%22:%7B%22from%22:%22now-15m%22,%22to%22:%22now%22%7D%7D',
        icon: <FormOutlined style={{ fontSize: '110%'}}/>,
        label: 'Service Log',
    },
]


const App: React.FC = () => {
    const [collapsed, setCollapsed] = useState(false);
    const {
        token: {colorBgContainer},
    } = theme.useToken();

    return (
        <Layout style={{"height": "100vh"}}>
            <Sider trigger={null} width={"250"} collapsed={collapsed}>
                <div className="logo" style={{justifyContent: collapsed ? "center" : "flex-start" }}>
                    <img src={"logo.png"} style={{"width": "50px"}}/>
                    <span style={{marginLeft: '10px', display: collapsed ? "none" : "block"}}>GraphQL Portal</span>
                </div>
                <br/>
                <Menu
                    theme="dark"
                    mode="inline"
                    defaultSelectedKeys={['1']}
                    selectedKeys={['1']}
                >
                    {
                        items.map(item => {
                            return (
                                <Menu.Item key={item.key} className="menu-item">
                                    <Link to={item.path} target={item.key === '1' ? "" : "_blank"}>
                                        {item.icon}
                                        <span>{item.label}</span>
                                    </Link>
                                </Menu.Item>
                            )
                        })
                    }
                </Menu>
            </Sider>
            <Layout className="site-layout">
                <Header style={{padding: 0, background: colorBgContainer}}>
                    <Button
                        type="text"
                        icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                        onClick={() => setCollapsed(!collapsed)}
                        style={{
                            fontSize: '16px',
                            width: 64,
                            height: 64,
                        }}
                    />
                </Header>
                <Content
                    // id="voyager"
                    style={{
                        margin: '24px 16px',
                        padding: 24,
                        minHeight: 280,
                        background: colorBgContainer,
                    }}
                >
                    <ServiceList/>
                </Content>
            </Layout>
        </Layout>
    );
};

export default App;