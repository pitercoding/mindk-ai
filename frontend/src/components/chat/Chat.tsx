import ChatInput from "@/components/chat/ChatInput";
import ChatMessageList from "@/components/chat/ChatMessageList";

import type { Message } from "@/types/message";

interface ChatProps {
    messages: Message[];
    isLoading: boolean;
    onSend: (message: string) => Promise<void>;
}

export default function Chat({
    messages,
    isLoading,
    onSend,
}: ChatProps) {

    return (

        <>
            <section className="chat-content">

                <ChatMessageList
                    messages={messages}
                />

            </section>

            <footer className="chat-footer">

                <ChatInput
                    onSend={onSend}
                    disabled={isLoading}
                />

            </footer>
        </>
    );
}
