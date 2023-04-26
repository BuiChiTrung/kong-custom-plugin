import {createBrowserRouter} from "react-router-dom";
import App from "./App";
import ServicePlayground from "./containers/ServicePlayground";
import ServiceVisualize from "./containers/ServiceVisualize";
import React from "react";
import RequestGenerator from "./containers/RequestGenerator";

const router = createBrowserRouter([
    {
        path: "/",
        element: <App/>,
    },
    {
        path: ":name/playground",
        element: <ServicePlayground/>,
    },
    {
        path: ":name/visualize",
        element: <ServiceVisualize/>,
    },
    {
        path: ":name/request-generator",
        element: <RequestGenerator/>
    },
]);

export default router
