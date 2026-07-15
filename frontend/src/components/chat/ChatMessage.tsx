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
            <strong>
                {isUser ? "You" : "AI"}
            </strong>

            <p>
                {message.content}
            </p>
        </div>
    );
}