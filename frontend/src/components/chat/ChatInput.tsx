import { useState, type FormEvent } from "react";

interface ChatInputProps {
    onSend: (message: string) => void;
    disabled?: boolean;
}

export default function ChatInput({
    onSend,
    disabled,
}: ChatInputProps) {
    const [message, setMessage] = useState("");

    function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();

        const trimmedMessage = message.trim();

        if (!trimmedMessage) {
            return;
        }

        onSend(trimmedMessage);

        setMessage("");
    }

    function handleKeyDown(
        event: React.KeyboardEvent<HTMLTextAreaElement>
    ) {

        if (
            event.key === "Enter" &&
            !event.shiftKey
        ) {

            event.preventDefault();

            const trimmedMessage = message.trim();

            if (!trimmedMessage) {
                return;
            }

            onSend(trimmedMessage);

            setMessage("");
        }
    }

    return (
        <form
            className="chat-input"
            onSubmit={handleSubmit}>
            <textarea
                placeholder="Ask something..."
                value={message}
                onChange={(event) => setMessage(event.target.value)}
                onKeyDown={handleKeyDown}
                disabled={disabled}
            />

            <button
                type="submit"
                disabled={disabled}
            >
                Send
            </button>
        </form>
    );
}
