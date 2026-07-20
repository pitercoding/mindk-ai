import ChatInput from "@/components/chat/ChatInput";
import ChatMessageList from "@/components/chat/ChatMessageList";
import type { ChatContext } from "@/types/chat";

import type { Message } from "@/types/message";

interface ChatProps {
    messages: Message[];
    isLoading: boolean;
    onSend: (message: string) => Promise<void>;
    context?: ChatContext;
}

export default function Chat({
    messages,
    isLoading,
    onSend,
    context,
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
