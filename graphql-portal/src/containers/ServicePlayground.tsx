import GraphiQL from "graphiql";
import React from "react";
import {createGraphiQLFetcher} from "@graphiql/toolkit";
import {useParams} from "react-router-dom";
import BackToHomeBtn from "../component/BackToHomeBtn";
import {kongProxyURL} from "./constants";

const ServicePlayground: React.FC = () => {
    let {name} = useParams()
    const fetcher = createGraphiQLFetcher({url: `${kongProxyURL}/${name}/graphql`});

    return <div id="graphiql" style={{"height": "100vh"}}>
        <GraphiQL fetcher={fetcher}/>
        <BackToHomeBtn/>
    </div>
}

export default ServicePlayground
