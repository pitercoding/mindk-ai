import { useEffect, useRef } from "react";

import MessageAvatar from "@/components/chat/MessageAvatar";

import ChatMessage from "./ChatMessage";

import type { Message } from "@/types/message";

interface ChatMessageListProps {
    messages: Message[];
    isSending: boolean;
}

export default function ChatMessageList({
    messages,
    isSending,
}: ChatMessageListProps) {

    const messagesEndRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({
            behavior: "smooth",
        });
    }, [messages, isSending]);

    return (
        <div className="message-list">

            {messages.map((message) => (
                <ChatMessage
                    key={message.id}
                    message={message}
                />
            ))}

            {isSending && (
                <div
                    className="message message-assistant message-pending"
                    role="status"
                    aria-label="MindK AI is thinking"
                >
                    <MessageAvatar role="assistant" />

                    <div className="message-body">
                        <span className="typing-indicator" aria-hidden="true">
                            <span />
                            <span />
                            <span />
                        </span>
                    </div>
                </div>
            )}

            <div ref={messagesEndRef} />

        </div>
    );
}
