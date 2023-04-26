import GraphiQL from "graphiql";
import React from "react";
import {createGraphiQLFetcher} from "@graphiql/toolkit";
import {useParams} from "react-router-dom";
import BackToHomeBtn from "../component/BackToHomeBtn";

const ServicePlayground: React.FC = () => {
    let {name} = useParams()
    let host = 'localhost'
    // host = 'kong-gateway'
    const fetcher = createGraphiQLFetcher({url: `http://${host}:8000/${name}/graphql`});

    return <div id="graphiql" style={{"height": "100vh"}}>
        <GraphiQL fetcher={fetcher}/>
        <BackToHomeBtn/>
    </div>
}

export default ServicePlayground