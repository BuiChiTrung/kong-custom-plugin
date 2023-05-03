import {useParams} from "react-router-dom";
import {Voyager} from "graphql-voyager";
import BackToHomeBtn from "../component/BackToHomeBtn";
import {kongProxyURL} from "./constants";


const ServiceVisualize: React.FC = () => {
    let {name} = useParams()
    function introspectionProvider(query: any) {
        return fetch(`${kongProxyURL}/${name}/graphql`, {
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
