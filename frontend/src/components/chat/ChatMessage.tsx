import { useState } from "react";

import MessageAvatar from "@/components/chat/MessageAvatar";
import MarkdownViewer from "@/components/common/MarkdownViewer";
import {
    AlertIcon,
    CheckIcon,
    CopyIcon,
    DocumentsIcon,
} from "@/components/common/icons";

import { useToast } from "@/context/toast-context";

import type { Message } from "@/types/message";


interface ChatMessageProps {
    message: Message;
}

function formatTime(value?: string): string {

    if (!value) {
        return "";
    }

    const date = new Date(value);

    if (Number.isNaN(date.getTime())) {
        return "";
    }

    return date.toLocaleTimeString(undefined, {
        hour: "2-digit",
        minute: "2-digit",
    });
}

export default function ChatMessage({
    message,
}: ChatMessageProps) {

    const toast = useToast();

    const [copied, setCopied] = useState(false);

    const isUser = message.role === "user";
    const isError = message.status === "error";
    const time = formatTime(message.createdAt);

    async function handleCopy() {

        try {
            await navigator.clipboard.writeText(message.content);

            setCopied(true);

            window.setTimeout(() => setCopied(false), 2000);

        } catch {
            toast.error("Unable to copy to the clipboard.");
        }
    }

    return (

        <div
            className={`message ${isUser
                    ? "message-user"
                    : "message-assistant"
                }${isError ? " message-error" : ""}`}
        >

            <MessageAvatar role={message.role} />

            <div className="message-body">

                <header className="message-meta">

                    <strong className="message-author">
                        {isError && <AlertIcon size={14} />}
                        {isError ? "System" : isUser ? "You" : "MindK AI"}
                    </strong>

                    {time && (
                        <time
                            className="message-time"
                            dateTime={message.createdAt}
                        >
                            {time}
                        </time>
                    )}

                    {!isUser && !isError && (
                        <button
                            type="button"
                            className="btn btn-ghost btn-icon message-copy"
                            aria-label={copied ? "Copied" : "Copy answer"}
                            title={copied ? "Copied" : "Copy answer"}
                            onClick={handleCopy}
                        >
                            {copied
                                ? <CheckIcon size={14} />
                                : <CopyIcon size={14} />}
                        </button>
                    )}

                </header>

                <MarkdownViewer content={message.content} />

                {message.sources && message.sources.length > 0 && (

                    <div className="message-sources">

                        <span className="message-sources-title">
                            Sources
                        </span>

                        <ul className="message-source-list">
                            {message.sources.map((source) => (
                                <li
                                    key={source.document_id}
                                    className="message-source-chip"
                                    title={source.name}
                                >
                                    <DocumentsIcon size={12} />

                                    <span className="message-source-name">
                                        {source.name}
                                    </span>

                                    <span className="message-source-score">
                                        {Math.round(source.score * 100)}%
                                    </span>
                                </li>
                            ))}
                        </ul>

                    </div>
                )}

            </div>

        </div>
    );
}
