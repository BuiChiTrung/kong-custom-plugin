import React from "react";
import {Button, Tooltip} from "antd";
import {HomeOutlined} from '@ant-design/icons';
import {Link} from "react-router-dom";

const BackToHomeBtn: React.FC = () => {
    return <Tooltip title="Home page" className="back-to-home-btn">
            <Link to="/">
                <Button shape="circle" size={"large"} icon={<HomeOutlined /> } />
            </Link>
    </Tooltip>
}

export default BackToHomeBtn