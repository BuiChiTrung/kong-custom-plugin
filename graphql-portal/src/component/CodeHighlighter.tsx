import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
import Prism from 'prismjs';
import 'prismjs/themes/prism.css';
import 'prismjs/components/prism-json';
import 'prismjs/components/prism-graphql';

interface Prop {
    code: string
    language: string
}
const CodeHighlighter = (prop: Prop) => {
    useEffect(() => {
        Prism.highlightAll();
    }, [prop.code]);

    return (
        <pre className="code-highlighter">
      <code className={`language-${prop.language}`}>{prop.code}</code>
    </pre>
    );
};

CodeHighlighter.propTypes = {
    code: PropTypes.string.isRequired,
    language: PropTypes.string.isRequired,
};

export default CodeHighlighter;
