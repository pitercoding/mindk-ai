import ReactMarkdown from "react-markdown";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { oneDark } from "react-syntax-highlighter/dist/esm/styles/prism";
import remarkGfm from "remark-gfm";

import type { Components } from "react-markdown";

interface MarkdownViewerProps {
    content: string;
}

export default function MarkdownViewer({
    content,
}: MarkdownViewerProps) {

    const components: Components = {
        a({
            children,
            href,
            ...props
        }) {
            return (
                <a
                    href={href}
                    target="_blank"
                    rel="noreferrer"
                    {...props}
                >
                    {children}
                </a>
            );
        },
        code({
            children,
            className,
            ...props
        }) {
            const match = /language-(.+)/.exec(className ?? "");

            if (match) {
                return (
                    <SyntaxHighlighter
                        language={match[1]}
                        PreTag="div"
                        style={oneDark}
                        customStyle={{
                            margin: 0,
                            padding: "18px",
                            background: "transparent",
                        }}
                        codeTagProps={{
                            style: {
                                fontFamily: "inherit",
                            },
                        }}
                    >
                        {String(children).replace(/\n$/, "")}
                    </SyntaxHighlighter>
                );
            }

            return (
                <code
                    className={className}
                    {...props}
                >
                    {children}
                </code>
            );
        },
    };

    return (
        <div className="markdown-viewer">

            <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                components={components}
            >
                {content}
            </ReactMarkdown>

        </div>
    );
}
