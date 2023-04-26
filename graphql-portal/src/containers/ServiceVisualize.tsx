// import {GraphQLVoyager, Voyager} from "graphql-voyager";
import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import {Voyager} from "graphql-voyager";
import BackToHomeBtn from "../component/BackToHomeBtn";
import {buildClientSchema} from "graphql/index";


const ServiceVisualize: React.FC = () => {
    let {name} = useParams()
    let host = 'localhost'
    function introspectionProvider(query: any) {
        return fetch(`http://${host}:8000/${name}/graphql`, {
            method: 'post',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ query: query }),
        })
            .then((response) => {
                return response.json()
            })
            .catch((err) => {
                console.log(err);
            });
    }

    return <div>
        <Voyager introspection={introspectionProvider}/>
        <BackToHomeBtn/>
    </div>
}

export default ServiceVisualize
