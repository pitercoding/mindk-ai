import ChatMessage from "./ChatMessage";

import type { Message } from "@/types/message";

interface ChatMessageListProps {
  messages: Message[];
}

export default function ChatMessageList({
    messages
}: ChatMessageListProps) {
    return (
        <div>
            {messages.map((message) => (
                <ChatMessage
                key={message.id}
                message={message}
                />
            ))}
        </div>
    );
}