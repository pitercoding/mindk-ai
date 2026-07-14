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

    return (
        <form onSubmit={handleSubmit}>
            <input
                type="text"
                placeholder="Ask something..."
                value={message}
                onChange={(event) => setMessage(event.target.value)}
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