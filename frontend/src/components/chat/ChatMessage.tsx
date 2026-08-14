import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { oneDark } from "react-syntax-highlighter/dist/esm/styles/prism";

import type { Message } from "@/types/message";


interface ChatMessageProps {
    message: Message;
}

export default function ChatMessage({
    message,
}: ChatMessageProps) {

    const isUser = message.role === "user";
    const isError = message.status === "error";

    return (

        <div
            className={`message ${isUser
                    ? "message-user"
                    : "message-assistant"
                }${isError ? " message-error" : ""}`}
        >

            <strong>
                {isError ? "⚠ System" : isUser ? "You" : "AI"}
            </strong>

            <div className="markdown-viewer">

                <ReactMarkdown
                    remarkPlugins={[remarkGfm]}

                    components={{

                        code({
                            className,
                            children,
                        }) {

                            const match =
                                /language-(\w+)/.exec(
                                    className || "",
                                );

                            return match ? (

                                <SyntaxHighlighter
                                    style={oneDark}
                                    language={match[1]}
                                >
                                    {String(children).replace(
                                        /\n$/,
                                        "",
                                    )}
                                </SyntaxHighlighter>

                            ) : (

                                <code>{children}</code>

                            );
                        },
                    }}
                >

                    {message.content}

                </ReactMarkdown>

            </div>

            {message.sources && message.sources.length > 0 && (

                <div className="message-sources">

                    <span className="message-sources-title">
                        Sources
                    </span>

                    <ul>
                        {message.sources.map((source) => (
                            <li key={source.document_id}>
                                📄 {source.name}

                                <span className="message-source-score">
                                    {Math.round(source.score * 100)}%
                                </span>
                            </li>
                        ))}
                    </ul>

                </div>
            )}
        </div>
    );
}
