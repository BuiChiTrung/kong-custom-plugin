import React, { useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import Prism from 'prismjs';
import 'prismjs/themes/prism.css';
import 'prismjs/components/prism-json';
import 'prismjs/components/prism-graphql';

interface Prop {
    code: string
    language: string
}

let req = 'query RandomQuery {\n            capsules {\n            dragon {\n            active\n            description\n            dry_mass_lb\n            first_flight\n            id\n            type\n        }\n            missions {\n            flight\n            name\n        }\n            status\n        }\n        }'
req = ''
const CodeHighlighter = (prop: Prop) => {
    useEffect(() => {
        Prism.highlightAll();
    }, [prop.code]);

    return (
        <pre>
      <code className={`language-${prop.language}`}>{prop.code}</code>
    </pre>
    //     <pre><code className="language-graphql">
    //                 {req}
    // </code></pre>
    );
};

CodeHighlighter.propTypes = {
    code: PropTypes.string.isRequired,
    language: PropTypes.string.isRequired,
};

export default CodeHighlighter;
