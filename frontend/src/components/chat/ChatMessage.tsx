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

    return (

        <div
            className={`message ${isUser
                    ? "message-user"
                    : "message-assistant"
                }`}
        >

            <strong>{isUser ? "You" : "AI"}</strong>

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
        </div>
    );
}
